// Package strings contains functions for dealing with strings.
package strings

import (
	"strings"

	"github.com/DeedleFake/wdte"
)

// Contains returns true if the second argument is a substring of the
// first. If only given one argument, it returns a function that
// checks if that argument is a substring of its own argument.
func Contains(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("contains")

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
	frame = frame.WithID("prefix")

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
	frame = frame.WithID("suffix")

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
	frame = frame.WithID("index")

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

// Len returns the length of its argument.
func Len(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("len")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Len)
	}

	return wdte.Number(len(string(args[0].Call(frame).(wdte.String))))
}

// At returns the ith character of a string. The order of arguments
// is the string and then the index. If given only one argument, it
// uses that argument as the index to get characters from in strings
// that it is passed.
func At(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("at")

	switch len(args) {
	case 0:
		return wdte.GoFunc(At)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return At(frame, append(next, args...)...)
		})
	}

	str := args[0].Call(frame).(wdte.String)
	i := args[1].Call(frame).(wdte.Number)

	return wdte.String(str[int(i)])
}

// Upper returns its argument converted to uppercase.
func Upper(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("upper")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Upper)
	}

	return wdte.String(strings.ToUpper(string(args[0].Call(frame).(wdte.String))))
}

// Lower returns its argument converted to lowercase.
func Lower(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("lower")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Lower)
	}

	return wdte.String(strings.ToLower(string(args[0].Call(frame).(wdte.String))))
}

func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"contains": wdte.GoFunc(Contains),
			"prefix":   wdte.GoFunc(Prefix),
			"suffix":   wdte.GoFunc(Suffix),
			"index":    wdte.GoFunc(Index),

			"len": wdte.GoFunc(Len),
			"at":  wdte.GoFunc(At),

			"upper": wdte.GoFunc(Upper),
			"lower": wdte.GoFunc(Lower),

			"format": wdte.GoFunc(Format),
		},
	}
}
