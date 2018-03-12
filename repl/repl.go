package repl

import (
	"bufio"
	"bytes"
	"io"

	"github.com/DeedleFake/wdte"
)

// A NextFunc returns the next piece of code to be interpreted. When
// reading from stdin, this is likely the next line entered.
type NextFunc func() ([]byte, error)

// SimpleNext returns a next func that scans lines from r.
func SimpleNext(r io.Reader) NextFunc {
	s := bufio.NewScanner(r)
	return func() ([]byte, error) {
		more := s.Scan()
		if !more {
			err := s.Err()
			if err != nil {
				err = io.EOF
			}
			return nil, err
		}

		return s.Bytes(), nil
	}
}

// REPL provides a means to track a global state and interpret
// separate lines of a WDTE script, allowing for the implementation of
// a read-eval-print loop.
type REPL struct {
	next  NextFunc
	im    wdte.Importer
	scope *wdte.Scope
}

// New creates a new REPL which reads with next, imports using im, and
// executes the first line with the scope start.
func New(next NextFunc, im wdte.Importer, start *wdte.Scope) *REPL {
	return &REPL{
		next:  next,
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

	src, err := r.next()
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return nil, err
	}

	m, err := wdte.Parse(bytes.NewReader(src), r.im)
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
