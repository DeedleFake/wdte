package wdte

import (
	"bytes"
	"fmt"
)

// A String is a string, as parsed from a string literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type String string

func (s String) Call(frame Frame, args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

func (s String) Compare(other Func) (int, bool) {
	o, ok := other.(String)
	if !ok {
		return -1, false
	}

	switch {
	case s < o:
		return -1, true
	case s > o:
		return 1, true
	}

	return 0, true
}

// A Number is a number, as parsed from a number literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type Number float64

func (n Number) Call(frame Frame, args ...Func) Func {
	// TODO: Use the arguments for something, perhaps.
	return n
}

func (n Number) Compare(other Func) (int, bool) {
	o, ok := other.(Number)
	if !ok {
		return -1, false
	}

	return int(n - o), true
}

// An Array represents a WDTE array type. It's similar to a Compound,
// but doesn't evaluate its own members. Instead, evaluation simply
// yields the array, much like how strings and numbers work.
type Array []Func

func (a Array) Call(frame Frame, args ...Func) Func {
	return a
}

//func (a Array)Compare(other Func) (int, bool) {
//	TODO: Implement this. I'm not sure if it should support ordering
//	or not. I'm also not sure if it should call its elements in order
//	to get their underlying values. It probably should.
//}

// An Error is returned by any of the built-in functions when they run
// into an error.
type Error struct {
	// Err is the error that generated the Error. In a lot of cases,
	// this is just a simple error message.
	Err error

	// Frame is the frame of the function that the error was first
	// generated in.
	Frame Frame
}

func (e Error) Call(frame Frame, args ...Func) Func {
	return e
}

func (e Error) Error() string {
	var buf bytes.Buffer
	e.Frame.Backtrace(&buf)

	return fmt.Sprintf("WDTE Error: %v\n%s", e.Err, buf.Bytes())
}

type Bool bool

func (b Bool) Call(frame Frame, args ...Func) Func {
	return b
}

func (b Bool) Compare(other Func) (int, bool) {
	if b == other {
		return 0, false
	}
	return -1, false
}
