package wdte

import (
	"fmt"

	"github.com/DeedleFake/wdte/ast"
	"github.com/DeedleFake/wdte/scanner"
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
			def := fromFuncDecl(d.Children()[0].(*ast.NTerm), m)
			m.Funcs[def.ID] = def

		default:
			panic(fmt.Errorf("Malformed AST with bad <decl>: %q", dtype))
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

func fromFuncDecl(decl *ast.NTerm, m *Module) *DeclFunc {
	id := ID(decl.Children()[0].(*ast.Term).Tok().Val.(string))
	args := fromArgDecls(flatten(decl.Children()[1].(*ast.NTerm), 1, 0))
	expr := fromExpr(decl.Children()[3].(*ast.NTerm), m, scopeMap(args))

	return &DeclFunc{
		ID:   id,
		Expr: expr,
		Args: len(args),
	}
}

func fromArgDecls(argdecls []ast.Node) []ID {
	ids := make([]ID, 0, len(argdecls))
	for _, arg := range argdecls {
		ids = append(ids, ID(arg.(*ast.Term).Tok().Val.(string)))
	}

	return ids
}

func fromExpr(expr *ast.NTerm, m *Module, scope map[ID]int) Func {
	first := fromSingle(expr.Children()[0].(*ast.NTerm), m, scope)
	in := fromArgs(flatten(expr.Children()[1].(*ast.NTerm), 1, 0), m, scope)

	return fromChain(expr.Children()[2].(*ast.NTerm), &Expr{
		Func: first,
		Args: in,
	}, m, scope)
}

func fromSingle(single *ast.NTerm, m *Module, scope map[ID]int) Func {
	switch s := single.Children()[0].(type) {
	case *ast.Term:
		switch s.Tok().Type {
		case scanner.Number:
			return Number(s.Tok().Val.(float64))
		case scanner.String:
			return String(s.Tok().Val.(string))
		}

	case *ast.NTerm:
		switch s.Name() {
		case "func":
			return fromFunc(s, m, scope)
		case "array":
			return fromArray(s, m, scope)
		case "switch":
			return fromSwitch(s, m, scope)
		case "compound":
			return fromCompound(s, m, scope)
		}
	}

	panic(fmt.Errorf("Malformed AST with bad <single>: %#v", single))
}

func fromFunc(f *ast.NTerm, m *Module, scope map[ID]int) Func {
	id := ID(f.Children()[0].(*ast.Term).Tok().Val.(string))
	sub := fromSubfunc(f.Children()[1].(*ast.NTerm))
	if sub != "" {
		return &External{
			Module: m,
			Import: id,
			Func:   sub,
		}
	}

	if arg, ok := scope[id]; ok {
		return Arg(arg)
	}

	return &Local{
		Module: m,
		Func:   id,
	}
}

func fromSubfunc(subfunc *ast.NTerm) ID {
	if _, ok := subfunc.Children()[0].(*ast.Epsilon); ok {
		return ""
	}

	return ID(subfunc.Children()[1].(*ast.Term).Tok().Val.(string))
}

func fromArray(array *ast.NTerm, m *Module, scope map[ID]int) Func {
	panic("Not implemented.")
}

func fromSwitch(s *ast.NTerm, m *Module, scope map[ID]int) Func {
	check := fromExpr(s.Children()[1].(*ast.NTerm), m, scope)
	switches := fromSwitches(flatten(s.Children()[3].(*ast.NTerm), 4, 0, 2), m, scope)

	return &Switch{
		Check: check,
		Cases: switches,
	}
}

func fromSwitches(switches []ast.Node, m *Module, scope map[ID]int) [][2]Func {
	cases := make([][2]Func, 0, len(switches)/2)
loop:
	for i := 0; i < len(switches); i += 2 {
		switch c := switches[i].(type) {
		case *ast.Term:
			cases = append(cases, [...]Func{
				nil,
				fromExpr(switches[i+1].(*ast.NTerm), m, scope),
			})
			break loop

		case *ast.NTerm:
			cases = append(cases, [...]Func{
				fromExpr(c, m, scope),
				fromExpr(switches[i+1].(*ast.NTerm), m, scope),
			})
		}
	}
	return cases
}

func fromCompound(compound *ast.NTerm, m *Module, scope map[ID]int) Func {
	return Compound(fromExprList(compound.Children()[1].(*ast.NTerm), m, scope))
}

func fromExprList(exprlist *ast.NTerm, m *Module, scope map[ID]int) []Func {
	if _, ok := exprlist.Children()[0].(*ast.Epsilon); ok {
		return nil
	}

	first := fromExpr(exprlist.Children()[0].(*ast.NTerm), m, scope)
	exprs := flatten(exprlist.Children()[1].(*ast.NTerm), 2, 1)

	funcs := make([]Func, 0, 1+len(exprs))
	funcs = append(funcs, first)
	return append(funcs, fromExprs(exprs, m, scope)...)
}

func fromExprs(exprs []ast.Node, m *Module, scope map[ID]int) []Func {
	funcs := make([]Func, 0, len(exprs))
	for _, expr := range exprs {
		funcs = append(funcs, fromExpr(expr.(*ast.NTerm), m, scope))
	}

	return funcs
}

func fromArgs(args []ast.Node, m *Module, scope map[ID]int) []Func {
	singles := make([]Func, 0, len(args))
	for _, arg := range args {
		single := fromSingle(arg.(*ast.NTerm), m, scope)
		singles = append(singles, single)
	}

	return singles
}

func fromChain(chain *ast.NTerm, prev Func, m *Module, scope map[ID]int) Func {
	if _, ok := chain.Children()[0].(*ast.Epsilon); ok {
		return prev
	}

	first := fromSingle(chain.Children()[1].(*ast.NTerm), m, scope)
	in := fromArgs(flatten(chain.Children()[2].(*ast.NTerm), 1, 0), m, scope)

	return fromChain(chain.Children()[3].(*ast.NTerm), &Chain{
		Func: first,
		Args: in,
		Prev: prev,
	}, m, scope)
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

func scopeMap(args []ID) map[ID]int {
	m := make(map[ID]int, len(args))
	for i, arg := range args {
		m[arg] = i
	}

	return m
}

func defaultImporter(from string) (*Module, error) {
	// TODO: This should probably do something else.
	return nil, nil
}
