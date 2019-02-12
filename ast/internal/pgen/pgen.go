package pgen

import "fmt"

//go:generate pgen -out table.go ../../../res/grammar.ebnf

const (
	ID uint = iota
	String
	Number
	Keyword
)

func newTerm(t string) Term {
	switch t {
	case "id":
		return Term{Type: ID}

	case "string":
		return Term{Type: String}

	case "number":
		return Term{Type: Number}
	}

	return Term{Type: Keyword, Keyword: t}
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
	Type    uint
	Keyword string
}

func (t Term) String() string {
	switch t.Type {
	case ID:
		return "id"
	case String:
		return "string"
	case Number:
		return "number"
	case Keyword:
		return t.Keyword
	}

	panic(fmt.Errorf("Unknown terminal type: %v", t.Type))
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
