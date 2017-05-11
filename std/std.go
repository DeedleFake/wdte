package std

import (
	"math"

	"github.com/DeedleFake/wdte"
)

func save(f wdte.Func, saved ...wdte.Func) wdte.Func {
	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		return f.Call(frame, append(saved, args...)...)
	})
}

// Add returns the sum of its arguments. If called with only 1
// argument, it returns a function which adds arguments given to that
// one argument.
func Add(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Add)

	case 1:
		return save(wdte.GoFunc(Add), args[0])
	}

	frame = frame.WithID("+")

	var sum wdte.Number
	for _, arg := range args {
		sum += arg.Call(frame).(wdte.Number)
	}
	return sum
}

// Sub returns args[0] - args[1]. If called with only 1 argument, it
// returns a function which returns that argument minus the argument
// given.
func Sub(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Sub)

	case 1:
		return save(wdte.GoFunc(Sub), args[0])
	}

	frame = frame.WithID("-")

	a1 := args[0].Call(frame).(wdte.Number)
	a2 := args[1].Call(frame).(wdte.Number)
	return a1 - a2
}

// Mult returns the product of its arguments. If called with only 1
// argument, it returns a function that multiplies that argument by
// its own arguments.
func Mult(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Mult)

	case 1:
		return save(wdte.GoFunc(Mult), args[0])
	}

	frame = frame.WithID("*")

	p := wdte.Number(1)
	for _, arg := range args {
		p *= arg.Call(frame).(wdte.Number)
	}
	return p
}

// Div returns args[0] / args[1]. If called with only 1 argument, it
// returns a function which divides its own argument by the original
// argument.
func Div(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Div)

	case 1:
		return save(wdte.GoFunc(Div), args[0])
	}

	frame = frame.WithID("/")

	a1 := args[0].Call(frame).(wdte.Number)
	a2 := args[1].Call(frame).(wdte.Number)
	return a1 / a2
}

// Mod returns args[0] % args[1]. If called with only 1 argument, it
// returns a function which divides its own argument by the original
// argument.
func Mod(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Mod)

	case 1:
		return save(wdte.GoFunc(Mod), args[0])
	}

	frame = frame.WithID("%")

	a1 := args[0].Call(frame).(wdte.Number)
	a2 := args[1].Call(frame).(wdte.Number)
	return wdte.Number(math.Mod(float64(a1), float64(a2)))
}

// Insert adds the functions in this package to m. It maps them to the
// corresponding mathematical operators. For example, Add() becomes
// `+`, Sub() becomes `-`, and so on.
func Insert(m *wdte.Module) {
	m.Funcs["+"] = wdte.GoFunc(Add)
	m.Funcs["-"] = wdte.GoFunc(Sub)
	m.Funcs["*"] = wdte.GoFunc(Mult)
	m.Funcs["/"] = wdte.GoFunc(Div)
	m.Funcs["%"] = wdte.GoFunc(Mod)
}
