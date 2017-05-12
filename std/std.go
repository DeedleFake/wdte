package std

import (
	"fmt"
	"math"

	"github.com/DeedleFake/wdte"
)

func save(f wdte.Func, saved ...wdte.Func) wdte.Func {
	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		return f.Call(frame, append(args, saved...)...)
	})
}

// Add returns the sum of its arguments. If called with only 1
// argument, it returns a function which adds arguments given to that
// one argument.
func Add(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Add), args...)
	}

	frame = frame.WithID("+")

	var sum wdte.Number
	for _, arg := range args {
		arg = arg.Call(frame)
		if _, ok := arg.(error); ok {
			return arg
		}
		sum += arg.(wdte.Number)
	}
	return sum
}

// Sub returns args[0] - args[1]. If called with only 1 argument, it
// returns a function which returns that argument minus the argument
// given.
func Sub(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Sub), args...)
	}

	frame = frame.WithID("-")

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	return a1.(wdte.Number) - a2.(wdte.Number)
}

// Mult returns the product of its arguments. If called with only 1
// argument, it returns a function that multiplies that argument by
// its own arguments.
func Mult(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Mult), args...)
	}

	frame = frame.WithID("*")

	p := wdte.Number(1)
	for _, arg := range args {
		arg = arg.Call(frame)
		if _, ok := arg.(error); ok {
			return arg
		}
		p *= arg.(wdte.Number)
	}
	return p
}

// Div returns args[0] / args[1]. If called with only 1 argument, it
// returns a function which divides its own argument by the original
// argument.
func Div(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Div), args...)
	}

	frame = frame.WithID("/")

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	return a1.(wdte.Number) / a2.(wdte.Number)
}

// Mod returns args[0] % args[1]. If called with only 1 argument, it
// returns a function which divides its own argument by the original
// argument.
func Mod(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Mod), args...)
	}

	frame = frame.WithID("%")

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	return wdte.Number(math.Mod(
		float64(a1.(wdte.Number)),
		float64(a2.(wdte.Number)),
	))
}

// Equals checks if two values are equal or not. If called with only
// one argument, it returns a function which checks that argument for
// equality with other values.
//
// If the first argument given implements wdte.Comparer, it is used
// for the comparison. If not, and the second does, then that is used.
// If neither does, a simple, direct Go equality check is used.
func Equals(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Equals), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, _ := cmp.Compare(a2)
		return wdte.Bool(c == 0)
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, _ := cmp.Compare(a1)
		return wdte.Bool(c == 0)
	}

	return wdte.Bool(a1 == a2)
}

// Less returns true if the first argument is less than the second.
// It returns an error if two arguments can't be compared.
//
// TODO: Document usage of wdte.Comparer.
func Less(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Less), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a2)
		if ord {
			return wdte.Bool(c < 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c > 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// Greater returns true if the first argument is less than the second.
// It returns an error if two arguments can't be compared.
//
// TODO: Document usage of wdte.Comparer.
func Greater(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(Less), args...)
	}

	a1 := args[0].Call(frame)
	if _, ok := a1.(error); ok {
		return a1
	}

	a2 := args[1].Call(frame)
	if _, ok := a2.(error); ok {
		return a2
	}

	if cmp, ok := a1.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a2)
		if ord {
			return wdte.Bool(c > 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c < 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// Insert adds the functions in this package to m. It maps
// mathematical functions to the corresponding mathematical symbols.
// For example, Add() becomes `+`, Sub() becomes `-`, and so on.
// Comparisons get mapped to the cooresponding comparison symbols from
// C-style languages. For example, Equals() becomes `==`.
func Insert(m *wdte.Module) {
	m.Funcs["+"] = wdte.GoFunc(Add)
	m.Funcs["-"] = wdte.GoFunc(Sub)
	m.Funcs["*"] = wdte.GoFunc(Mult)
	m.Funcs["/"] = wdte.GoFunc(Div)
	m.Funcs["%"] = wdte.GoFunc(Mod)

	m.Funcs["=="] = wdte.GoFunc(Equals)
	m.Funcs["<"] = wdte.GoFunc(Less)
	m.Funcs[">"] = wdte.GoFunc(Greater)
}
