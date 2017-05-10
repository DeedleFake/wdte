package wdte

import (
	"io"

	"github.com/DeedleFake/wdte/ast"
)

type Module struct {
	Imports map[ID]*Module
	Funcs   map[ID]Func
}

func Parse(r io.Reader, im Importer) (*Module, error) {
	root, err := ast.Parse(r)
	if err != nil {
		return nil, err
	}

	return FromAST(root, im)
}

func FromAST(root ast.Node, im Importer) (*Module, error) {
	return fromScript(root.(*ast.NTerm), im)
}

type Importer interface {
	Import(from string) (*Module, error)
}

type ImportFunc func(from string) (*Module, error)

func (f ImportFunc) Import(from string) (*Module, error) {
	return f(from)
}

type ID string

type Func interface {
	// TODO: Handle errors.
	Call(frame []Func, args ...Func) Func
	Equals(other Func) bool
}

type GoFunc func(frame []Func, args ...Func) Func

func (f GoFunc) Call(frame []Func, args ...Func) Func {
	return f(frame, args...)
}

func (f GoFunc) Equals(other Func) bool {
	panic("Not implemented.")
}

type DeclFunc struct {
	Expr   Func
	Args   int
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

	frame = append(f.Stored, args...)
	return f.Expr.Call(frame, frame...)
}

func (f DeclFunc) Equals(other Func) bool {
	panic("Not implemented.")
}

type Expr struct {
	Func Func
	Args []Func
}

func (f Expr) Call(frame []Func, args ...Func) Func {
	return f.Func.Call(frame, f.Args...)
}

func (f Expr) Equals(other Func) bool {
	panic("Not implemented.")
}

type Chain struct {
	Func Func
	Args []Func
	Prev Func
}

func (f Chain) Call(frame []Func, args ...Func) Func {
	return f.Func.Call(frame, f.Args...).Call(frame, f.Prev.Call(frame))
}

func (f Chain) Equals(other Func) bool {
	panic("Not implemented.")
}

type String string

func (s String) Call(frame []Func, args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

func (s String) Equals(other Func) bool {
	o, ok := other.(String)
	return ok && (s == o)
}

type Number float64

func (n Number) Call(frame []Func, args ...Func) Func {
	// TODO: Use the arguments for something, perhaps.
	return n
}

func (n Number) Equals(other Func) bool {
	o, ok := other.(Number)
	return ok && (n == o)
}

type External struct {
	Module *Module
	Import ID
	Func   ID
}

func (e External) Call(frame []Func, args ...Func) Func {
	return e.Module.Imports[e.Import].Funcs[e.Func].Call(frame, args...)
}

func (e External) Equals(other Func) bool {
	o, ok := other.(External)
	return ok && (e.Import == o.Import) && (e.Func == o.Func)
}

type Local struct {
	Module *Module
	Func   ID
}

func (local Local) Call(frame []Func, args ...Func) Func {
	return local.Module.Funcs[local.Func].Call(frame, args...)
}

func (local Local) Equals(other Func) bool {
	o, ok := other.(Local)
	return ok && (local.Func == o.Func)
}

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

type Switch struct {
	Check Func
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
