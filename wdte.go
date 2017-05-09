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
	ast, err := ast.Parse(r)
	if err != nil {
		return nil, err
	}

	return FromAST(ast, im)
}

func FromAST(ast ast.Node, im Importer) (*Module, error) {
	panic("Not implemented.")
}

type Importer interface {
	Import(from string) (*Module, error)
}

type ID string

type Func interface {
	Call(args ...Expr) Expr
}

type Expr interface {
	Eval(scope map[ID]Expr) Value
}

type Value interface{}
