// Package math contains wdte.Funcs for performing mathematical
// operations.
package math

import (
	"math"

	"github.com/DeedleFake/wdte"
)

// Pi ignores its arguments and returns π as a wdte.Number.
func Pi(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return wdte.Number(math.Pi)
}

// Sin returns sin(args[0]).
func Sin(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Sin)
	}

	frame = frame.Sub("sin")

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Sin(float64(a)))
}

// Cos returns cos(args[0]).
func Cos(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Cos)
	}

	frame = frame.Sub("cos")

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Cos(float64(a)))
}

// Tan returns tan(args[0]).
func Tan(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Tan)
	}

	frame = frame.Sub("tan")

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Tan(float64(a)))
}

// Floor returns ⌊args[0]⌋.
func Floor(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Floor)
	}

	frame = frame.Sub("floor")

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Floor(float64(a)))
}

// Ceil returns ⌈args[0]⌉.
func Ceil(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Ceil)
	}

	frame = frame.Sub("ceil")

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Ceil(float64(a)))
}

// Abs returns |args[0]|.
func Abs(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Abs)
	}

	frame = frame.Sub("abs")

	a := args[0].Call(frame).(wdte.Number)
	return wdte.Number(math.Abs(float64(a)))
}

// Module returns a module that contains the functions in this
// package. This can be used by an Importer to import them more
// easily. The functions in the returned module have the same names as
// those in this package except that they are lowercase.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"pi": wdte.GoFunc(Pi),

			"sin": wdte.GoFunc(Sin),
			"cos": wdte.GoFunc(Cos),
			"tan": wdte.GoFunc(Tan),

			"floor": wdte.GoFunc(Floor),
			"ceil":  wdte.GoFunc(Ceil),
			"abs":   wdte.GoFunc(Abs),
		},
	}
}
