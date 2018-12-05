package wdte

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
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

// An Atter is a Func that can be indexed, like an array or a string.
type Atter interface {
	// At returns the value at index i. If the index is out of range, it
	// should return false as its second return value.
	At(i Func) (Func, bool)
}

// A Reflector is a Func that can determine if it can be treated as
// the named type or not. For example,
//
//    s := wdte.String("example")
//    return s.Reflect("string")
//
// returns true.
type Reflector interface {
	Reflect(name string) bool
}

// Reflect checks if a Func can be considered to be of a given type.
// If v implements Reflector, v.Reflect(name) is used to check for
// compatability. If not, a simple string comparison is done against
// whatever Go's reflect package claims the short name of the
// underlying type to be.
func Reflect(f Func, name string) bool {
	if r, ok := f.(Reflector); ok {
		return r.Reflect(name)
	}

	return reflect.TypeOf(f).Name() == name
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

func (s String) At(i Func) (Func, bool) { // nolint
	if i, ok := i.(Number); ok {
		if (int(i) < 0) || (int(i) >= len(s)) {
			return nil, false
		}

		return String(s[int(i)]), true
	}

	return nil, false
}

func (s String) Reflect(name string) bool { // nolint
	return name == "String"
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

func (n Number) String() string { // nolint
	bn := big.NewFloat(float64(n))
	if bn.IsInt() {
		return bn.Text('f', -1)
	}

	return bn.Text('g', 10)
}

func (n Number) Reflect(name string) bool { // nolint
	return name == "Number"
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

func (a Array) At(i Func) (Func, bool) { // nolint
	if i, ok := i.(Number); ok {
		if (int(i) < 0) || (int(i) >= len(a)) {
			return nil, false
		}

		return a[int(i)], true
	}

	return nil, false
}

func (a Array) String() string { // nolint
	var buf strings.Builder

	buf.WriteByte('[')
	var pre string
	for _, f := range a {
		buf.WriteString(pre)
		fmt.Fprint(&buf, f)
		pre = "; "
	}
	buf.WriteByte(']')

	return buf.String()
}

func (a Array) Reflect(name string) bool { // nolint
	return name == "Array"
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
	return e.Err.Error()
}

func (e Error) Reflect(name string) bool { // nolint
	return name == "Error"
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

func (b Bool) Reflect(name string) bool { // nolint
	return name == "Bool"
}
