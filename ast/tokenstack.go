package ast

import "github.com/DeedleFake/wdte/ast/internal/pgen"

type tokenStack []pgen.Token

func (ts *tokenStack) Pop() pgen.Token {
	t := (*ts)[len(*ts)-1]
	*ts = (*ts)[:len(*ts)-1]
	return t
}

func (ts *tokenStack) Push(t pgen.Token) {
	*ts = append(*ts, t)
}

func (ts *tokenStack) PushRule(r pgen.Rule) {
	ts.Push(nil)
	for i := len(r) - 1; i >= 0; i-- {
		ts.Push(r[i])
	}
}
