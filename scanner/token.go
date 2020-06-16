package scanner

import (
	"fmt"
)

// A Token is a usable element parsed from a string.
type Token struct {
	Line, Col int
	Type      TokenType
	Val       interface{}
}

// TokenType is the type of a token.
type TokenType uint

const (
	Invalid TokenType = iota
	Number
	String
	ID
	Keyword
	Macro
	EOF
)

func (t TokenType) String() string {
	switch t {
	case Invalid:
		return "invalid"
	case Number:
		return "number"
	case String:
		return "string"
	case ID:
		return "id"
	case Keyword:
		return "keyword"
	case Macro:
		return "macro"
	case EOF:
		return "EOF"
	}

	panic(fmt.Errorf("Invalid token type: %v", uint(t)))
}
