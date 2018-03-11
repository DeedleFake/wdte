package wdte

import (
	"fmt"
	"io"

	"github.com/DeedleFake/wdte/ast"
)

// A Module is the result of parsing a WDTE script. It is the main
// type of this entire library.
type Module struct {
	// Funcs maps IDs to functions. WDTE import statements and function
	// declarations create these when parsed.
	Funcs map[ID]Func
}

// Parse parses an AST from r and then translates it into a module. im
// is used to handle import statements. If im is nil, a no-op importer
// is used.
func Parse(r io.Reader, im Importer) (*Module, error) {
	return new(Module).Parse(r, im)
}

// FromAST translates an AST into a module. im is used to handle
// import statements. If im is nil, a no-op importer is used.
func FromAST(root ast.Node, im Importer) (*Module, error) {
	return new(Module).FromAST(root, im)
}

// Parse parses an AST from r and then translates it into a module. im
// is used to handle import statements. If im is nil, a no-op importer
// is used.
func (m *Module) Parse(r io.Reader, im Importer) (*Module, error) {
	root, err := ast.ParseScript(r)
	if err != nil {
		return nil, err
	}

	return m.FromAST(root, im)
}

// FromAST translates an AST into a module. im is used to handle
// import statements. If im is nil, a no-op importer is used.
func (m *Module) FromAST(root ast.Node, im Importer) (*Module, error) {
	return m.fromScript(root.(*ast.NTerm), im)
}

// Insert inserts the functions from n into m. This is different from
// an import as the functions are inserted into m's namespace. This is
// the preferred way of using the standard library.
//
// Insert returns m to allow for chaining.
func (m *Module) Insert(n *Module) *Module {
	if n == nil {
		return m
	}

	if m.Funcs == nil {
		m.Funcs = make(map[ID]Func)
	}

	for id := range n.Funcs {
		m.Funcs[id] = Local{
			Module: n,
			Func:   id,
		}
	}

	return m
}

// Eval parses an expression in the context of the module. It returns
// the expression unevaluated, despite the name.
//
// BUG: Due to #30, this doesn't usually work as expected, and should
// probably be avoided for now.
func (m *Module) Eval(r io.Reader) (Func, error) {
	expr, err := ast.ParseExpr(r)
	if err != nil {
		return nil, err
	}

	return m.fromExpr(expr.(*ast.NTerm), nil), nil
}

func (m *Module) Call(frame Frame, args ...Func) Func { // nolint
	return m
}

// An Importer creates modules from strings. When parsing a WDTE
// script, an importer is used to import modules.
//
// When the WDTE import statement
//
//    'example' => e;
//
// is parsed, the associated Importer will be invoked as follows:
//
//    im.Import("example")
//
// The return value will then be added to the module's Funcs map.
type Importer interface {
	Import(from string) (*Module, error)
}

// ImportFunc is a wrapper around simple functions to allow them to be
// used as Importers.
type ImportFunc func(from string) (*Module, error)

func (f ImportFunc) Import(from string) (*Module, error) { // nolint
	return f(from)
}

// ID represents a WDTE ID, such as a function or imported module
// name.
type ID string

// Func is the base type through which all data is handled by WDTE. It
// represents everything that can be passed around in the language.
// This includes functions, of course, expressions, strings, numbers,
// Go functions, and anything else the client wants to pass into WDTE.
type Func interface {
	// Call calls the function with the given arguments, returning its
	// return value. frame represents the current call frame. This is
	// used to keep track of function arguments during the evaluation of
	// expressions, and can largely be ignored by clients.
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

// A Frame tracks information about the current function call.
type Frame struct {
	id    ID
	scope *Scope

	p *Frame
}

// F returns a top-level frame. This can be used by Go code calling
// WDTE functions directly if another frame is not available.
func F() Frame {
	return Frame{
		id: "unknown function, maybe Go",
	}
}

// New creates a new frame from a previous frame. id should be the ID
// of the function that generated the frame, and args should be the
// arguments given to that function. The previous frame is set as the
// new frame's parent.
func (f Frame) New(id ID, args map[ID]Func) Frame {
	return Frame{
		id:    id,
		scope: S().Map(args),
		p:     &f,
	}
}

// Sub returns a sub-scoped frame that has the variables in scope
// added to its scope as new tier.
func (f Frame) Sub(scope map[ID]Func) Frame {
	f.scope = f.scope.Map(scope)
	return f
}

// WithID is a convienence function for creating a frame from another
// frame with a new ID, the same scope, and the other frame as its
// parent.
func (f Frame) WithID(id ID) Frame {
	return Frame{
		id:    id,
		scope: f.scope,
		p:     &f,
	}
}

// WithScope returns a frame with the same ID and parent frame as the
// current one, but with the given scope.
func (f Frame) WithScope(scope *Scope) Frame {
	return Frame{
		id:    f.id,
		scope: scope,
		p:     f.p,
	}
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
	getFunc func(id ID) Func
}

// S is a convienence function that returns a blank, top-level scope.
func S() *Scope {
	return nil
}

// Get returns value of the variable with the given id. If the
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
		getFunc: func(g ID) Func {
			if g == id {
				return val
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
		getFunc: func(g ID) Func {
			if v, ok := vars[g]; ok {
				return v
			}

			return s.Get(g)
		},
	}
}

// Custom returns a new subscope that uses the given lookup function
// to retrieve values.
func (s *Scope) Custom(getFunc func(ID) Func) *Scope {
	return &Scope{
		p:       s,
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

// Freeze returns a new function executes in the scope s regardless of
// whatever Frame it is called with.
func (s *Scope) Freeze(f Func) Func {
	return &ScopedFunc{
		Func:  f,
		Scope: s,
	}
}

// A GoFunc is an implementation of Func that calls a Go function.
// This is the easiest way to implement lower-level systems for WDTE
// scripts to make use of.
//
// For example, to implement a simple, non-type-safe addition
// function:
//
//    module.Funcs["+"] = GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
//      var sum wdte.Number
//      for _, arg := range(args) {
//        sum += arg.Call(frame).(wdte.Number)
//      }
//      return sum
//    })
//
// This can then be called from WDTE as follows:
//
//    + 3 6 9
//
// As shown, it is recommended that arguments be passed the given
// frame when evaluating them. Failing to do so without knowing what
// you're doing can cause unexpected behavior, including sending the
// evaluation system into infinite loops or causing panics.
type GoFunc func(frame Frame, args ...Func) Func

func (f GoFunc) Call(frame Frame, args ...Func) (r Func) { // nolint
	defer func() {
		if err, ok := recover().(error); ok {
			r = Error{
				Err: err,

				// Hmmm...
				Frame: frame.WithID("panic in GoFunc"),
			}
		}
	}()

	return f(frame, args...)
}

// A DeclFunc is a function that was declared in a WDTE function
// declaration. This is the primary source of the frame argument that
// is passed around everywhere.
type DeclFunc struct {
	// ID is the name that the function was declared with. This is
	// primarily for debugging.
	ID ID

	// Expr is the expression that the function maps to.
	Expr Func

	Args []ID

	// Stored is the arguments that have already been passed to a
	// function if it was given less arguments than it was declared
	// with.
	Stored []Func
}

func (f DeclFunc) Call(frame Frame, args ...Func) Func { // nolint
	if len(f.Stored)+len(args) < len(f.Args) {
		return &DeclFunc{
			ID:     f.ID,
			Expr:   f.Expr,
			Args:   f.Args,
			Stored: append(f.Stored, args...),
		}
	}

	next := make([]Func, 0, len(f.Stored)+len(args))
	for _, arg := range f.Stored {
		next = append(next, frame.Scope().Freeze(arg))
	}
	for _, arg := range args {
		next = append(next, frame.Scope().Freeze(arg))
	}

	vars := make(map[ID]Func, len(f.Args))
	for i, arg := range next {
		vars[f.Args[i]] = arg
	}

	return f.Expr.Call(frame.New(f.ID, vars), next...)
}

// An Expr is an unevaluated expression. This is usually the
// right-hand side of a function declaration, but could also be any of
// various pieces of switches, compounds, or arrays.
type Expr struct {
	// Func is the underlying function.
	Func Func

	// Args are the arguments to pass to Func.
	Args []Func

	Chain Func

	Slot ID
}

func (f Expr) Call(frame Frame, args ...Func) Func { // nolint
	n := f.Func.Call(frame, f.Args...)
	frame = frame.Sub(map[ID]Func{
		f.Slot: frame.Scope().Freeze(n),
	})

	return f.Chain.Call(frame, n)
}

// Chain is an unevaluated chain expression.
type Chain struct {
	// Func is the expression at the current part of the chain.
	Func Func

	// Args is the arguments to Func.
	Args []Func

	Chain Func

	Slot ID
}

func (f Chain) Call(frame Frame, args ...Func) Func { // nolint
	n := f.Func.Call(frame, f.Args...).Call(frame, args[0])
	frame = frame.Sub(map[ID]Func{
		f.Slot: frame.Scope().Freeze(n),
	})

	return f.Chain.Call(frame, n)
}

// IgnoredChain is an unevaluated chain expression that returns the
// previous expression's return value, rather than its own.
type IgnoredChain struct {
	// Func is the expression at the current part of the chain.
	Func Func

	// Args is the arguments to Func.
	Args []Func

	Chain Func

	Slot ID
}

func (f IgnoredChain) Call(frame Frame, args ...Func) Func { // nolint
	n := f.Func.Call(frame, f.Args...).Call(frame, args[0])
	frame = frame.Sub(map[ID]Func{
		f.Slot: frame.Scope().Freeze(n),
	})

	return f.Chain.Call(frame, args[0])
}

// An EndChain is a no-op that just returns its own first argument.
// This is used as the last element of a chain.
type EndChain struct {
}

func (f EndChain) Call(frame Frame, args ...Func) Func { // nolint
	return args[0]
}

// External represents a function from an imported module. It looks
// the function up when called, so it is safe to pass Externals around
// before importing, so long as they are not evaluated until after
// importing.
type External struct {
	// Module is the module that the function was called from. This is
	// *not* the module that the function was declared in. Unless, for
	// some reason, the module has itself as an import.
	Module *Module

	// Import is the expression on the left-hand side of an external
	// function call. If this does not return a module, the call will
	// fail.
	Import Func

	// Func is the ID of the function in the module it was declared in.
	Func ID
}

func (e External) Call(frame Frame, args ...Func) Func { // nolint
	im := e.Import.Call(frame)
	i, ok := im.(*Module)
	if !ok {
		return Error{
			Err:   fmt.Errorf("Function called on non-module %#v", im),
			Frame: frame,
		}
	}
	f, ok := i.Funcs[e.Func]
	if !ok {
		return Error{
			Err:   fmt.Errorf("Function %q does not exist in import %q", e.Func, e.Import),
			Frame: frame,
		}
	}

	return f.Call(frame, args...)
}

func (e External) Compare(other Func) (int, bool) { // nolint
	o, ok := other.(External)
	if ok && (e.Import == o.Import) && (e.Func == o.Func) {
		return 0, false
	}

	return -1, false
}

// Local represents a function from a module, usually the current one.
// It looks the function up when called, so it is safe to pass Locals
// around before importing, so long as they are not evaluated until
// after importing.
type Local struct {
	// Module is the module that the function was declared in.
	Module *Module

	// Func is the ID of the function in the module.
	Func ID
}

func (local Local) Call(frame Frame, args ...Func) Func { // nolint
	f, ok := local.Module.Funcs[local.Func]
	if !ok {
		return Error{
			Err:   fmt.Errorf("Function %q does not exist", local.Func),
			Frame: frame,
		}
	}

	return f.Call(frame, args...)
}

func (local Local) Compare(other Func) (int, bool) { // nolint
	o, ok := other.(Local)
	if ok && (local.Func == o.Func) {
		return 0, false
	}

	return -1, false
}

// A Compound represents a compound expression. Calling it calls each
// of the expressions in the compound, returning the value of the last
// one. If the compound is empty, nil is returned.
type Compound []Func

func (c Compound) Call(frame Frame, args ...Func) Func { // nolint
	var last Func
	for _, f := range c {
		last = f.Call(frame)
		if _, ok := last.(error); ok {
			return last
		}
	}

	return last
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
// with itself and its own arguments appended to its frame.
type Lambda struct {
	ID ID

	// Expr is the expression that the lambda maps to.
	Expr Func

	Args []ID

	// Stored is the arguments that have already been passed to a
	// lambda if it was given less arguments than it was declared with.
	Stored []Func
}

func (lambda *Lambda) Call(frame Frame, args ...Func) Func { // nolint
	if len(lambda.Stored)+len(args) < len(lambda.Args) {
		return &Lambda{
			ID:     lambda.ID,
			Expr:   lambda.Expr,
			Args:   lambda.Args,
			Stored: append(lambda.Stored, args...),
		}
	}

	next := make([]Func, 0, len(lambda.Args))
	for _, arg := range lambda.Stored {
		next = append(next, frame.Scope().Freeze(arg))
	}
	for _, arg := range args {
		next = append(next, frame.Scope().Freeze(arg))
	}

	vars := make(map[ID]Func, 1+len(lambda.Args))
	vars[lambda.ID] = lambda
	for i, arg := range next {
		vars[lambda.Args[i]] = arg
	}

	return lambda.Expr.Call(frame.Sub(vars), next...)
}
