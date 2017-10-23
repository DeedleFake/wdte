package wdte

import (
	"fmt"

	"github.com/DeedleFake/wdte/ast"
	"github.com/DeedleFake/wdte/scanner"
)

func (m *Module) fromScript(script *ast.NTerm, im Importer) (*Module, error) {
	if im == nil {
		im = ImportFunc(defaultImporter)
	}

	return m.fromDecls(flatten(script.Children()[0].(*ast.NTerm), 1, 0), im)
}

func (m *Module) fromDecls(decls []ast.Node, im Importer) (*Module, error) {
	if m.Funcs == nil {
		m.Funcs = make(map[ID]Func)
	}

	for _, d := range decls {
		switch dtype := d.Children()[0].(*ast.NTerm).Name(); dtype {
		case "import":
			id, sub, err := m.fromImport(d.Children()[0].(*ast.NTerm), im)
			if err != nil {
				return nil, err
			}
			m.Funcs[id] = sub

		case "funcdecl":
			def := m.fromFuncDecl(d.Children()[0].(*ast.NTerm))
			m.Funcs[def.ID] = def

		default:
			panic(fmt.Errorf("Malformed AST with bad <decl>: %q", dtype))
		}
	}

	return m, nil
}

func (m *Module) fromImport(i *ast.NTerm, im Importer) (ID, *Module, error) {
	name := i.Children()[0].(*ast.Term).Tok().Val.(string)
	id := ID(i.Children()[2].(*ast.Term).Tok().Val.(string))

	m, err := im.Import(name)
	return id, m, err
}

func (m *Module) fromFuncDecl(decl *ast.NTerm) *DeclFunc {
	mods := m.fromFuncMods(decl.Children()[0].(*ast.NTerm))
	id := ID(decl.Children()[1].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(flatten(decl.Children()[2].(*ast.NTerm), 1, 0))
	expr := m.fromExpr(decl.Children()[4].(*ast.NTerm), args)

	if mods&funcModMemo != 0 {
		expr = &Memo{
			Func: expr,
		}
	}

	return &DeclFunc{
		ID:   id,
		Expr: expr,
		Args: args,
	}
}

func (m *Module) fromFuncMods(funcMods *ast.NTerm) funcMod {
	switch mod := funcMods.Children()[0].(type) {
	case *ast.NTerm:
		return m.fromFuncMod(mod) | m.fromFuncMods(funcMods.Children()[1].(*ast.NTerm))

	case *ast.Epsilon:
		return 0

	default:
		panic(fmt.Errorf("Malformed AST with bad <funcmods>: %T", mod))
	}
}

func (m *Module) fromFuncMod(funcMod *ast.NTerm) funcMod {
	switch mod := funcMod.Children()[0].(*ast.Term).Tok().Val; mod {
	case "memo":
		return funcModMemo

	default:
		panic(fmt.Errorf("Malformed AST with bad <funcmod>: %v", mod))
	}
}

func (m *Module) fromArgDecls(argdecls []ast.Node) []ID {
	ids := make([]ID, 0, len(argdecls))
	for _, arg := range argdecls {
		ids = append(ids, ID(arg.(*ast.Term).Tok().Val.(string)))
	}

	return ids
}

func (m *Module) fromExpr(expr *ast.NTerm, scope []ID) Func {
	first := m.fromSingle(expr.Children()[0].(*ast.NTerm), scope)
	in := m.fromArgs(flatten(expr.Children()[1].(*ast.NTerm), 1, 0), scope)

	slot := m.fromSlot(expr.Children()[2].(*ast.NTerm))
	scope = append(scope, slot)

	return &Expr{
		Func:  first,
		Args:  in,
		Slot:  slot,
		Chain: m.fromChain(expr.Children()[3].(*ast.NTerm), scope),
	}
}

func (m *Module) fromSlot(expr *ast.NTerm) ID {
	if _, ok := expr.Children()[0].(*ast.Epsilon); ok {
		return ""
	}

	return ID(expr.Children()[1].(*ast.Term).Tok().Val.(string))
}

func (m *Module) fromSingle(single *ast.NTerm, scope []ID) Func {
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
			return m.fromFunc(s, scope)
		case "array":
			return m.fromArray(s, scope)
		case "switch":
			return m.fromSwitch(s, scope)
		case "compound":
			return m.fromCompound(s, scope)
		case "lambda":
			return m.fromLambda(s, scope)
		}
	}

	panic(fmt.Errorf("Malformed AST with bad <single>: %#v", single))
}

func (m *Module) fromFunc(f *ast.NTerm, scope []ID) Func {
	tok := f.Children()[0].(*ast.Term).Tok()
	id := ID(tok.Val.(string))

	sub := m.fromSubfunc(f.Children()[1].(*ast.NTerm))
	if sub != "" {
		var im Func

		ok := inScope(scope, id)
		switch ok {
		case true:
			im = Var(id)
		case false:
			im = &Local{
				Module: m,
				Func:   id,
			}
		}

		return &External{
			Module: m,
			Import: im,
			Func:   sub,
		}
	}

	if ok := inScope(scope, id); ok {
		return Var(id)
	}

	return &Local{
		Module: m,
		Func:   id,
	}
}

func (m *Module) fromSubfunc(subfunc *ast.NTerm) ID {
	if _, ok := subfunc.Children()[0].(*ast.Epsilon); ok {
		return ""
	}

	return ID(subfunc.Children()[1].(*ast.Term).Tok().Val.(string))
}

func (m *Module) fromArray(array *ast.NTerm, scope []ID) Func {
	return Array(m.fromExprs(flatten(array.Children()[1].(*ast.NTerm), 2, 0), scope))
}

func (m *Module) fromSwitch(s *ast.NTerm, scope []ID) Func {
	check := m.fromExpr(s.Children()[1].(*ast.NTerm), scope)
	switches := m.fromSwitches(flatten(s.Children()[3].(*ast.NTerm), 4, 0, 2), scope)

	return &Switch{
		Check: check,
		Cases: switches,
	}
}

func (m *Module) fromSwitches(switches []ast.Node, scope []ID) [][2]Func {
	cases := make([][2]Func, 0, len(switches)/2)
loop:
	for i := 0; i < len(switches); i += 2 {
		switch c := switches[i].(type) {
		case *ast.Term:
			cases = append(cases, [...]Func{
				nil,
				m.fromExpr(switches[i+1].(*ast.NTerm), scope),
			})
			break loop

		case *ast.NTerm:
			cases = append(cases, [...]Func{
				m.fromExpr(c, scope),
				m.fromExpr(switches[i+1].(*ast.NTerm), scope),
			})
		}
	}
	return cases
}

func (m *Module) fromCompound(compound *ast.NTerm, scope []ID) Func {
	exprs := flatten(compound.Children()[1].(*ast.NTerm), 2, 0)
	if len(exprs) == 1 {
		return m.fromExpr(exprs[0].(*ast.NTerm), scope)
	}

	return Compound(m.fromExprs(exprs, scope))
}

func (m *Module) fromLambda(lambda *ast.NTerm, scope []ID) Func {
	mods := m.fromFuncMods(lambda.Children()[1].(*ast.NTerm))
	id := ID(lambda.Children()[2].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(flatten(lambda.Children()[3].(*ast.NTerm), 1, 0))
	scope = append(scope, id)
	scope = append(scope, args...)

	expr := m.fromExpr(lambda.Children()[5].(*ast.NTerm), scope)
	if mods&funcModMemo != 0 {
		expr = &Memo{
			Func: expr,
		}
	}

	return &Lambda{
		ID:   id,
		Expr: expr,
		Args: args,
	}
}

func (m *Module) fromExprs(exprs []ast.Node, scope []ID) []Func {
	funcs := make([]Func, 0, len(exprs))
	for _, expr := range exprs {
		funcs = append(funcs, m.fromExpr(expr.(*ast.NTerm), scope))
	}

	return funcs
}

func (m *Module) fromArgs(args []ast.Node, scope []ID) []Func {
	singles := make([]Func, 0, len(args))
	for _, arg := range args {
		single := m.fromSingle(arg.(*ast.NTerm), scope)
		singles = append(singles, single)
	}

	return singles
}

func (m *Module) fromChain(chain *ast.NTerm, scope []ID) Func {
	if _, ok := chain.Children()[0].(*ast.Epsilon); ok {
		return &EndChain{}
	}

	// TODO: Make this properly recursive with m.fromExpr().
	expr := chain.Children()[1].(*ast.NTerm)

	first := m.fromSingle(expr.Children()[0].(*ast.NTerm), scope)
	in := m.fromArgs(flatten(expr.Children()[1].(*ast.NTerm), 1, 0), scope)

	slot := m.fromSlot(expr.Children()[2].(*ast.NTerm))
	scope = append(scope, slot)

	switch t := chain.Children()[0].(*ast.Term).Tok().Val.(string); t {
	case "->":
		return &Chain{
			Func:  first,
			Args:  in,
			Slot:  slot,
			Chain: m.fromChain(expr.Children()[3].(*ast.NTerm), scope),
		}

	case "--":
		return &IgnoredChain{
			Func:  first,
			Args:  in,
			Slot:  slot,
			Chain: m.fromChain(expr.Children()[3].(*ast.NTerm), scope),
		}

	default:
		panic(fmt.Errorf("Malformed AST with unexpected chain type: %q", t))
	}
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

func inScope(scope []ID, id ID) bool {
	for _, s := range scope {
		if s == id {
			return true
		}
	}

	return false
}

func defaultImporter(from string) (*Module, error) {
	// TODO: This should probably do something else.
	return nil, nil
}

type funcMod uint

const (
	funcModMemo funcMod = 1 << iota
)
