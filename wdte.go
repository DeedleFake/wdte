package wdte

import (
	"bytes"
	"fmt"
	"io"

	"github.com/DeedleFake/wdte/ast"
)

// A Module is the result of parsing a WDTE script. It is the main
// type of this entire library.
type Module struct {
	// Imports maps IDs to other modules. WDTE import statements create
	// these when parsed.
	Imports map[ID]*Module

	// Funcs maps IDs to functions. WDTE function declarations create
	// these when parsed.
	Funcs map[ID]Func
}

// Parse parses an AST from r and then translates it into a module. im
// is used to handle import statements. If im is nil, a no-op importer
// is used.
func Parse(r io.Reader, im Importer) (*Module, error) {
	root, err := ast.Parse(r)
	if err != nil {
		return nil, err
	}

	return FromAST(root, im)
}

// FromAST translates an AST into a module. im is used to handle
// import statements. If im is nil, a no-op importer is used.
func FromAST(root ast.Node, im Importer) (*Module, error) {
	return fromScript(root.(*ast.NTerm), im)
}

// An Importer creates modules from strings. When parsing a WDTE script, an importer is used to import modules.
//
// When the WDTE import statement
//
//    'example' => e;
//
// is parsed, the associated Importer will be invoked as follows:
//
//    im.Import(example)
//
// The return value will then be added to the module's Imports map.
type Importer interface {
	Import(from string) (*Module, error)
}

// ImportFunc is a wrapper around simple functions to allow them to be
// used as Importers.
type ImportFunc func(from string) (*Module, error)

func (f ImportFunc) Import(from string) (*Module, error) {
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

	// Equals returns true if one function equals another one. This is
	// meaningless for the majority of implementations, and should
	// largely be ignored by clients. When implementing custom
	// functions, if you're not sure what to do for this, you probably
	// won't have any issues. This is only used, by default, by switches
	// as a way of checking if the condition matches one of the cases.
	Equals(other Func) bool
}

// A Frame tracks information about the current function call.
type Frame struct {
	id   ID
	args []Func

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
// arguments given to that function.
func (f Frame) New(id ID, args []Func) Frame {
	return Frame{
		id:   id,
		args: args,
		p:    &f,
	}
}

// WithID is a convienence function for creating a frame with a new ID
// but the same arguments as the previous frame.
func (f Frame) WithID(id ID) Frame {
	return f.New(id, f.args)
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

	_, err := fmt.Fprintf(w, "\tCalled from %v\n", id)
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

func (f GoFunc) Call(frame Frame, args ...Func) (r Func) {
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

func (f GoFunc) Equals(other Func) bool {
	panic("Not implemented.")
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

func (f DeclFunc) Call(frame Frame, args ...Func) Func {
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

func (f DeclFunc) Equals(other Func) bool {
	panic("Not implemented.")
}

// An Expr is an unevaluated expression. This is usually the
// right-hand side of a function declaration, but could also be any of
// various pieces of switches, compounds, or arrays.
type Expr struct {
	// Func is the underlying function.
	Func Func

	// Args are the arguments to pass to Func.
	Args []Func
}

func (f Expr) Call(frame Frame, args ...Func) Func {
	return f.Func.Call(frame, f.Args...)
}

func (f Expr) Equals(other Func) bool {
	panic("Not implemented.")
}

// Chain is an unevaluated chain expression.
type Chain struct {
	// Func is the expression at the current part of the chain.
	Func Func

	// Args is the arguments to Func.
	Args []Func

	// Prev is the previous part of the chain.
	Prev Func
}

func (f Chain) Call(frame Frame, args ...Func) Func {
	return f.Func.Call(frame, f.Args...).Call(frame, f.Prev.Call(frame))
}

func (f Chain) Equals(other Func) bool {
	panic("Not implemented.")
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

	// Import is the import ID of the module that the function was
	// declared in.
	Import ID

	// Func is the ID of the function in the module it was declared in.
	Func ID
}

func (e External) Call(frame Frame, args ...Func) Func {
	i, ok := e.Module.Imports[e.Import]
	if !ok {
		return Error{
			Err:   fmt.Errorf("Import %q does not exist", e.Import),
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

func (e External) Equals(other Func) bool {
	o, ok := other.(External)
	return ok && (e.Import == o.Import) && (e.Func == o.Func)
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

func (local Local) Call(frame Frame, args ...Func) Func {
	f, ok := local.Module.Funcs[local.Func]
	if !ok {
		return Error{
			Err:   fmt.Errorf("Function %q does not exist", local.Func),
			Frame: frame,
		}
	}

	return f.Call(frame, args...)
}

func (local Local) Equals(other Func) bool {
	o, ok := other.(Local)
	return ok && (local.Func == o.Func)
}

// A Compound represents a compound expression. Calling it calls each
// of the expressions in the compound, returning the value of the last
// one. If the compound is empty, nil is returned.
type Compound []Func

func (c Compound) Call(frame Frame, args ...Func) Func {
	var last Func
	for _, f := range c {
		last = f.Call(frame)
		if _, ok := last.(error); ok {
			return last
		}
	}

	return last
}

func (c Compound) Equals(other Func) bool {
	panic("Not implemented.")
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

func (s Switch) Call(frame Frame, args ...Func) Func {
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

		if check.Equals(lhs) {
			return c[1].Call(frame)
		}
	}

	return nil
}

func (s Switch) Equals(other Func) bool {
	panic("Not implemented.")
}

// Arg represents an argument in the current frame. It is the opposite
// end from DeclFunc of the frame argument that gets passed around all
// over the place.
type Arg int

func (a Arg) Call(frame Frame, args ...Func) Func {
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

func (a Arg) Equals(other Func) bool {
	panic("Not implemented.")
}

// A FramedFunc is a function which keeps track of its own calling
// frame.
type FramedFunc struct {
	// Func is the actual function.
	Func Func

	// Frame is the frame to call Func with.
	Frame Frame
}

func (f FramedFunc) Call(frame Frame, args ...Func) Func {
	return f.Func.Call(f.Frame, args...)
}

func (f FramedFunc) Equals(other Func) bool {
	panic("Not implemented.")
}

type Memo struct {
	Func Func

	cache memoCache
}

func (m *Memo) Call(frame Frame, args ...Func) Func {
	if m.cache == nil {
		m.cache = make(memoCache)
	}

	for i := range args {
		args[i] = args[i].Call(frame)
	}

	key := m.cache.Key(args)

	cached, ok := m.cache[key]
	if ok {
		return cached
	}

	r := m.Func.Call(frame, args...)
	m.cache[key] = r
	return r
}

func (m Memo) Equals(other Func) bool {
	panic("Not implemented.")
}

type memoCache map[string]Func

func (cache memoCache) Key(args []Func) string {
	var buf bytes.Buffer
	for _, arg := range args {
		// TODO: Figure out a way to do this that isn't stupid.
		fmt.Fprint(&buf, arg)
	}
	return buf.String()
}
