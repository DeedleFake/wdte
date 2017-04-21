package scanner

type Token interface {
	setPos(line, col int)

	Line() int
	Col() int
}

type commonToken struct {
	line int
	col  int
}

func (t *commonToken) setPos(line, col int) {
	t.line, t.col = line, col
}

func (t commonToken) Line() int {
	return t.line
}

func (t commonToken) Col() int {
	return t.col
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
