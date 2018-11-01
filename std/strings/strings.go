// Package strings contains functions for dealing with strings.
package strings

import (
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/auto"
	"github.com/DeedleFake/wdte/std"
)

// Contains is a WDTE function with the following signatures:
//
//    contains outer inner
//    (contains inner) outer
//
// Returns true if inner is a substring of outer.
func Contains(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("contains")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Contains), args...)
	}

	haystack := args[0].Call(frame).(wdte.String)
	needle := args[1].Call(frame).(wdte.String)

	return wdte.Bool(strings.Contains(string(haystack), string(needle)))
}

// Prefix is a WDTE function with the following signatures:
//
//    prefix s p
//    (prefix p) s
//
// Returns true if p is a prefix of s.
func Prefix(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("prefix")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Prefix), args...)
	}

	haystack := args[0].Call(frame).(wdte.String)
	needle := args[1].Call(frame).(wdte.String)

	return wdte.Bool(strings.HasPrefix(string(haystack), string(needle)))
}

// Suffix is a WDTE function with the following signatures:
//
//    suffix s p
//    (suffix p) s
//
// Returns true if p is a suffix of s.
func Suffix(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("suffix")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Suffix), args...)
	}

	haystack := args[0].Call(frame).(wdte.String)
	needle := args[1].Call(frame).(wdte.String)

	return wdte.Bool(strings.HasSuffix(string(haystack), string(needle)))
}

// Index is a WDTE function with the following signatures:
//
//    index outer inner
//    (index inner) outer
//
// It returns the index of the first character of the first instances
// of inner in outer. If inner is not a substring of outer, it returns
// -1.
func Index(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("index")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Index), args...)
	}

	haystack := args[0].Call(frame).(wdte.String)
	needle := args[1].Call(frame).(wdte.String)

	return wdte.Number(strings.Index(string(haystack), string(needle)))
}

// Upper is a WDTE function with the following signatures:
//
//    upper s
//
// It returns s converted to uppercase.
func Upper(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("upper")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Upper)
	}

	return wdte.String(strings.ToUpper(string(args[0].Call(frame).(wdte.String))))
}

// Lower is a WDTE function with the following signature:
//
//    lower s
//
// It returns s converted to lowercase.
func Lower(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("lower")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Lower)
	}

	return wdte.String(strings.ToLower(string(args[0].Call(frame).(wdte.String))))
}

// Repeat is a WDTE function with the following signatures:
//
//    repeat string times
//    (repeat string) times
//    repeat times string
//    (repeat times) string
//
// It returns a new string containing the given string repeated the
// number of times specified.
func Repeat(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("repeat")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Repeat), args...)
	}

	var str wdte.String
	var times wdte.Number
	switch a0 := args[0].Call(frame).(type) {
	case wdte.String:
		str = a0
		times = args[1].Call(frame).(wdte.Number)

	case wdte.Number:
		times = a0
		str = args[1].Call(frame).(wdte.String)
	}

	return wdte.String(strings.Repeat(string(str), int(times)))
}

// Split is a WDTE function with the following signatures:
//
//    split string sep
//    (split sep) string
//    split string sep n
//    (split sep n) string
//    (split sep) string n
//
// It splits the given string around instances of the given seperator
// string. If n is provided and is positive, the returned array of
// strings will have at most n elements. Note that this behavior
// differs from the Go standard library's string splitting function in
// that a zero value for n does not cause the function to return an
// empty array.
func Split(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("split")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Split), args...)
	}

	str := args[0].Call(frame).(wdte.String)

	var sep wdte.String
	n := wdte.Number(-1)
	switch arg := args[1].Call(frame).(type) {
	case wdte.String:
		sep = arg
		if len(args) > 2 {
			n = args[2].Call(frame).(wdte.Number)
		}

	case wdte.Number:
		if len(args) < 3 {
			return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
				return Split(frame, append(next, args...)...)
			})
		}

		sep = args[2].Call(frame).(wdte.String)
		n = arg
	}
	if n == 0 {
		n = -1
	}

	split := strings.SplitN(string(str), string(sep), int(n))

	out := make(wdte.Array, 0, len(split))
	for _, part := range split {
		out = append(out, wdte.String(part))
	}
	return out
}

// Join is a WDTE function with the following signatures:
//
//    join strings sep
//    (join sep) strings
//
// It returns a new string containing the strings in the provided
// array with sep inserted between each.
func Join(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("join")

	if len(args) < 2 {
		return auto.SaveArgsReverse(wdte.GoFunc(Join), args...)
	}

	a := args[0].Call(frame).(wdte.Array)
	s := make([]string, 0, len(a))
	for _, str := range a {
		s = append(s, string(str.(wdte.String)))
	}

	return wdte.String(strings.Join(s, string(args[1].Call(frame).(wdte.String))))
}

type reader struct {
	*strings.Reader
}

func (r reader) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return r
}

func (r reader) Reflect(name string) bool { // nolint
	return name == "Reader"
}

// Read is a WDTE function with the following signature:
//
//    read s
//
// Returns a reader which reads from the string s.
func Read(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("read")

	if len(args) == 0 {
		return wdte.GoFunc(Read)
	}

	s := args[0].Call(frame).(wdte.String)
	return reader{Reader: strings.NewReader(string(s))}
}

// Scope is a scope containing the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"contains": wdte.GoFunc(Contains),
	"prefix":   wdte.GoFunc(Prefix),
	"suffix":   wdte.GoFunc(Suffix),
	"index":    wdte.GoFunc(Index),

	"upper":  wdte.GoFunc(Upper),
	"lower":  wdte.GoFunc(Lower),
	"repeat": wdte.GoFunc(Repeat),
	"split":  wdte.GoFunc(Split),
	"join":   wdte.GoFunc(Join),

	"read":   wdte.GoFunc(Read),
	"format": wdte.GoFunc(Format),
})

func init() {
	std.Register("strings", Scope)
}
