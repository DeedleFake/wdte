package scanner

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
	StmtEnd
)
