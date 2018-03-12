package pgen

import "github.com/DeedleFake/wdte/scanner"

//go:generate pgen -out table.go ../../../res/grammar.ebnf

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

type Token interface{}

type Term struct {
	Type    scanner.TokenType
	Keyword string
}

func (t Term) String() string {
	if t.Type == scanner.Keyword {
		return t.Keyword
	}

	return t.Type.String()
}

type NTerm string

type Epsilon struct{}

func (e Epsilon) String() string {
	return "ε"
}

type EOF struct{}

func (e EOF) String() string {
	return "Ω"
}

type Rule []Token

type Lookup struct {
	Term  Token
	NTerm NTerm
}
