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
	Call(args ...Func) Value
}

type GoFunc func(args ...Func) Func

func (f GoFunc) Call(args ...Func) Func {
	return f(args...)
}

type Value interface{}
