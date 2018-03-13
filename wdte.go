package wdte

import (
	"fmt"
	"io"
	"sort"

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

// A Comparer is a Func that is able to be compared to other
// functions.
type Comparer interface {
	// Compare returns two values. The meaning of the first is dependent
	// upon the second. If the second is true, then the first indicates
	// ordering via the standard negative, positive, and zero results to
	// indicate less than, greater than, and equal, respectively. If the
	// second is false, then the first indicates only equality, with
	// zero still meaning equal, but other values simply meaning unequal.
	Compare(other Func) (int, bool)
}

// A Frame tracks information about the current function call, such as
// the scope that the function is being executed in and debugging
// info.
type Frame struct {
	id    ID
	scope *Scope

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

// ID returns the ID of the frame. This is generally the function that
// created the frame.
func (f Frame) ID() ID {
	return f.id
}

// Scope returns the scope associated with the frame.
func (f Frame) Scope() *Scope {
	return f.scope
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
	known   map[ID]struct{}
	getFunc func(id ID) Func
}

// S is a convenience function that returns a blank, top-level scope.
func S() *Scope {
	return nil
}

// Get returns the value of the variable with the given id. If the
// variable doesn't exist in either the current scope or any of its
// parent scopes, nil is returned.
func (s *Scope) Get(id ID) Func {
	if s == nil {
		return nil
	}

	return s.getFunc(id)
}

// Sub returns a new subscope with the given variable stored in it.
func (s *Scope) Sub(id ID, val Func) *Scope {
	return &Scope{
		p: s,
		known: map[ID]struct{}{
			id: {},
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
	known := make(map[ID]struct{}, len(vars))
	for v := range vars {
		known[v] = struct{}{}
	}

	return &Scope{
		p:     s,
		known: known,
		getFunc: func(g ID) Func {
			if v, ok := vars[g]; ok {
				return s.Freeze(v)
			}

			return s.Get(g)
		},
	}
}

// Custom returns a new subscope that uses the given lookup function
// to retrieve values. vars is an optional list of known variables for
// listing purposes.
func (s *Scope) Custom(getFunc func(ID) Func, vars ...ID) *Scope {
	known := make(map[ID]struct{}, len(vars))
	for _, v := range vars {
		known[v] = struct{}{}
	}

	return &Scope{
		p:       s,
		known:   known,
		getFunc: getFunc,
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
	return &ScopedFunc{
		Func:  f,
		Scope: s,
	}
}

func (s *Scope) knownSet(vars map[ID]struct{}) {
	if s == nil {
		return
	}

	for v := range s.known {
		vars[v] = struct{}{}
	}

	s.p.knownSet(vars)
}

// Known returns a list of variables that are in scope.
func (s *Scope) Known() []ID {
	vars := make(map[ID]struct{})
	s.knownSet(vars)

	var list []ID
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

// An Expr is an unevaluated expression. This is usually the
// right-hand side of a function declaration, but could also be any of
// various pieces of switches, compounds, or arrays.
type Expr struct {
	Func  Func
	Args  []Func
	Chain Func
	Slot  ID
}

func (f Expr) Call(frame Frame, args ...Func) Func { // nolint
	next := make([]Func, len(f.Args))
	for i := range f.Args {
		next[i] = frame.Scope().Freeze(f.Args[i])
	}

	n := f.Func.Call(frame, next...)
	frame = frame.WithScope(frame.Scope().Sub(f.Slot, n))

	return f.Chain.Call(frame, n)
}

// Chain is an unevaluated chain expression.
type Chain struct {
	Func  Func
	Args  []Func
	Chain Func
	Slot  ID
}

func (f Chain) Call(frame Frame, args ...Func) Func { // nolint
	next := make([]Func, len(f.Args))
	for i := range f.Args {
		next[i] = frame.Scope().Freeze(f.Args[i])
	}

	n := f.Func.Call(frame, next...).Call(frame, frame.Scope().Freeze(args[0]))
	frame = frame.WithScope(frame.Scope().Sub(f.Slot, n))

	return f.Chain.Call(frame, n)
}

// IgnoredChain is an unevaluated chain expression that returns the
// previous expression's return value, rather than its own.
type IgnoredChain struct {
	Func  Func
	Args  []Func
	Chain Func
	Slot  ID
}

func (f IgnoredChain) Call(frame Frame, args ...Func) Func { // nolint
	next := make([]Func, len(f.Args))
	for i := range f.Args {
		next[i] = frame.Scope().Freeze(f.Args[i])
	}

	n := f.Func.Call(frame, next...).Call(frame, frame.Scope().Freeze(args[0]))
	frame = frame.WithScope(frame.Scope().Sub(f.Slot, n))

	return f.Chain.Call(frame, args[0])
}

// An EndChain is a no-op that just returns its own first argument.
// This is used as the last element of a chain.
type EndChain struct {
}

func (f EndChain) Call(frame Frame, args ...Func) Func { // nolint
	return args[0]
}

// A Sub is a function that is in a subscope. This is most commonly an
// imported function.
type Sub struct {
	// Scope is a function that returns the subscope. If it does not
	// return a *Scope, calling the Sub will fail.
	Scope Func

	// Func is the ID of the function being called.
	Func ID
}

func (sub Sub) Call(frame Frame, args ...Func) Func { // nolint
	ns := sub.Scope.Call(frame)
	scope, ok := ns.(*Scope)
	if !ok {
		return Error{
			Err:   fmt.Errorf("Function called on non-scope %#v", ns),
			Frame: frame,
		}
	}

	f := scope.Get(sub.Func)
	if f == nil {
		return Error{
			Err:   fmt.Errorf("Function %q does not exist in subscope", sub.Func),
			Frame: frame,
		}
	}

	return f.Call(frame, args...)
}

func (sub Sub) Compare(other Func) (int, bool) { // nolint
	o, ok := other.(Sub)
	if ok && (sub.Scope == o.Scope) && (sub.Func == o.Func) {
		return 0, false
	}

	return -1, false
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
// returned as a lambda.
type Compound []Func

// Collect executes the compound the same as Call, but also returns
// the collected scope that has been modified by let expressions
// alongside the usual return value. This is useful when dealing with
// scopes as modules, as it allows you to evaluate specific functions
// in a script.
func (c Compound) Collect(frame Frame, args ...Func) (*Scope, Func) {
	var last Func
	for _, f := range c {
		switch f := f.(type) {
		case *Let:
			frame = frame.WithScope(frame.Scope().Sub(f.ID, f.Expr))
			last = f
		default:
			last = f.Call(frame)
			if _, ok := last.(error); ok {
				return frame.Scope(), last
			}
		}
	}

	if let, ok := last.(*Let); ok {
		last = let.Call(frame)
	}

	if len(args) > 0 {
		last = last.Call(frame, args...)
	}

	return frame.Scope(), last
}

func (c Compound) Call(frame Frame, args ...Func) Func { // nolint
	_, f := c.Collect(frame, args...)
	return f
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
	//
	// A default case is represented by a nil in the first index. It is
	// possible to have cases after a default, but pointless, as a
	// default is always run when it is encountered.
	Cases [][2]Func
}

func (s Switch) Call(frame Frame, args ...Func) Func { // nolint
	check := s.Check.Call(frame)
	if _, ok := check.(error); ok {
		return check
	}

	for _, c := range s.Cases {
		if c[0] == nil {
			return c[1].Call(frame)
		}

		lhs := c[0].Call(frame)
		if _, ok := lhs.(error); ok {
			return lhs
		}

		if lhs.Call(frame, check) == Bool(true) {
			return c[1].Call(frame)
		}
	}

	return nil
}

// A Var represents a local variable. When called, it looks itself up
// in the frame that it's given and calls whatever it finds.
type Var ID

func (v Var) Call(frame Frame, args ...Func) Func { // nolint
	return frame.Scope().Get(ID(v)).Call(frame, args...)
}

// A ScopedFunc is an expression that uses a predefined scope instead
// of the one that comes with its frame. This is to make sure that a
// lazily evaluated expression has access to the correct scope.
type ScopedFunc struct {
	// Func is the actual function.
	Func Func

	Scope *Scope
}

func (f ScopedFunc) Call(frame Frame, args ...Func) Func { // nolint
	return f.Func.Call(frame.WithScope(f.Scope), args...)
}

// A Memo wraps another function, caching the results of calls with
// the same arguments.
type Memo struct {
	Func Func

	cache memoCache
}

func (m *Memo) Call(frame Frame, args ...Func) Func { // nolint
	check := make([]Func, 0, len(args))
	for _, arg := range args {
		check = append(check, arg.Call(frame))
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

	scope = scope.Map(vars).Sub(original.ID, original)
	return lambda.Expr.Call(frame.WithScope(scope), append(stored, args...)...)
}

// A Let is an expression that maps an expression to an ID. It's used
// inside compounds to create subscopes, essentially allowing for
// read-only, shadowable variable declaration.
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
