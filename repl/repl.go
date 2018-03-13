package repl

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/scanner"
)

var (
	// ErrIncomplete is returned by REPL.Next() if it expects more input
	// before it will begin evaluation.
	ErrIncomplete = errors.New("input incomplete")
)

// A NextFunc returns the next piece of code to be interpreted. When
// reading from stdin, this is likely the next line entered.
//
// If NextFunc should not be called again, it should return
// nil, io.EOF.
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

	stack []string
	buf   []byte
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

func (r *REPL) peek() string {
	if len(r.stack) == 0 {
		return ""
	}

	return r.stack[len(r.stack)-1]
}

func (r *REPL) push(v string) {
	r.stack = append(r.stack, v)
}

func (r *REPL) pop() string {
	if len(r.stack) == 0 {
		return ""
	}

	p := r.stack[len(r.stack)-1]
	r.stack = r.stack[:len(r.stack)-1]
	return p
}

// Next reads and evaluates the next line of input. It returns the
// value returned from that line, or an error if one is encountered.
// If the end of the input has been reached, it will return nil, nil.
//
// If an input ends in a partial expression, such as a single line of
// a mult-line expression, nil, ErrIncomplete is returned.
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

	s := scanner.New(bytes.NewReader(src))
	var prev scanner.Token
	for s.Scan() {
		switch tok := s.Tok(); tok.Type {
		case scanner.Keyword:
			switch v := tok.Val.(string); v {
			case "(", "(@":
				r.push(")")
			case "[":
				r.push("]")
			case "{":
				r.push("}")

			case ")", "]", "}":
				if v == r.peek() {
					r.pop()
				}
			}

		case scanner.EOF:
			if (len(r.stack) > 0) || (prev.Val != ";") {
				r.buf = append(r.buf, src...)
				return nil, ErrIncomplete
			}
		}

		prev = s.Tok()
	}

	src = append(r.buf, src...)
	r.buf = r.buf[:0]

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

// Cancel cancels a partial expression. This is useful if, for
// example, a user sends an interrupt to a command-line REPL while
// entering a subsequent line of a multi-line expression.
//
// If a partial expression is not in progress, this has no effect.
func (r *REPL) Cancel() {
	r.stack = r.stack[:0]
	r.buf = r.buf[:0]
}
