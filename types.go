package wdte

import "fmt"

// A String is a string, as parsed from a string literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type String string

func (s String) Call(frame Frame, args ...Func) Func {
	// TODO: Use the arguments for something. Probably concatenation.
	return s
}

func (s String) Equals(other Func) bool {
	o, ok := other.(String)
	return ok && (s == o)
}

// A Number is a number, as parsed from a number literal. That's about
// it. Like everything else, it's a function. It simply returns itself
// when called.
type Number float64

func (n Number) Call(frame Frame, args ...Func) Func {
	// TODO: Use the arguments for something, perhaps.
	return n
}

func (n Number) Equals(other Func) bool {
	o, ok := other.(Number)
	return ok && (n == o)
}

// An Array represents a WDTE array type. It's similar to a Compound,
// but doesn't evaluate its own members. Instead, evaluation simply
// yields the array, much like how strings and numbers work.
type Array []Func

func (a Array) Call(frame Frame, args ...Func) Func {
	return a
}

func (a Array) Equals(other Func) bool {
	o, ok := other.(Array)
	if !ok || (len(a) != len(o)) {
		return false
	}

	for i := range a {
		if !a[i].Equals(o[i]) {
			return false
		}
	}

	return true
}

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

func (e Error) Equals(other Func) bool {
	panic("Not implemented.")
}

func (e Error) Error() string {
	return fmt.Sprintf("Error in %v: %v", e.Frame.ID(), e.Err)
}
