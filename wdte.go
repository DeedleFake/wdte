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
	id   ID
	args []Func

	p           *Frame
	cline, ccol int
}

// F returns a top-level frame. This can be used by Go code calling
// WDTE functions directly if another frame is not available.
func F() Frame {
	return Frame{
		id: "unknown function, maybe Go",
	}
}

func CustomFrame(id ID, args []Func, parent *Frame) Frame { // nolint
	return Frame{
		id:   id,
		args: args,
		p:    parent,
	}
}

// New creates a new frame from a previous frame. id should be the ID
// of the function that generated the frame, and args should be the
// arguments given to that function.
func (f Frame) New(id ID, args []Func) Frame {
	return Frame{
		id:   id,
		args: args,
		p:    &f,
	}
}

// Sub returns a sub-scoped frame that has args appended to its
// argument list.
func (f Frame) Sub(args []Func) Frame {
	f.args = append(f.args, args...)
	return f
}

// Pos builds a new frame with position information. This is primarily
// intended for internal use.
func (f Frame) Pos(line, col int) Frame {
	f.cline = line
	f.ccol = col
	return f
}

// WithID is a convienence function for creating a frame with a new ID
// but the same arguments as the previous frame.
func (f Frame) WithID(id ID) Frame {
	return Frame{
		id:   id,
		args: f.args,

		p:     &f,
		cline: f.cline,
		ccol:  f.ccol,
	}
}

// ID returns the ID of the frame. This is generally the function that
// created the frame.
func (f Frame) ID() ID {
	return f.id
}

// Args returns the arguments of the frame.
func (f Frame) Args() []Func {
	return f.args
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

	_, err := fmt.Fprintf(w, "\tCalled from %v (%v:%v)\n", id, f.cline, f.ccol)
	if err != nil {
		return err
	}

	return f.p.backtrace(w)
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

	// Args is the number of arguments the function expects.
	Args int

	// Stored is the arguments that have already been passed to a
	// function if it was given less arguments than it was declared
	// with.
	Stored []Func
}

func (f DeclFunc) Call(frame Frame, args ...Func) Func { // nolint
	if len(args) < f.Args {
		return &DeclFunc{
			ID:     f.ID,
			Expr:   f,
			Args:   f.Args - len(args),
			Stored: args,
		}
	}

	next := make([]Func, 0, len(f.Stored)+len(args))
	for _, arg := range f.Stored {
		next = append(next, &FramedFunc{
			Func:  arg,
			Frame: frame,
		})
	}
	for _, arg := range args {
		next = append(next, &FramedFunc{
			Func:  arg,
			Frame: frame,
		})
	}

	return f.Expr.Call(frame.New(f.ID, next), next...)
}

// An Expr is an unevaluated expression. This is usually the
// right-hand side of a function declaration, but could also be any of
// various pieces of switches, compounds, or arrays.
type Expr struct {
	// Func is the underlying function.
	Func Func

	// Args are the arguments to pass to Func.
	Args []Func

	// Slots is a chain-specific global storage space for expression
	// slots. This pointer is the same for every instance of Chain in
	// the same chain, as well as for the initial expression.
	Slots *[]Func
}

func (f Expr) Call(frame Frame, args ...Func) Func { // nolint
	r := f.Func.Call(frame, f.Args...)
	*f.Slots = append(*f.Slots, r)
	return r
}

// Chain is an unevaluated chain expression.
type Chain struct {
	// Func is the expression at the current part of the chain.
	Func Func

	// Args is the arguments to Func.
	Args []Func

	// Prev is the previous part of the chain.
	Prev Func

	// Slots is a chain-specific global storage space for expression
	// slots. This pointer is the same for every instance of Chain in
	// the same chain, as well as for the initial expression.
	Slots *[]Func
}

func (f Chain) Call(frame Frame, args ...Func) Func { // nolint
	prev := f.Prev.Call(frame)

	frame = frame.Sub(*f.Slots)
	r := f.Func.Call(frame, f.Args...).Call(frame, prev)

	*f.Slots = append(*f.Slots, r)
	return r
}

// IgnoredChain is an unevaluated chain expression that returns the
// previous expression's return value, rather than its own.
type IgnoredChain struct {
	// Func is the expression at the current part of the chain.
	Func Func

	// Args is the arguments to Func.
	Args []Func

	// Prev is the previous part of the chain.
	Prev Func

	// Slots is a chain-specific global storage space for expression
	// slots. This pointer is the same for every instance of Chain in
	// the same chain, as well as for the initial expression.
	Slots *[]Func
}

func (f IgnoredChain) Call(frame Frame, args ...Func) Func { // nolint
	prev := f.Prev.Call(frame)

	frame = frame.Sub(*f.Slots)
	r := f.Func.Call(frame, f.Args...).Call(frame, prev)

	*f.Slots = append(*f.Slots, r)
	return prev
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

	Line, Col int
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

	return f.Call(frame.Pos(e.Line, e.Col), args...)
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

	Line, Col int
}

func (local Local) Call(frame Frame, args ...Func) Func { // nolint
	f, ok := local.Module.Funcs[local.Func]
	if !ok {
		return Error{
			Err:   fmt.Errorf("Function %q does not exist", local.Func),
			Frame: frame,
		}
	}

	return f.Call(frame.Pos(local.Line, local.Col), args...)
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

// Arg represents an argument in the current frame. It is the opposite
// end from DeclFunc of the frame argument that gets passed around all
// over the place.
type Arg int

func (a Arg) Call(frame Frame, args ...Func) Func { // nolint
	if int(a) >= len(frame.Args()) {
		// I don't think this can happen normally, but some GoFunc
		// somewhere could generate an Arg for some bizarre reason and
		// cause this.
		return Error{
			Err: fmt.Errorf(
				"Attempted to access %vth argument in a frame containing %v",
				a,
				len(frame.Args()),
			),
			Frame: frame,
		}
	}

	return frame.Args()[a].Call(frame, args...)
}

// A FramedFunc is a function which keeps track of its own calling
// frame.
type FramedFunc struct {
	// Func is the actual function.
	Func Func

	// Frame is the frame to call Func with.
	Frame Frame
}

func (f FramedFunc) Call(frame Frame, args ...Func) Func { // nolint
	return f.Func.Call(f.Frame, args...)
}

// A Memo wraps another function, caching the results of calls with
// the same arguments.
type Memo struct {
	Func Func

	cache memoCache
}

func (m *Memo) Call(frame Frame, args ...Func) Func { // nolint
	for i := range args {
		args[i] = args[i].Call(frame)
	}

	cached, ok := m.cache.Get(args)
	if ok {
		return cached
	}

	r := m.Func.Call(frame, args...)
	m.cache.Set(args, r)
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
	// Expr is the expression that the lambda maps to.
	Expr Func

	// Args is the number of arguments the lambda expects.
	Args int

	// Stored is the arguments that have already been passed to a
	// lambda if it was given less arguments than it was declared with.
	Stored []Func
}

func (lambda *Lambda) Call(frame Frame, args ...Func) Func { // nolint
	if len(args) < lambda.Args {
		return &Lambda{
			Expr:   lambda,
			Args:   lambda.Args - len(args),
			Stored: args,
		}
	}

	next := make([]Func, 1, 1+len(lambda.Stored)+len(args))
	next[0] = lambda
	for _, arg := range lambda.Stored {
		next = append(next, &FramedFunc{
			Func:  arg,
			Frame: frame,
		})
	}
	for _, arg := range args {
		next = append(next, &FramedFunc{
			Func:  arg,
			Frame: frame,
		})
	}

	return lambda.Expr.Call(frame.Sub(next), next...)
}
