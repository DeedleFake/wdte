package repl

import (
	"bufio"
	"bytes"
	"io"

	"github.com/DeedleFake/wdte"
)

// REPL provides a means to track a global state and interpret
// separate lines of a WDTE script, allowing for the implementation of
// a read-eval-print loop.
type REPL struct {
	s     *bufio.Scanner
	im    wdte.Importer
	scope *wdte.Scope
}

// New creates a new REPL which reads from r, imports using im, and
// executes the first line with the scope start.
func New(r io.Reader, im wdte.Importer, start *wdte.Scope) *REPL {
	return &REPL{
		s:     bufio.NewScanner(r),
		im:    im,
		scope: start,
	}
}

// Next reads and evaluates the next line of input. It returns the
// value returned from that line, or an error if one is encountered.
// If the end of the input has been reached, it will return nil, nil.
//
// BUG: Currently, this simply reads the input line-by-line. This
// means that it can't be used to interpret multi-line, complex
// expressions. See #62.
func (r *REPL) Next() (ret wdte.Func, err error) {
	defer func() {
		switch e := recover().(type) {
		case error:
			err = e

		case nil:

		default:
			panic(e)
		}
	}()

	if !r.s.Scan() {
		err := r.s.Err()
		return nil, err
	}

	src := bytes.NewReader(r.s.Bytes())
	m, err := wdte.Parse(src, r.im)
	if err != nil {
		return nil, err
	}

	next, ret := m.Collect(wdte.F().WithScope(r.scope))
	if err, ok := ret.(error); ok {
		return nil, err
	}

	r.scope = next
	return ret, nil
}
