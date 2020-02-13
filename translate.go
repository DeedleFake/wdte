package wdte

import (
	"fmt"
	"runtime"

	"github.com/DeedleFake/wdte/ast"
	"github.com/DeedleFake/wdte/scanner"
)

type translator struct {
	im Importer
}

func (m *translator) fromScript(script *ast.NTerm) (c Compound, err error) {
	defer func() {
		switch e := recover().(type) {
		case runtime.Error:
			panic(e)

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

func (m *translator) fromExpr(expr *ast.NTerm, flags uint, chain Chain) (r Func) {
	first := m.fromSingle(expr.Children()[0].(*ast.NTerm))
	in := m.fromArgs(expr.Children()[1].(*ast.NTerm), nil)
	slots, assignFunc := m.fromSlot(expr.Children()[3].(*ast.NTerm))

	r = &FuncCall{
		Func: first,
		Args: in,
	}
	r = m.fromSwitch(expr.Children()[2].(*ast.NTerm), r)

	piece := &ChainPiece{
		Expr: r,

		Flags:      flags,
		Slots:      slots,
		AssignFunc: assignFunc,
	}

	fc := m.fromChain(expr.Children()[4].(*ast.NTerm), append(chain, piece))
	if chain, ok := fc.(Chain); ok && (len(chain) == 1) {
		return r
	}
	return fc
}

func (m *translator) fromLetExpr(expr *ast.NTerm) Func {
	assign := expr.Children()[1].(*ast.NTerm)

	switch first := assign.Children()[0].(type) {
	case *ast.NTerm:
		mods := m.fromFuncMods(first)
		id := ID(assign.Children()[1].(*ast.Term).Tok().Val.(string))
		args := m.fromArgDecls(assign.Children()[2].(*ast.NTerm), nil)
		inner := m.fromExpr(assign.Children()[4].(*ast.NTerm), 0, nil)

		if mods&funcModMemo != 0 {
			inner = &Memo{
				Func: inner,
				Args: args,
			}
		}

		var right Func
		switch len(args) {
		case 0:
			right = inner

		default:
			right = &Lambda{
				ID:   id,
				Expr: inner,
				Args: args,
			}
		}

		return &Assigner{
			AssignFunc: AssignSimple,

			IDs:  []ID{id},
			Expr: right,
		}

	case *ast.Term:
		return &Assigner{
			AssignFunc: AssignPattern,

			IDs: m.fromArgDecls(
				assign.Children()[2].(*ast.NTerm),
				[]ID{ID(assign.Children()[1].(*ast.Term).Tok().Val.(string))},
			),
			Expr: m.fromExpr(assign.Children()[6].(*ast.NTerm), 0, nil),
		}
	}

	panic(fmt.Errorf("Malformed AST with bad <assign>: %#v", assign))
}

func (m *translator) fromSlot(expr *ast.NTerm) ([]ID, AssignFunc) {
	if _, ok := expr.Children()[0].(*ast.Epsilon); ok {
		return nil, func(frame Frame, scope *Scope, ids []ID, val Func) (*Scope, Func) {
			return scope, val
		}
	}

	assign := expr.Children()[1].(*ast.NTerm)

	first := assign.Children()[0].(*ast.Term)
	if first.Tok().Type == scanner.ID {
		return []ID{ID(first.Tok().Val.(string))}, AssignSimple
	}

	return m.fromArgDecls(
		assign.Children()[2].(*ast.NTerm),
		[]ID{ID(assign.Children()[1].(*ast.Term).Tok().Val.(string))},
	), AssignPattern
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
		case "array":
			return m.fromArray(s)

		case "lambda":
			return m.fromLambda(s)

		case "import":
			return m.fromImport(s)

		case "subbable":
			sub := m.fromSubbable(s, nil)
			if len(sub) == 1 {
				return sub[0]
			}
			return sub
		}
	}

	panic(fmt.Errorf("Malformed AST with bad <single>: %#v", single))
}

func (m *translator) fromSubbable(subbable *ast.NTerm, acc Sub) Sub {
	var found bool
	switch s := subbable.Children()[0].(type) {
	case *ast.Term:
		switch s.Tok().Type {
		case scanner.ID:
			acc = append(acc, Var(s.Tok().Val.(string)))
			found = true
		}

	case *ast.NTerm:
		switch s.Name() {
		case "compound":
			acc = append(acc, m.fromCompound(s))
			found = true
		}
	}

	if found {
		return m.fromSub(subbable.Children()[1].(*ast.NTerm), acc)
	}

	panic(fmt.Errorf("Malformed AST with bad <subbable>: %#v", subbable))
}

func (m *translator) fromSub(sub *ast.NTerm, acc Sub) Sub {
	if _, ok := sub.Children()[0].(*ast.Epsilon); ok {
		return acc
	}

	return m.fromSubbable(sub.Children()[1].(*ast.NTerm), acc)
}

func (m *translator) fromArray(array *ast.NTerm) Func {
	aexprs := array.Children()[1].(*ast.NTerm)
	if _, ok := aexprs.Children()[0].(*ast.Term); ok {
		return Array{}
	}

	return Array(m.fromExprs(aexprs.Children()[0].(*ast.NTerm), nil))
}

func (m *translator) fromSwitch(s *ast.NTerm, check Func) Func {
	if _, ok := s.Children()[0].(*ast.Epsilon); ok {
		return check
	}

	switches := m.fromSwitches(s.Children()[1].(*ast.NTerm), nil)

	return &Switch{
		Check: check,
		Cases: switches,
	}
}

func (m *translator) fromSwitches(switches *ast.NTerm, cases [][2]Func) [][2]Func {
	switch sw := switches.Children()[0].(type) {
	case *ast.NTerm:
		cases = append(cases, [...]Func{
			m.fromExpr(sw, 0, nil),
			m.fromExpr(switches.Children()[2].(*ast.NTerm), 0, nil),
		})
		return m.fromSwitches(switches.Children()[4].(*ast.NTerm), cases)

	case *ast.Epsilon:
		return cases

	default:
		panic(fmt.Errorf("Malformed AST with bad <switches>: %T", sw))
	}
}

func (m *translator) fromCompound(compound *ast.NTerm) Func {
	mode := compound.Children()[0].(*ast.Term).Tok().Val.(string)
	c := Compound(m.fromExprs(compound.Children()[1].(*ast.NTerm), nil))
	if (len(c) == 1) && (mode != "(|") {
		if _, ok := c[0].(*Assigner); !ok {
			return c[0]
		}
	}

	r := Func(c)
	if mode == "(|" {
		r = Collector{Compound: c}
	}

	return r
}

func (m *translator) fromLambda(lambda *ast.NTerm) (f Func) {
	mods := m.fromFuncMods(lambda.Children()[1].(*ast.NTerm))
	id := ID(lambda.Children()[2].(*ast.Term).Tok().Val.(string))
	args := m.fromArgDecls(lambda.Children()[3].(*ast.NTerm), nil)
	expr := Compound(m.fromExprs(lambda.Children()[5].(*ast.NTerm), nil))

	inner := Func(expr)
	if len(expr) == 1 {
		if _, ok := expr[0].(*Assigner); !ok {
			inner = expr[0]
		}
	}

	if mods&funcModMemo != 0 {
		inner = &Memo{
			Func: inner,
			Args: args,
		}
	}

	return &Lambda{
		ID:   id,
		Expr: inner,
		Args: args,
	}
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
			funcs = append(funcs, m.fromExpr(expr, 0, nil))
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

func (m *translator) fromChain(chain *ast.NTerm, acc Chain) Func {
	if _, ok := chain.Children()[0].(*ast.Epsilon); ok {
		return acc
	}

	oper := chain.Children()[0].(*ast.Term).Tok().Val

	var flags uint
	if oper == "--" {
		flags |= IgnoredChain
	}
	if oper == "-|" {
		flags |= ErrorChain
	}

	return m.fromExpr(chain.Children()[1].(*ast.NTerm), flags, acc)
}

type funcMod uint

const (
	funcModMemo funcMod = 1 << iota
)
