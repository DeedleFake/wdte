// Package math contains wdte.Funcs for performing mathematical
// operations.
package math

import (
	"math"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

// A number of useful constants. To see the IDs under which they are
// exported, see Scope.
const (
	E     wdte.Number = math.E
	Pi    wdte.Number = math.Pi
	Sqrt2 wdte.Number = math.Sqrt2
)

// Sin is a WDTE function with the following signature:
//
//    sin n
//
// Returns the sine of n.
func Sin(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Sin)
	}

	frame = frame.Sub("sin")

	a := args[0].(wdte.Number)
	return wdte.Number(math.Sin(float64(a)))
}

// Cos is a WDTE function with the following signature:
//
//    cos n
//
// Returns the cosine of n.
func Cos(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Cos)
	}

	frame = frame.Sub("cos")

	a := args[0].(wdte.Number)
	return wdte.Number(math.Cos(float64(a)))
}

// Tan is a WDTE function with the following signature:
//
//    tan n
//
// Returns the tangent of n.
func Tan(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Tan)
	}

	frame = frame.Sub("tan")

	a := args[0].(wdte.Number)
	return wdte.Number(math.Tan(float64(a)))
}

// Floor is a WDTE function with the following signature:
//
//    floor n
//
// Returns ⌊n⌋.
func Floor(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Floor)
	}

	frame = frame.Sub("floor")

	a := args[0].(wdte.Number)
	return wdte.Number(math.Floor(float64(a)))
}

// Ceil is a WDTE function with the following signature:
//
//    ceil n
//
// Returns ⌈n⌉.
func Ceil(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Ceil)
	}

	frame = frame.Sub("ceil")

	a := args[0].(wdte.Number)
	return wdte.Number(math.Ceil(float64(a)))
}

// Abs is a WDTE function with the following signature:
//
//    abs n
//
// Returns |n|.
func Abs(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Abs)
	}

	frame = frame.Sub("abs")

	a := args[0].(wdte.Number)
	return wdte.Number(math.Abs(float64(a)))
}

// Scope is a scope that contains the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"e":     E,
	"pi":    Pi,
	"sqrt2": Sqrt2,

	"sin": wdte.GoFunc(Sin),
	"cos": wdte.GoFunc(Cos),
	"tan": wdte.GoFunc(Tan),

	"floor": wdte.GoFunc(Floor),
	"ceil":  wdte.GoFunc(Ceil),
	"abs":   wdte.GoFunc(Abs),
})

func init() {
	std.Register("math", Scope)
}
