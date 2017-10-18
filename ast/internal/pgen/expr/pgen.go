package pgen

import (
	"github.com/DeedleFake/wdte/ast/internal/pgen"
	"github.com/DeedleFake/wdte/scanner"
)

//go:generate pgen -out table.go ../../../../res/expr.grammar

type (
	Lookup  = pgen.Lookup
	Rule    = pgen.Rule
	Token   = pgen.Token
	Term    = pgen.Term
	NTerm   = pgen.NTerm
	Epsilon = pgen.Epsilon
	EOF     = pgen.EOF
)

func newTerm(t string) Term {
	switch t {
	case "id":
		return Term{
			Type: scanner.ID,
		}

	case "string":
		return Term{
			Type: scanner.String,
		}

	case "number":
		return Term{
			Type: scanner.Number,
		}
	}

	return Term{
		Type:    scanner.Keyword,
		Keyword: t,
	}
}

func newNTerm(nt string) NTerm {
	return NTerm(nt)
}

func newEpsilon() Epsilon {
	return Epsilon{}
}

func newEOF() EOF {
	return EOF{}
}

func newRule(tokens ...Token) (r Rule) {
	return Rule(tokens)
}
