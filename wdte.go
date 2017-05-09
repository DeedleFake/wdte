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
	Equals(other Func) bool
}

type GoFunc func(scope []Func, args ...Func) Func

func (f GoFunc) Call(scope []Func, args ...Func) Func {
	return f(scope, args...)
}

func (f GoFunc) Equals(other Func) bool {
	panic("Not implemented.")
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

func (f DeclFunc) Equals(other Func) bool {
	panic("Not implemented.")
}

type Expr struct {
	Func Func
	Args []Func
}

func (f Expr) Call(scope []Func, args ...Func) Func {
	return f.Func.Call(scope, f.Args...)
}

func (f Expr) Equals(other Func) bool {
	panic("Not implemented.")
}

type Chain struct {
	Func Func
	Args []Func
	Prev Func
}

func (f Chain) Call(scope []Func, args ...Func) Func {
	return f.Func.Call(scope, f.Args...).Call(scope, f.Prev.Call(scope))
}

func (f Chain) Equals(other Func) bool {
	panic("Not implemented.")
}

type String string

func (s String) Call(scope []Func, args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

func (s String) Equals(other Func) bool {
	o, ok := other.(String)
	return ok && (s == o)
}

type Number float64

func (n Number) Call(scope []Func, args ...Func) Func {
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

func (e External) Call(scope []Func, args ...Func) Func {
	return e.Module.Imports[e.Import].Funcs[e.Func].Call(scope, args...)
}

func (e External) Equals(other Func) bool {
	o, ok := other.(External)
	return ok && (e.Import == o.Import) && (e.Func == o.Func)
}

type Local struct {
	Module *Module
	Func   ID
}

func (local Local) Call(scope []Func, args ...Func) Func {
	return local.Module.Funcs[local.Func].Call(scope, args...)
}

func (local Local) Equals(other Func) bool {
	o, ok := other.(Local)
	return ok && (local.Func == o.Func)
}

type Compound []Func

func (c Compound) Call(scope []Func, args ...Func) Func {
	var last Func
	for _, f := range c {
		last = f.Call(scope)
	}

	return last
}

func (c Compound) Equals(other Func) bool {
	panic("Not implemented.")
}

type Arg int

func (a Arg) Call(scope []Func, args ...Func) Func {
	if int(a) >= len(scope) {
		// TODO: Handle this properly.
		panic("Argument out of scope.")
	}

	return scope[a].Call(scope, args...)
}

func (a Arg) Equals(other Func) bool {
	panic("Not implemented.")
}

type Switch struct {
	Check Func
	Cases [][2]Func
}

func (s Switch) Call(scope []Func, args ...Func) Func {
	check := s.Check.Call(scope)
	for _, c := range s.Cases {
		if (c[0] == nil) || (check.Equals(c[0].Call(scope))) {
			return c[1].Call(scope)
		}
	}

	return nil
}

func (s Switch) Equals(other Func) bool {
	panic("Not implemented.")
}
