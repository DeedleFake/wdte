package pgen

import "github.com/DeedleFake/wdte/scanner"

//go:generate pgen -out table.go ../../../res/grammar.ebnf

type Term struct {
	Type    scanner.TokenType
	Keyword string
}

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

type NTerm string

func newNTerm(nt string) NTerm {
	return NTerm(nt)
}

type Epsilon struct{}

func newEpsilon() Epsilon {
	return Epsilon{}
}

func newRule(tokens ...Token) (r Rule) {
	return Rule(tokens)
}
