package wdte

import (
	"bytes"
	"fmt"
)

// A Comparer is a Func that is able to be compared to other
// functions.
type Comparer interface {
	// Compare returns two values. The meaning of the first is dependent
	// upon the second. If the second is true, then the first indicates
	// ordering via the standard negative, positive, and zero results to
	// indicate less than, greater than, and equal, respectively. If the
	// second is false, then the first indicates only equality, with
	// zero still meaning equal, but other values simply meaning unequal.
	Compare(other Func) (int, bool)
}

// A Lenner is a Func that has a length, such as arrays and strings.
type Lenner interface {
	Len() int
}

// A String is a string, as parsed from a string literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type String string

func (s String) Call(frame Frame, args ...Func) Func { // nolint
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

func (s String) Compare(other Func) (int, bool) { // nolint
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

func (s String) Len() int { // nolint
	return len(s)
}

// A Number is a number, as parsed from a number literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type Number float64

func (n Number) Call(frame Frame, args ...Func) Func { // nolint
	// TODO: Use the arguments for something, perhaps.
	return n
}

func (n Number) Compare(other Func) (int, bool) { // nolint
	o, ok := other.(Number)
	if !ok {
		return -1, false
	}

	return int(n - o), true
}

// An Array represents a WDTE array type. It's similar to a Compound,
// but when evaluated, it returns itself with its own members replaced
// with their own evaluations. This allows it to be passed around as a
// value in the same way as strings and numbers.
type Array []Func

func (a Array) Call(frame Frame, args ...Func) Func { // nolint
	n := make(Array, 0, len(a))
	for i := range a {
		n = append(n, a[i].Call(frame))
	}
	return n
}

func (a Array) Len() int { // nolint
	return len(a)
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

func (e Error) Call(frame Frame, args ...Func) Func { // nolint
	return e
}

func (e Error) Error() string {
	var buf bytes.Buffer
	_ = e.Frame.Backtrace(&buf)

	return fmt.Sprintf("WDTE Error: %v\n%s", e.Err, buf.Bytes())
}

// Bool is a boolean. Like other primitive types, it simply returns
// itself when called.
type Bool bool

func (b Bool) Call(frame Frame, args ...Func) Func { // nolint
	return b
}

func (b Bool) Compare(other Func) (int, bool) { // nolint
	if b == other {
		return 0, false
	}
	return -1, false
}
