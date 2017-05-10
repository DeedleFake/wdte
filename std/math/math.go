// Package math contains wdte.Funcs for performing mathematical
// operations.
package math

import (
	"math"

	"github.com/DeedleFake/wdte"
)

// Pi ignores its arguments and returns π as a wdte.Number.
func Pi(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	return wdte.Number(math.Pi)
}

// Sin returns sin(args[0]).
func Sin(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Sin)
	}

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Sin(a))
}

// Cos returns cos(args[0]).
func Cos(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Cos)
	}

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Cos(a))
}

// Tan returns tan(args[0]).
func Tan(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Tan)
	}

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Tan(a))
}

// Import returns a module that contains the functions in this
// package. This can be used by an Importer to import them more
// easily. The functions in the returned module have the same names as
// those in this package except that they are lowercase.
func Import() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"pi": wdte.GoFunc(Pi),

			"sin": wdte.GoFunc(Sin),
			"cos": wdte.GoFunc(Cos),
			"tan": wdte.GoFunc(Tan),
		},
	}
}
