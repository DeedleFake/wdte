package scanner

type Token interface {
	setPos(line, col int)
	Pos() (line, col int)
}

type commonToken struct {
	line int
	col  int
}

func (t *commonToken) setPos(line, col int) {
	t.line, t.col = line, col
}

func (t commonToken) Pos() (line, col int) {
	return t.line, t.col
}

type Number struct {
	commonToken

	Val float64
}

type String struct {
	commonToken

	Val string
}

type Nil struct {
	commonToken
}

type ID struct {
	commonToken

	Val string
}

type Keyword struct {
	commonToken

	Val string
}
