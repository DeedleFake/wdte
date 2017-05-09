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
			m.Imports[id] = sub

		case "funcdecl":
			id, def, err := fromFuncDecl(d.Children()[0].(*ast.NTerm), m)
			if err != nil {
				return nil, err
			}
			m.Funcs[id] = def

		default:
			return nil, fmt.Errorf("Unexpected decl type: %q", dtype)
		}
	}

	return m, nil
}

func fromImport(i *ast.NTerm, im Importer) (ID, *Module, error) {
	name := i.Children()[0].(*ast.Term).Tok().Val.(string)
	id := ID(i.Children()[2].(*ast.Term).Tok().Val.(string))

	m, err := im.Import(name)
	return id, m, err
}

func fromFuncDecl(decl *ast.NTerm, m *Module) (ID, Func, error) {
	id := ID(decl.Children()[0].(*ast.Term).Tok().Val.(string))
	args := fromArgDecls(flatten(decl.Children()[1].(*ast.NTerm), 1, 0))

	expr, err := fromExpr(decl.Children()[3].(*ast.NTerm), m, args)
	if err != nil {
		return "", nil, err
	}

	return id, &DeclFunc{
		Expr: expr,
		Args: len(args),
	}, nil
}

func fromExpr(expr *ast.NTerm, m *Module, args []ID) (Func, error) {
	panic("Not implemented.")
}

func fromArgDecls(argdecls []ast.Node) []ID {
	ids := make([]ID, 0, len(argdecls))
	for _, arg := range argdecls {
		ids = append(ids, ID(arg.(*ast.Term).Tok().Val.(string)))
	}

	return ids
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
