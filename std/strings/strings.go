// Package strings contains functions for dealing with strings.
package strings

import (
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

// Contains returns true if the second argument is a substring of the
// first. If only given one argument, it returns a function that
// checks if that argument is a substring of its own argument.
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

// Prefix returns true if the first argument starts with the second.
// If given only one argument, it returns a function that checks if
// that argument is a prefix of its own argument.
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

// Suffix returns true if the first argument ends with the second. If
// given only one argument, it returns a function that checks if that
// argument is a suffix of its own argument.
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

// Index searches the first argument for the second argument,
// returning the index of the beginning of its first instance, or -1
// if its not present. If only given one argument, Index returns a
// function which searches other strings for that argument.
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

// Upper returns its argument converted to uppercase.
func Upper(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("upper")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Upper)
	}

	return wdte.String(strings.ToUpper(string(args[0].Call(frame).(wdte.String))))
}

// Lower returns its argument converted to lowercase.
func Lower(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("lower")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Lower)
	}

	return wdte.String(strings.ToLower(string(args[0].Call(frame).(wdte.String))))
}

// S returns a scope containing the various functions in this package.
func S() *wdte.Scope {
	return wdte.S().Map(map[wdte.ID]wdte.Func{
		"contains": wdte.GoFunc(Contains),
		"prefix":   wdte.GoFunc(Prefix),
		"suffix":   wdte.GoFunc(Suffix),
		"index":    wdte.GoFunc(Index),

		"upper": wdte.GoFunc(Upper),
		"lower": wdte.GoFunc(Lower),

		"format": wdte.GoFunc(Format),
	})
}

func init() {
	std.Register("strings", S())
}
