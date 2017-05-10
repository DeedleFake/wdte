package std

import "github.com/DeedleFake/wdte"

func save(f wdte.Func, saved ...wdte.Func) wdte.Func {
	return wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		return f(frame, append(saved, args...))
	})
}

// Add returns the sum of its arguments. If called with only 1
// argument, it returns a function which adds arguments given to that
// one argument.
func Add(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return Add

	case 1:
		return save(Add, args[0])
	}

	var sum wdte.Number
	for _, arg := range args {
		sum += arg.Call(frame).(wdte.Number)
	}
	return sum
}

// Sub returns args[0] - args[1]. If called with only 1 argument, it
// returns a function which returns that argument minus the argument
// given.
func Sub(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return Sub

	case 1:
		return save(Sub, args[0])
	}

	a1 := args[0].Call(frame).(wdte.Number)
	a2 := args[1].Call(frame).(wdte.Number)
	return a1 - a2
}

// Mult returns the product of its arguments. If called with only 1
// argument, it returns a function that multiplies that argument by
// its own arguments.
func Mult(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return Mult

	case 1:
		return save(Mult, args[0])
	}

	p := wdte.Number(1)
	for _, arg := range args {
		p *= arg.Call(frame).(wdte.Number)
	}
	return p
}

// Div returns args[0] / args[1]. If called with only 1 argument, it
// returns a function which divides its own argument by the original
// argument.
func Div(frame []wdte.Func, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return Div

	case 1:
		return save(Div, args[0])
	}

	a1 := args[0].Call(frame).(wdte.Number)
	a2 := args[1].Call(frame).(wdte.Number)
	return a1 / a2
}

// Insert adds the functions in this package to m. It maps them to the
// cooresponding mathematical operators. For example, Add() becomes
// `+`, Sub() becomes `-`, and so on.
func Insert(m *wdte.Module) {
	m.Funcs["+"] = wdte.GoFunc(Add)
	m.Funcs["-"] = wdte.GoFunc(Sub)
	m.Funcs["*"] = wdte.GoFunc(Mult)
	m.Funcs["/"] = wdte.GoFunc(Div)
}
