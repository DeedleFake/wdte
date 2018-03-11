package wdte

import (
	"fmt"

	"github.com/DeedleFake/wdte/ast"
	"github.com/DeedleFake/wdte/scanner"
)

type translator struct {
	im Importer
}

func (m *translator) fromScript(script *ast.NTerm) (c Compound, err error) {
	defer func() {
		switch e := recover().(type) {
		case error:
			err = e

		case nil:

		default:
			panic(e)
		}
	}()

	return Compound(m.fromExprs(script.Children()[0].(*ast.NTerm), nil)), nil
}

func (m *translator) fromFuncMods(funcMods *ast.NTerm) funcMod {
	switch mod := funcMods.Children()[0].(type) {
	case *ast.NTerm:
		return m.fromFuncMod(mod) | m.fromFuncMods(funcMods.Children()[1].(*ast.NTerm))

	case *ast.Epsilon:
		return 0

	default:
		panic(fmt.Errorf("Malformed AST with bad <funcmods>: %T", mod))
	}
}

func (m *translator) fromFuncMod(funcMod *ast.NTerm) funcMod {
	switch mod := funcMod.Children()[0].(*ast.Term).Tok().Val; mod {
	case "memo":
		return funcModMemo

	default:
		panic(fmt.Errorf("Malformed AST with bad <funcmod>: %v", mod))
	}
}

func (m *translator) fromArgDecls(argdecls *ast.NTerm, ids []ID) []ID {
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

func (m *translator) fromExpr(expr *ast.NTerm) Func {
	first := m.fromSingle(expr.Children()[0].(*ast.NTerm))
	in := m.fromArgs(expr.Children()[1].(*ast.NTerm), nil)

	slot := m.fromSlot(expr.Children()[2].(*ast.NTerm))

	return &Expr{
		Func:  first,
		Args:  in,
		Slot:  slot,
		Chain: m.fromChain(expr.Children()[3].(*ast.NTerm)),
	}
}

func (m *translator) fromLetExpr(expr *ast.NTerm) Func {
	mods := m.fromFuncMods(expr.Children()[1].(*ast.NTerm))
	id := ID(expr.Children()[2].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(expr.Children()[3].(*ast.NTerm), nil)
	inner := m.fromExpr(expr.Children()[5].(*ast.NTerm))

	var ret Func
	ret = &Lambda{
		ID:   id,
		Expr: inner,
		Args: args,
	}

	if mods&funcModMemo != 0 {
		ret = &Memo{
			Func: ret,
		}
	}

	return &Let{
		ID:   id,
		Expr: ret,
	}
}

func (m *translator) fromSlot(expr *ast.NTerm) ID {
	if _, ok := expr.Children()[0].(*ast.Epsilon); ok {
		return ""
	}

	return ID(expr.Children()[1].(*ast.Term).Tok().Val.(string))
}

func (m *translator) fromSingle(single *ast.NTerm) Func {
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
			return m.fromFunc(s)
		case "array":
			return m.fromArray(s)
		case "switch":
			return m.fromSwitch(s)
		case "compound":
			return m.fromCompound(s)
		case "lambda":
			return m.fromLambda(s)
		case "import":
			return m.fromImport(s)
		}
	}

	panic(fmt.Errorf("Malformed AST with bad <single>: %#v", single))
}

func (m *translator) fromFunc(f *ast.NTerm) Func {
	id := ID(f.Children()[0].(*ast.Term).Tok().Val.(string))

	sub := m.fromSubFunc(f.Children()[1].(*ast.NTerm))
	if sub == "" {
		return Var(id)
	}

	return Sub{
		Module: Var(id),
		Func:   sub,
	}
}

func (m *translator) fromSubFunc(subfunc *ast.NTerm) ID {
	if _, ok := subfunc.Children()[0].(*ast.Epsilon); ok {
		return ""
	}

	return ID(subfunc.Children()[1].(*ast.Term).Tok().Val.(string))
}

func (m *translator) fromArray(array *ast.NTerm) Func {
	return Array(m.fromExprs(array.Children()[1].(*ast.NTerm), nil))
}

func (m *translator) fromSwitch(s *ast.NTerm) Func {
	check := m.fromExpr(s.Children()[1].(*ast.NTerm))
	switches := m.fromSwitches(s.Children()[3].(*ast.NTerm), nil)

	return &Switch{
		Check: check,
		Cases: switches,
	}
}

func (m *translator) fromSwitches(switches *ast.NTerm, cases [][2]Func) [][2]Func {
	switch sw := switches.Children()[0].(type) {
	case *ast.Term:
		cases = append(cases, [...]Func{
			nil,
			m.fromExpr(switches.Children()[2].(*ast.NTerm)),
		})
		return cases

	case *ast.NTerm:
		cases = append(cases, [...]Func{
			m.fromExpr(sw),
			m.fromExpr(switches.Children()[2].(*ast.NTerm)),
		})
		return m.fromSwitches(switches.Children()[4].(*ast.NTerm), cases)

	case *ast.Epsilon:
		return cases

	default:
		panic(fmt.Errorf("Malformed AST with bad <switches>: %T", sw))
	}
}

func (m *translator) fromCompound(compound *ast.NTerm) Func {
	return Compound(m.fromExprs(compound.Children()[1].(*ast.NTerm), nil))
}

func (m *translator) fromLambda(lambda *ast.NTerm) (f Func) {
	mods := m.fromFuncMods(lambda.Children()[1].(*ast.NTerm))
	id := ID(lambda.Children()[2].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(lambda.Children()[3].(*ast.NTerm), nil)
	expr := m.fromExpr(lambda.Children()[5].(*ast.NTerm))

	var ret Func
	ret = &Lambda{
		ID:   id,
		Expr: expr,
		Args: args,
	}

	if mods&funcModMemo != 0 {
		ret = &Memo{
			Func: ret,
		}
	}

	return ret
}

func (m *translator) fromImport(im *ast.NTerm) Func {
	s, err := m.im.Import(im.Children()[1].(*ast.Term).Tok().Val.(string))
	if err != nil {
		panic(err)
	}

	return s
}

func (m *translator) fromExprs(exprs *ast.NTerm, funcs []Func) []Func {
	switch expr := exprs.Children()[0].(type) {
	case *ast.NTerm:
		switch expr.Name() {
		case "expr":
			funcs = append(funcs, m.fromExpr(expr))
		case "letexpr":
			let := m.fromLetExpr(expr)
			funcs = append(funcs, let)
		}
		return m.fromExprs(exprs.Children()[2].(*ast.NTerm), funcs)

	case *ast.Epsilon:
		return funcs

	default:
		panic(fmt.Errorf("Malformed AST with bad <exprs>: %T", expr))
	}
}

func (m *translator) fromArgs(args *ast.NTerm, funcs []Func) []Func {
	switch arg := args.Children()[0].(type) {
	case *ast.NTerm:
		funcs = append(funcs, m.fromSingle(arg))
		return m.fromArgs(args.Children()[1].(*ast.NTerm), funcs)

	case *ast.Epsilon:
		return funcs

	default:
		panic(fmt.Errorf("Malformed AST with bad <args>: %T", arg))
	}
}

func (m *translator) fromChain(chain *ast.NTerm) Func {
	if _, ok := chain.Children()[0].(*ast.Epsilon); ok {
		return &EndChain{}
	}

	// TODO: Make this properly recursive with m.fromExpr().
	expr := chain.Children()[1].(*ast.NTerm)

	first := m.fromSingle(expr.Children()[0].(*ast.NTerm))
	in := m.fromArgs(expr.Children()[1].(*ast.NTerm), nil)

	slot := m.fromSlot(expr.Children()[2].(*ast.NTerm))

	switch t := chain.Children()[0].(*ast.Term).Tok().Val.(string); t {
	case "->":
		return &Chain{
			Func:  first,
			Args:  in,
			Slot:  slot,
			Chain: m.fromChain(expr.Children()[3].(*ast.NTerm)),
		}

	case "--":
		return &IgnoredChain{
			Func:  first,
			Args:  in,
			Slot:  slot,
			Chain: m.fromChain(expr.Children()[3].(*ast.NTerm)),
		}

	default:
		panic(fmt.Errorf("Malformed AST with unexpected chain type: %q", t))
	}
}

type funcMod uint

const (
	funcModMemo funcMod = 1 << iota
)
