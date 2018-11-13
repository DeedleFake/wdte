package wdte

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/DeedleFake/wdte/ast"
)

// Parse parses an AST from r and then translates it into a top-level
// compound. im is used to handle import statements. If im is nil, a
// no-op importer is used. In most cases, std.Import is a good default.
func Parse(r io.Reader, im Importer) (Compound, error) {
	root, err := ast.Parse(r)
	if err != nil {
		return nil, err
	}

	return FromAST(root, im)
}

// FromAST translates an AST into a top-level compound. im is used to
// handle import statements. If im is nil, a no-op importer is used.
func FromAST(root ast.Node, im Importer) (Compound, error) {
	if im == nil {
		im = ImportFunc(defaultImporter)
	}

	return (&translator{
		im: im,
	}).fromScript(root.(*ast.NTerm))
}

// An Importer creates scopes from strings. When parsing a WDTE
// script, an importer is used to import scopes into namespaces.
//
// When the WDTE import expression
//
//    import 'example'
//
// is parsed, the associated Importer will be invoked as follows:
//
//    im.Import("example")
type Importer interface {
	Import(from string) (*Scope, error)
}

func defaultImporter(from string) (*Scope, error) {
	// TODO: This should probably do something else.
	return nil, nil
}

// ImportFunc is a wrapper around simple functions to allow them to be
// used as Importers.
type ImportFunc func(from string) (*Scope, error)

func (f ImportFunc) Import(from string) (*Scope, error) { // nolint
	return f(from)
}

// ID represents a WDTE ID, such as a local variable.
type ID string

// Func is the base type through which all data is handled by WDTE. It
// represents everything that can be passed around in the language.
// This includes functions, of course, expressions, strings, numbers,
// Go functions, and anything else the client wants to pass into WDTE.
type Func interface {
	// Call calls the function with the given arguments, returning its
	// return value. frame represents the current call frame, which
	// tracks scope as well as debugging info.
	Call(frame Frame, args ...Func) Func
}

// A Frame tracks information about the current function call, such as
// the scope that the function is being executed in and debugging
// info.
type Frame struct {
	id    ID
	scope *Scope
	ctx   context.Context

	p *Frame
}

// F returns a top-level frame. This can be used by Go code calling
// WDTE functions directly if another frame is not available.
//
// In many cases, it may be preferable to use std.F() instead.
func F() Frame {
	return Frame{
		id: "unknown function, maybe Go",
	}
}

// Sub returns a new child frame of f with the given ID and the same
// scope as f.
//
// Under most circumstances, a GoFunc should call this before calling
// any WDTE functions, as it is useful for debugging. For example:
//
//    func Example(frame wdte.Frame, args ...wdte.Func) wdte.Func {
//        frame = frame.Sub("example")
//        ...
//    }
func (f Frame) Sub(id ID) Frame {
	return Frame{
		id:    id,
		scope: f.scope,
		p:     &f,
	}
}

// WithScope returns a copy of f with the given scope.
func (f Frame) WithScope(scope *Scope) Frame {
	f.scope = scope
	return f
}

// WithContext returns a copy of f with the given context.
func (f Frame) WithContext(ctx context.Context) Frame {
	f.ctx = ctx
	return f
}

// ID returns the ID of the frame. This is generally the function that
// created the frame.
func (f Frame) ID() ID {
	return f.id
}

// Scope returns the scope associated with the frame.
func (f Frame) Scope() *Scope {
	return f.scope
}

func (f Frame) Context() context.Context {
	if f.ctx == nil {
		return context.Background()
	}

	return f.ctx
}

// Parent returns the frame that this frame was created from, or a
// blank frame if there was none.
func (f Frame) Parent() Frame {
	if f.p == nil {
		return Frame{}
	}

	return *f.p
}

// Backtrace prints a backtrace to w.
func (f Frame) Backtrace(w io.Writer) error {
	_, err := fmt.Fprintf(w, "\t%v\n", f.ID())
	if err != nil {
		return err
	}

	return f.p.backtrace(w)
}

func (f *Frame) backtrace(w io.Writer) error {
	if f == nil {
		return nil
	}

	id := f.ID()
	if id == "" {
		return nil
	}

	_, err := fmt.Fprintf(w, "\tCalled from %v\n", id)
	if err != nil {
		return err
	}

	return f.p.backtrace(w)
}

// Scope is a tiered storage space for local variables. This includes
// function parameters and chain slots. A nil *Scope is equivalent to
// a blank, top-level scope.
type Scope struct {
	p       *Scope
	known   func(m map[ID]struct{})
	getFunc func(id ID) Func
}

// S is a convenience function that returns a blank, top-level scope.
func S() *Scope {
	return nil
}

// Get returns the value of the variable with the given ID. If the
// variable doesn't exist in either the current scope or any of its
// parent scopes, nil is returned.
func (s *Scope) Get(id ID) Func {
	if s == nil {
		return nil
	}

	if s.getFunc == nil {
		return s.p.Get(id)
	}

	return s.getFunc(id)
}

// Sub subscopes sub to s such that variables in sub will shadow
// variables in s.
func (s *Scope) Sub(sub *Scope) *Scope {
	return &Scope{
		p: s,
		known: func(m map[ID]struct{}) {
			sub.knownSet(m)
		},
		getFunc: func(g ID) Func {
			if v := sub.Get(g); v != nil {
				return v
			}

			return s.Get(g)
		},
	}
}

// Add returns a new subscope with the given variable stored in it.
func (s *Scope) Add(id ID, val Func) *Scope {
	return &Scope{
		p: s,
		known: func(m map[ID]struct{}) {
			m[id] = struct{}{}
		},
		getFunc: func(g ID) Func {
			if g == id {
				return s.Freeze(val)
			}

			return s.Get(g)
		},
	}
}

// Map returns a subscope that includes the given mapping of variable
// names to functions. Note that no copy is made of vars, so changing
// the map after passing it to this method may result in undefined
// behavior.
func (s *Scope) Map(vars map[ID]Func) *Scope {
	return &Scope{
		p: s,
		known: func(m map[ID]struct{}) {
			for v := range vars {
				m[v] = struct{}{}
			}
		},
		getFunc: func(g ID) Func {
			if v, ok := vars[g]; ok {
				return s.Freeze(v)
			}

			return s.Get(g)
		},
	}
}

// Custom returns a new subscope that uses the given lookup function
// to retrieve values. If getFunc returns nil, the parent of s will be
// searched. known is an optional function which adds all variables
// known to this layer of the scope into the map that it is passed as
// keys.
func (s *Scope) Custom(getFunc func(ID) Func, known func(map[ID]struct{})) *Scope {
	return &Scope{
		p: s,
		known: func(m map[ID]struct{}) {
			if known == nil {
				return
			}

			known(m)
		},
		getFunc: func(g ID) Func {
			if v := getFunc(g); v != nil {
				return v
			}

			return s.Get(g)
		},
	}
}

// Parent returns the parent of the current scope.
func (s *Scope) Parent() *Scope {
	if s == nil {
		return nil
	}

	return s.p
}

// Freeze returns a new function which executes in the scope s
// regardless of whatever Frame it is called with.
func (s *Scope) Freeze(f Func) Func {
	if sf, ok := f.(*ScopedFunc); ok {
		return sf
	}

	return &ScopedFunc{
		Func:  f,
		Scope: s,
	}
}

func (s *Scope) knownSet(vars map[ID]struct{}) {
	if s == nil {
		return
	}

	if s.known != nil {
		s.known(vars)
	}

	s.p.knownSet(vars)
}

// Known returns a sorted list of variables that are in scope.
func (s *Scope) Known() []ID {
	vars := make(map[ID]struct{})
	s.knownSet(vars)
	if len(vars) == 0 {
		return nil
	}

	list := make([]ID, 0, len(vars))
	for v := range vars {
		list = append(list, v)
	}

	sort.Slice(list, func(i1, i2 int) bool {
		return list[i1] < list[i2]
	})

	return list
}

func (s *Scope) Call(frame Frame, args ...Func) Func { // nolint
	return s
}

func (s *Scope) At(i Func) (Func, bool) { // nolint
	v := s.Get(ID(i.(String)))
	return v, v != nil
}

func (s *Scope) String() string { // nolint
	var buf strings.Builder

	buf.WriteString("scope(")
	var pre string
	for _, id := range s.Known() {
		buf.WriteString(pre)
		buf.WriteString(string(id))
		buf.WriteString(": ")
		fmt.Fprint(&buf, s.Get(id))

		pre = "; "
	}
	buf.WriteByte(')')

	return buf.String()
}

func (s *Scope) Reflect(name string) bool { // nolint
	return name == "Scope"
}

// A GoFunc is an implementation of Func that calls a Go function.
// This is the easiest way to implement lower-level systems for WDTE
// scripts to make use of.
//
// For example, to implement a simple, non-type-safe addition
// function:
//
//    GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
//      frame = frame.Sub("+")
//      var sum wdte.Number
//      for _, arg := range(args) {
//        sum += arg.Call(frame).(wdte.Number)
//      }
//      return sum
//    })
//
// If placed into a scope with the ID "+", this function can then be
// called from WDTE as follows:
//
//    + 3 6 9
//
// As shown, it is recommended that arguments be passed the given
// frame when evaluating them. Failing to do so without knowing what
// you're doing can cause unexpected behavior, including sending the
// evaluation system into infinite loops or causing panics.
//
// In the event that a GoFunc panics with an error value, it will be
// automatically caught and converted into an Error, which will then
// be returned.
type GoFunc func(frame Frame, args ...Func) Func

func (f GoFunc) Call(frame Frame, args ...Func) (r Func) { // nolint
	defer func() {
		if err, ok := recover().(error); ok {
			r = Error{
				Err: err,

				// Hmmm...
				Frame: frame.Sub("panic in GoFunc"),
			}
		}
	}()

	return f(frame, args...)
}

func (f GoFunc) String() string { // nolint
	return "<go func>"
}

// A FuncCall is an unevaluated function call. This is usually the
// right-hand side of a function declaration, but could also be any of
// various pieces of switches, compounds, or arrays.
type FuncCall struct {
	Func Func
	Args []Func

	Slot    ID
	Ignored bool
}

func (f FuncCall) Call(frame Frame, args ...Func) Func { // nolint
	if err := frame.Context().Err(); err != nil {
		return &Error{
			Frame: frame,
			Err:   err,
		}
	}

	next := make([]Func, len(f.Args))
	for i := range f.Args {
		next[i] = frame.Scope().Freeze(f.Args[i])
	}

	return f.Func.Call(frame).Call(frame, next...)
}

func (f FuncCall) String() string { // nolint
	if inner, ok := f.Func.(fmt.Stringer); ok {
		return inner.String()
	}

	return fmt.Sprint(f.Func)
}

// Chain is an unevaluated chain expression.
type Chain []*FuncCall

func (f Chain) Call(frame Frame, args ...Func) Func { // nolint
	ret := f[0].Call(frame)
	if f[0].Slot != "" {
		frame = frame.WithScope(frame.Scope().Add(f[0].Slot, ret))
	}

	for i := 1; i < len(f); i++ {
		tmp := f[i].Call(frame).Call(frame, ret)
		if f[i].Slot != "" {
			frame = frame.WithScope(frame.Scope().Add(f[i].Slot, tmp))
		}
		if !f[i].Ignored {
			ret = tmp
		}
	}
	return ret
}

// A Sub is a function that is in a subscope. This is most commonly an
// imported function.
type Sub []Func

func (sub Sub) Call(frame Frame, args ...Func) Func { // nolint
	scope := frame.Scope()
	for _, f := range sub[:len(sub)-1] {
		next := f.Call(frame.WithScope(frame.Scope().Sub(scope)))

		switch tmp := next.(type) {
		case error:
			return next
		case *Scope:
			scope = tmp
		default:
			return Error{
				Err:   fmt.Errorf("Function called on non-scope %#v", next),
				Frame: frame,
			}
		}
	}

	return sub[len(sub)-1].Call(frame.WithScope(frame.Scope().Sub(scope)), args...)
}

// A Compound represents a compound expression. Calling it calls each
// of the expressions in the compound, returning the value of the last
// one. If the compound is empty, nil is returned.
//
// If an element of a compound is a *Let, then the unevaluated
// right-hand side is placed into a new subscope under the ID
// specified by the right-hand side. The remainder of the elements in
// the compound are then evaluated under this new subscope. If the
// last element in the compound is a *Let, the right-hand side is
// returned as a lambda if it has arugments.
type Compound []Func

// Collect executes the compound the same as Call, but also returns
// the collected scope that has been modified by let expressions
// alongside the usual return value. This is useful when dealing with
// scopes as modules, as it allows you to evaluate specific functions
// in a script.
func (c Compound) Collect(frame Frame) (letScope *Scope, last Func) {
	for _, f := range c {
		switch f := f.(type) {
		case assigner:
			letScope, last = f.assign(frame, letScope)
		default:
			last = f.Call(frame.WithScope(frame.Scope().Sub(letScope)))
		}

		if _, ok := last.(error); ok {
			return letScope, last
		}
	}

	return letScope, last
}

func (c Compound) Call(frame Frame, args ...Func) Func { // nolint
	s, f := c.Collect(frame)
	return f.Call(frame.WithScope(frame.Scope().Sub(s)), args...)
}

// Switch represents a switch expression.
type Switch struct {
	// Check is the condition at the front of the switch.
	Check Func

	// Cases is the switch's cases. Each contains two functions. The
	// first index is the left-hand side, while the second is the
	// right-hand side. When the switch is evaluated, the cases are run
	// in order. If any matches, the right-hand side is evaluated and
	// its return value is returned.
	Cases [][2]Func
}

func (s Switch) Call(frame Frame, args ...Func) Func { // nolint
	check := s.Check.Call(frame)
	if _, ok := check.(error); ok {
		return check
	}

	for _, c := range s.Cases {
		lhs := c[0].Call(frame)
		if _, ok := lhs.(error); ok {
			return lhs
		}

		if lhs.Call(frame, check) == Bool(true) {
			return c[1].Call(frame)
		}
	}

	return check
}

// A Var represents a local variable. When called, it looks itself up
// in the frame that it's given and calls whatever it finds.
type Var ID

func (v Var) Call(frame Frame, args ...Func) Func { // nolint
	f := frame.Scope().Get(ID(v))
	if f == nil {
		return &Error{
			Err:   fmt.Errorf("%q is not in scope", v),
			Frame: frame,
		}
	}

	return f.Call(frame, args...)
}

// A ScopedFunc is an expression that uses a predefined scope instead
// of the one that comes with its frame. This is to make sure that a
// lazily evaluated expression has access to the correct scope. It
// caches the result of its evaluation the first time it is evaluated,
// preventing expressions that are being lazily evaluated from being
// evaluated twice.
type ScopedFunc struct {
	Func  Func
	Scope *Scope
}

func (f *ScopedFunc) Call(frame Frame, args ...Func) Func { // nolint
	r := f.Func.Call(frame.WithScope(f.Scope), args...)
	f.Func = r
	return r
}

func (f ScopedFunc) String() string { // nolint
	if inner, ok := f.Func.(fmt.Stringer); ok {
		return inner.String()
	}

	return fmt.Sprint(f.Func)
}

// A Memo wraps another function, caching the results of calls with
// the same arguments.
type Memo struct {
	Func Func
	Args []ID

	cache memoCache
}

func (m *Memo) Call(frame Frame, args ...Func) Func { // nolint
	s := frame.Scope()

	check := make([]Func, 0, len(m.Args))
	for _, id := range m.Args {
		check = append(check, s.Get(id).Call(frame))
	}

	cached, ok := m.cache.Get(check)
	if ok {
		return cached
	}

	r := m.Func.Call(frame, check...)
	m.cache.Set(check, r)
	return r
}

type memoCache struct {
	val  Func
	next map[Func]*memoCache
}

func (cache *memoCache) Get(args []Func) (Func, bool) {
	if cache == nil {
		return nil, false
	}

	if len(args) == 0 {
		return cache.val, true
	}

	if cache.next == nil {
		return nil, false
	}

	return cache.next[args[0]].Get(args[1:])
}

func (cache *memoCache) Set(args []Func, val Func) {
	if len(args) == 0 {
		cache.val = val
		return
	}

	if cache.next == nil {
		cache.next = make(map[Func]*memoCache)
	}

	n := new(memoCache)
	n.Set(args[1:], val)
	cache.next[args[0]] = n
}

// A Lambda is a closure. When called, it calls its inner expression
// with itself and its own arguments placed into the scope. In other
// words, given the lambda
//
//    (@ ex x y => + x y)
//
// it will create a new subscope containing itself under the ID "ex",
// and its first and second arguments under the IDs "x" and "y",
// respectively. It will then evaluate `+ x y` in that new scope.
//
// The arguments in the subscope, not including the self-reference,
// are contained in the boundary "args". The self-reference is contained
// in the boundary "self".
type Lambda struct {
	ID   ID
	Expr Func
	Args []ID

	Stored   []Func
	Scope    *Scope
	Original *Lambda
}

func (lambda *Lambda) Call(frame Frame, args ...Func) Func { // nolint
	stored := lambda.Stored

	scope := lambda.Scope
	if scope == nil {
		scope = frame.Scope()
	}

	original := lambda.Original
	if original == nil {
		original = lambda
	}

	if len(args) < len(lambda.Args) {
		vars := make(map[ID]Func, len(args))
		for i := range args {
			vars[lambda.Args[i]] = args[i]
		}

		return &Lambda{
			ID:   lambda.ID,
			Expr: lambda.Expr,
			Args: lambda.Args[len(args):],

			Stored:   append(stored, args...),
			Scope:    scope.Map(vars),
			Original: original,
		}
	}

	vars := make(map[ID]Func, len(args))
	for i := range lambda.Args {
		vars[lambda.Args[i]] = args[i]
	}

	scope = scope.Map(vars)
	scope = scope.Add(original.ID, original)
	return lambda.Expr.Call(frame.WithScope(scope))
}

func (lambda *Lambda) String() string { // nolint
	var buf strings.Builder

	fmt.Fprintf(&buf, "(@ %v", lambda.ID)
	for _, arg := range lambda.Args {
		buf.WriteByte(' ')
		buf.WriteString(string(arg))
	}
	buf.WriteString(" => ...)")

	return buf.String()
}

type assigner interface {
	assign(Frame, *Scope) (*Scope, Func)
}

// A Let is an expression that maps an expression to an ID. It's used
// inside compounds to create subscopes, essentially allowing for
// read-only, shadowable variable declarations.
//
// When evaluated, a Let simply passes everything through to its inner
// expression.
type Let struct {
	ID   ID
	Expr Func
}

func (let *Let) Call(frame Frame, args ...Func) Func { // nolint
	return let.Expr.Call(frame, args...)
}

func (let *Let) assign(frame Frame, scope *Scope) (*Scope, Func) {
	frame = frame.WithScope(frame.Scope().Sub(scope))

	f := frame.Scope().Freeze(let.Expr)
	return scope.Add(let.ID, f), f
}

type LetPattern struct {
	IDs  []ID
	Expr Func
}

func (let *LetPattern) Call(frame Frame, args ...Func) Func { // nolint
	return let.Expr.Call(frame, args...)
}

func (let *LetPattern) assign(frame Frame, scope *Scope) (*Scope, Func) {
	frame = frame.WithScope(frame.Scope().Sub(scope))

	switch f := let.Expr.Call(frame).(type) {
	case interface {
		Func
		Atter
		Lenner
	}:
		if f.Len() < len(let.IDs) {
			return scope, &Error{
				Err:   errors.New("Array shorter than pattern"),
				Frame: frame,
			}
		}

		for i, id := range let.IDs {
			v, ok := f.At(Number(i))
			if !ok {
				return scope, &Error{
					Err:   errors.New("Array shorter than pattern"),
					Frame: frame,
				}
			}

			scope = scope.Add(id, frame.Scope().Freeze(v))
		}
		return scope, frame.Scope().Freeze(f)

	case interface {
		Func
		Atter
	}:
		for i, id := range let.IDs {
			v, ok := f.At(Number(i))
			if !ok {
				return scope, &Error{
					Err:   errors.New("Array shorter than pattern"),
					Frame: frame,
				}
			}

			scope = scope.Add(id, frame.Scope().Freeze(v))
		}
		return scope, frame.Scope().Freeze(f)

	default:
		return scope, &Error{
			Err:   fmt.Errorf("Invalid pattern matching type: %T", f),
			Frame: frame,
		}
	}
}
