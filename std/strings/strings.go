// Package strings contains functions for dealing with strings.
package strings

import (
	"strings"

	"github.com/DeedleFake/wdte"
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

	switch len(args) {
	case 0:
		return wdte.GoFunc(Contains)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Contains(frame, append(next, args...)...)
		})
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

	switch len(args) {
	case 0:
		return wdte.GoFunc(Prefix)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Prefix(frame, append(next, args...)...)
		})
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

	switch len(args) {
	case 0:
		return wdte.GoFunc(Suffix)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Suffix(frame, append(next, args...)...)
		})
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

	switch len(args) {
	case 0:
		return wdte.GoFunc(Index)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Index(frame, append(next, args...)...)
		})
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

	switch len(args) {
	case 0:
		return wdte.GoFunc(Repeat)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Repeat(frame, append(args, next...)...)
		})
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

// Join is a WDTE function with the following signatures:
//
//    join strings sep
//    (join sep) strings
//
// It returns a new string containing the strings in the provided
// array with sep inserted between each.
func Join(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("join")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Join)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Join(frame, append(next, args...)...)
		})
	}

	a := args[0].Call(frame).(wdte.Array)
	s := make([]string, 0, len(a))
	for _, str := range a {
		s = append(s, string(str.(wdte.String)))
	}

	return wdte.String(strings.Join(s, string(args[1].Call(frame).(wdte.String))))
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
	"join":   wdte.GoFunc(Join),

	"format": wdte.GoFunc(Format),
})

func init() {
	std.Register("strings", Scope)
}
