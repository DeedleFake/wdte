package wdte

import (
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
// is used to handle import statements.
func Parse(r io.Reader, im Importer) (*Module, error) {
	root, err := ast.Parse(r)
	if err != nil {
		return nil, err
	}

	return FromAST(root, im)
}

// FromAST translates an AST into a module. im is used to handle
// import statements.
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
	Call(frame []Func, args ...Func) Func

	// Equals returns true if one function equals another one. This is
	// meaningless for the majority of implementations, and should
	// largely be ignored by clients. When implementing custom
	// functions, if you're not sure what to do for this, you probably
	// won't have any issues. This is only used, by default, by switches
	// as a way of checking if the condition matches one of the cases.
	Equals(other Func) bool
}

// A GoFunc is an implementation of Func that calls a Go function.
// This is the easiest way to implement lower-level systems for WDTE
// scripts to make use of.
//
// For example, to implement a simple, non-type-safe addition
// function:
//
//    module.Funcs["+"] = GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
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
type GoFunc func(frame []Func, args ...Func) Func

func (f GoFunc) Call(frame []Func, args ...Func) Func {
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

func (f DeclFunc) Call(frame []Func, args ...Func) Func {
	if len(args) < f.Args {
		return &DeclFunc{
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
	return f.Expr.Call(next, next...)
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

func (f Expr) Call(frame []Func, args ...Func) Func {
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

func (f Chain) Call(frame []Func, args ...Func) Func {
	return f.Func.Call(frame, f.Args...).Call(frame, f.Prev.Call(frame))
}

func (f Chain) Equals(other Func) bool {
	panic("Not implemented.")
}

// A String is a string, as parsed from a string literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type String string

func (s String) Call(frame []Func, args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

func (s String) Equals(other Func) bool {
	o, ok := other.(String)
	return ok && (s == o)
}

// A Number is a number, as parsed from a number literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type Number float64

func (n Number) Call(frame []Func, args ...Func) Func {
	// TODO: Use the arguments for something, perhaps.
	return n
}

func (n Number) Equals(other Func) bool {
	o, ok := other.(Number)
	return ok && (n == o)
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

func (e External) Call(frame []Func, args ...Func) Func {
	return e.Module.Imports[e.Import].Funcs[e.Func].Call(frame, args...)
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

func (local Local) Call(frame []Func, args ...Func) Func {
	return local.Module.Funcs[local.Func].Call(frame, args...)
}

func (local Local) Equals(other Func) bool {
	o, ok := other.(Local)
	return ok && (local.Func == o.Func)
}

// A Compound represents a compound expression. Calling it calls each
// of the expressions in the compound, returning the value of the last
// one. If the compound is empty, nil is returned.
type Compound []Func

func (c Compound) Call(frame []Func, args ...Func) Func {
	var last Func
	for _, f := range c {
		last = f.Call(frame)
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

func (s Switch) Call(frame []Func, args ...Func) Func {
	check := s.Check.Call(frame)
	for _, c := range s.Cases {
		if (c[0] == nil) || (check.Equals(c[0].Call(frame))) {
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

func (a Arg) Call(frame []Func, args ...Func) Func {
	if int(a) >= len(frame) {
		// TODO: Handle this properly.
		panic("Argument out of frame.")
	}

	return frame[a].Call(frame, args...)
}

func (a Arg) Equals(other Func) bool {
	panic("Not implemented.")
}

// A FramedFunc is a function which keeps track of its own calling
// frame.
type FramedFunc struct {
	// Func is the actual function.
	Func Func

	// Frame is the function to call Func with.
	Frame []Func
}

func (f FramedFunc) Call(frame []Func, args ...Func) Func {
	return f.Func.Call(f.Frame, args...)
}

func (f FramedFunc) Equals(other Func) bool {
	panic("Not implemented.")
}

// An Array represents a WDTE array type. It's similar to a Compound,
// but doesn't evaluate its own members. Instead, evaluation simply
// returns itself, much like strings and numbers.
type Array []Func

func (a Array) Call(frame []Func, args ...Func) Func {
	return a
}

func (a Array) Equals(other Func) bool {
	o, ok := other.(Array)
	if !ok || (len(a) != len(o)) {
		return false
	}

	for i := range a {
		if !a[i].Equals(o[i]) {
			return false
		}
	}

	return true
}
