package pgen

import "github.com/DeedleFake/wdte/scanner"

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
