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
	Call(args ...Func) Func
}

type GoFunc func(args ...Func) Func

func (f GoFunc) Call(args ...Func) Func {
	return f(args...)
}

type DeclFunc struct {
	Expr   Func
	Args   int
	Stored []Func
}

func (f DeclFunc) Call(args ...Func) Func {
	if len(args) < f.Args {
		return &DeclFunc{
			Expr:   f,
			Args:   f.Args - len(args),
			Stored: args,
		}
	}

	return f.Expr.Call(append(f.Stored, args...)...)
}

type Expr struct {
	Func Func
	Args []Func
}

func (f Expr) Call(args ...Func) Func {
	return f.Func.Call(f.Args...)
}

type Chain struct {
	Func Func
	Args []Func
	Prev Func
}

func (f Chain) Call(args ...Func) Func {
	return f.Func.Call(f.Args...).Call(f.Prev.Call())
}

type String string

func (s String) Call(args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

type Number float64

func (n Number) Call(args ...Func) Func {
	// TODO: Use the arguments for something, perhaps.
	return n
}

type External struct {
	Module *Module
	Import ID
	Func   ID
}

func (e External) Call(args ...Func) Func {
	return e.Module.Imports[e.Import].Funcs[e.Func].Call(args...)
}

type Local struct {
	Module *Module
	Func   ID
}

func (local Local) Call(args ...Func) Func {
	return local.Module.Funcs[local.Func].Call(args...)
}

type Compound []Func

func (c Compound) Call(args ...Func) Func {
	var last Func
	for _, f := range c {
		last = f.Call()
	}

	return last
}
