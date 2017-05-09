package wdte

import (
	"fmt"

	"github.com/DeedleFake/wdte/ast"
)

func fromScript(script *ast.NTerm, im Importer) (*Module, error) {
	if im == nil {
		im = ImportFunc(defaultImporter)
	}

	return fromDecls(flatten(script.Children()[0].(*ast.NTerm), 1, 0), im)
}

func fromDecls(decls []ast.Node, im Importer) (*Module, error) {
	m := &Module{
		Imports: make(map[ID]*Module),
		Funcs:   make(map[ID]Func),
	}

	for _, d := range decls {
		switch dtype := d.Children()[0].(*ast.NTerm).Name(); dtype {
		case "import":
			id, sub, err := fromImport(d.Children()[0].(*ast.NTerm), im)
			if err != nil {
				return nil, err
			}
			m.Imports[ID(id)] = sub

		case "funcdecl":
			panic("Not implemented.")

		default:
			return nil, fmt.Errorf("Unexpected decl type: %q", dtype)
		}
	}

	return m, nil
}

func fromImport(i *ast.NTerm, im Importer) (string, *Module, error) {
	name := i.Children()[0].(*ast.Term).Tok().Val.(string)
	id := i.Children()[2].(*ast.Term).Tok().Val.(string)

	m, err := im.Import(name)
	return id, m, err
}

func flatten(top *ast.NTerm, rec int, get ...int) []ast.Node {
	if _, ok := top.Children()[0].(*ast.Epsilon); ok {
		return []ast.Node{}
	}

	c := make([]ast.Node, 0, len(get))
	for _, i := range get {
		c = append(c, top.Children()[i])
	}

	return append(c, flatten(top.Children()[rec].(*ast.NTerm), rec, get...)...)
}

func defaultImporter(from string) (*Module, error) {
	// TODO: This should probably do something else.
	return nil, nil
}
