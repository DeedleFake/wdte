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

	if m.Funcs == nil {
		m.Funcs = make(map[ID]Func)
	}

	return m.fromDecls(script.Children()[0].(*ast.NTerm), im)
}

func (m *Module) fromDecls(decls *ast.NTerm, im Importer) (*Module, error) {
	switch decl := decls.Children()[0].(type) {
	case *ast.NTerm:
		switch dtype := decl.Children()[0].(*ast.NTerm).Name(); dtype {
		case "import":
			id, sub, err := m.fromImport(decl.Children()[0].(*ast.NTerm), im)
			if err != nil {
				return nil, err
			}
			m.Funcs[id] = sub

		case "funcdecl":
			id, def := m.fromFuncDecl(decl.Children()[0].(*ast.NTerm))
			m.Funcs[id] = def

		default:
			panic(fmt.Errorf("Malformed AST with bad <decl>: %q", dtype))
		}

		return m.fromDecls(decls.Children()[1].(*ast.NTerm), im)

	case *ast.Epsilon:
		return m, nil

	default:
		panic(fmt.Errorf("Malformed AST with bad <decls>: %T", decl))
	}
}

func (m *Module) fromImport(i *ast.NTerm, im Importer) (ID, *Module, error) {
	name := i.Children()[0].(*ast.Term).Tok().Val.(string)
	id := ID(i.Children()[2].(*ast.Term).Tok().Val.(string))

	m, err := im.Import(name)
	return id, m, err
}

func (m *Module) fromFuncDecl(decl *ast.NTerm) (id ID, f Func) {
	mods := m.fromFuncMods(decl.Children()[0].(*ast.NTerm))
	id = ID(decl.Children()[1].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(decl.Children()[2].(*ast.NTerm), nil)
	expr := m.fromExpr(decl.Children()[4].(*ast.NTerm), args)

	if mods&funcModMemo != 0 {
		expr = &Memo{
			Func: expr,
		}
	}

	return id, &DeclFunc{
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

func (m *Module) fromArgDecls(argdecls *ast.NTerm, ids []ID) []ID {
	switch arg := argdecls.Children()[0].(type) {
	case *ast.Term:
		ids = append(ids, ID(arg.Tok().Val.(string)))
		return m.fromArgDecls(argdecls.Children()[1].(*ast.NTerm), ids)

	case *ast.Epsilon:
		return ids

	default:
		panic(fmt.Errorf("Malformed AST with bad <argdecls>: %T", arg))
	}
}

func (m *Module) fromExpr(expr *ast.NTerm, scope []ID) Func {
	first := m.fromSingle(expr.Children()[0].(*ast.NTerm), scope)
	in := m.fromArgs(expr.Children()[1].(*ast.NTerm), scope, nil)

	slot := m.fromSlot(expr.Children()[2].(*ast.NTerm))
	scope = append(scope, slot)

	return &Expr{
		Func:  first,
		Args:  in,
		Slot:  slot,
		Chain: m.fromChain(expr.Children()[3].(*ast.NTerm), scope),
	}
}

func (m *Module) fromLetExpr(expr *ast.NTerm, scope []ID) Func {
	mods := m.fromFuncMods(expr.Children()[1].(*ast.NTerm))
	id := ID(expr.Children()[2].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(expr.Children()[3].(*ast.NTerm), nil)
	scope = append(scope, id)
	scope = append(scope, args...)
	inner := m.fromExpr(expr.Children()[5].(*ast.NTerm), scope)

	if mods&funcModMemo != 0 {
		inner = &Memo{
			Func: inner,
		}
	}

	lambda := &Lambda{
		ID:   id,
		Expr: inner,
		Args: args,
	}

	return &Let{
		ID:   id,
		Expr: lambda,
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
	return Array(m.fromExprs(array.Children()[1].(*ast.NTerm), scope, nil))
}

func (m *Module) fromSwitch(s *ast.NTerm, scope []ID) Func {
	check := m.fromExpr(s.Children()[1].(*ast.NTerm), scope)
	switches := m.fromSwitches(s.Children()[3].(*ast.NTerm), scope, nil)

	return &Switch{
		Check: check,
		Cases: switches,
	}
}

func (m *Module) fromSwitches(switches *ast.NTerm, scope []ID, cases [][2]Func) [][2]Func {
	switch sw := switches.Children()[0].(type) {
	case *ast.Term:
		cases = append(cases, [...]Func{
			nil,
			m.fromExpr(switches.Children()[2].(*ast.NTerm), scope),
		})
		return cases

	case *ast.NTerm:
		cases = append(cases, [...]Func{
			m.fromExpr(sw, scope),
			m.fromExpr(switches.Children()[2].(*ast.NTerm), scope),
		})
		return m.fromSwitches(switches.Children()[4].(*ast.NTerm), scope, cases)

	case *ast.Epsilon:
		return cases

	default:
		panic(fmt.Errorf("Malformed AST with bad <switches>: %T", sw))
	}
}

func (m *Module) fromCompound(compound *ast.NTerm, scope []ID) Func {
	return Compound(m.fromExprs(compound.Children()[1].(*ast.NTerm), scope, nil))
}

func (m *Module) fromLambda(lambda *ast.NTerm, scope []ID) (f Func) {
	mods := m.fromFuncMods(lambda.Children()[1].(*ast.NTerm))
	id := ID(lambda.Children()[2].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(lambda.Children()[3].(*ast.NTerm), nil)
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

func (m *Module) fromExprs(exprs *ast.NTerm, scope []ID, funcs []Func) []Func {
	switch expr := exprs.Children()[0].(type) {
	case *ast.NTerm:
		switch expr.Name() {
		case "expr":
			funcs = append(funcs, m.fromExpr(expr, scope))
		case "letexpr":
			let := m.fromLetExpr(expr, scope)
			funcs = append(funcs, let)
			scope = append(scope, let.(*Let).ID)
		}
		return m.fromExprs(exprs.Children()[2].(*ast.NTerm), scope, funcs)

	case *ast.Epsilon:
		return funcs

	default:
		panic(fmt.Errorf("Malformed AST with bad <exprs>: %T", expr))
	}
}

func (m *Module) fromArgs(args *ast.NTerm, scope []ID, funcs []Func) []Func {
	switch arg := args.Children()[0].(type) {
	case *ast.NTerm:
		funcs = append(funcs, m.fromSingle(arg, scope))
		return m.fromArgs(args.Children()[1].(*ast.NTerm), scope, funcs)

	case *ast.Epsilon:
		return funcs

	default:
		panic(fmt.Errorf("Malformed AST with bad <args>: %T", arg))
	}
}

func (m *Module) fromChain(chain *ast.NTerm, scope []ID) Func {
	if _, ok := chain.Children()[0].(*ast.Epsilon); ok {
		return &EndChain{}
	}

	// TODO: Make this properly recursive with m.fromExpr().
	expr := chain.Children()[1].(*ast.NTerm)

	first := m.fromSingle(expr.Children()[0].(*ast.NTerm), scope)
	in := m.fromArgs(expr.Children()[1].(*ast.NTerm), scope, nil)

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
