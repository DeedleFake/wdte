package scanner

// A Token is a usable element parsed from a string.
type Token struct {
	Line, Col int
	Val       TokenValue
}

// A TokenValue is the value of a token.
type TokenValue interface {
	token()
}

type Number float64

func (Number) token() {}

type String string

func (String) token() {}

type ID string

func (ID) token() {}

type Keyword string

func (Keyword) token() {}

type Macro struct {
	Name  string
	Input string
}

func (Macro) token() {}

type EOF struct{}

func (EOF) token() {}
