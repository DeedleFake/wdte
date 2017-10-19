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

func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"contains": wdte.GoFunc(Contains),
		},
	}
}
