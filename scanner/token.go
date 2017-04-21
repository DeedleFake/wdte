package scanner

type Token struct {
	Line, Col int
	Type      TokenType
	Val       interface{}
}

type TokenType uint

const (
	Invalid TokenType = iota
	Number
	String
	ID
	Keyword
)
