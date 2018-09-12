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
	// Scope is the scope that the next line will be executed in. It is
	// automatically updated every time an executed line changes the
	// scope.
	Scope *wdte.Scope

	next NextFunc
	im   wdte.Importer

	stack []string
	buf   []byte
}

// New creates a new REPL which reads with next, imports using im, and
// executes the first line with the scope start.
func New(next NextFunc, im wdte.Importer, start *wdte.Scope) *REPL {
	return &REPL{
		next:  next,
		im:    im,
		Scope: start,
	}
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

	stack, partial := Partial(bytes.NewReader(src), r.stack)
	r.stack = stack
	if partial {
		r.buf = append(r.buf, src...)
		return nil, ErrIncomplete
	}

	src = append(r.buf, src...)
	r.buf = r.buf[:0]

	m, err := wdte.Parse(bytes.NewReader(src), r.im)
	if err != nil {
		return nil, err
	}

	frame := wdte.F().WithScope(r.Scope)
	next, ret := m.Collect(frame)
	if err, ok := ret.(error); ok {
		return nil, err
	}

	r.Scope = next
	return ret.Call(frame), nil
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

func peek(stack []string) string {
	if len(stack) == 0 {
		return ""
	}

	return stack[len(stack)-1]
}

func pop(stack []string) (string, []string) {
	if len(stack) == 0 {
		return "", stack
	}

	p := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return p, stack
}

// Partial checks if an expression, read from r, is incomplete. The
// initial value of stack should be nil, and subsequent values should
// be the value of the first return. The second return is true if the
// expression was incomplete.
func Partial(r io.Reader, stack []string) ([]string, bool) {
	s := scanner.New(r)
	var prev scanner.Token
	for s.Scan() {
		switch tok := s.Tok(); tok.Type {
		case scanner.Keyword:
			switch v := tok.Val.(string); v {
			case "(", "(@":
				stack = append(stack, ")")
			case "[":
				stack = append(stack, "]")
			case "{":
				stack = append(stack, "}")

			case ")", "]", "}":
				if v == peek(stack) {
					_, stack = pop(stack)
				}
			}

		case scanner.EOF:
			if (len(stack) > 0) || (prev.Val != ";") {
				return stack, true
			}
		}

		prev = s.Tok()
	}

	return stack, false
}
