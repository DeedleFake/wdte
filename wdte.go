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
	Call(scope []Func, args ...Func) Func
}

type GoFunc func(scope []Func, args ...Func) Func

func (f GoFunc) Call(scope []Func, args ...Func) Func {
	return f(scope, args...)
}

type DeclFunc struct {
	Expr   Func
	Args   int
	Stored []Func
}

func (f DeclFunc) Call(scope []Func, args ...Func) Func {
	if len(args) < f.Args {
		return &DeclFunc{
			Expr:   f,
			Args:   f.Args - len(args),
			Stored: args,
		}
	}

	scope = append(f.Stored, args...)
	return f.Expr.Call(scope, scope...)
}

type Expr struct {
	Func Func
	Args []Func
}

func (f Expr) Call(scope []Func, args ...Func) Func {
	return f.Func.Call(scope, f.Args...)
}

type Chain struct {
	Func Func
	Args []Func
	Prev Func
}

func (f Chain) Call(scope []Func, args ...Func) Func {
	return f.Func.Call(scope, f.Args...).Call(scope, f.Prev.Call(scope))
}

type String string

func (s String) Call(scope []Func, args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

type Number float64

func (n Number) Call(scope []Func, args ...Func) Func {
	// TODO: Use the arguments for something, perhaps.
	return n
}

type External struct {
	Module *Module
	Import ID
	Func   ID
}

func (e External) Call(scope []Func, args ...Func) Func {
	return e.Module.Imports[e.Import].Funcs[e.Func].Call(scope, args...)
}

type Local struct {
	Module *Module
	Func   ID
}

func (local Local) Call(scope []Func, args ...Func) Func {
	return local.Module.Funcs[local.Func].Call(scope, args...)
}

type Compound []Func

func (c Compound) Call(scope []Func, args ...Func) Func {
	var last Func
	for _, f := range c {
		last = f.Call(scope)
	}

	return last
}

type Arg int

func (a Arg) Call(scope []Func, args ...Func) Func {
	if int(a) >= len(scope) {
		// TODO: Handle this properly.
		panic("Argument out of scope.")
	}

	return scope[a].Call(scope, args...)
}
