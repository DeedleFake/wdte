package std

import (
	"fmt"
	"math"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std/arrays"
	stdio "github.com/DeedleFake/wdte/std/io"
	"github.com/DeedleFake/wdte/std/io/file"
	stdmath "github.com/DeedleFake/wdte/std/math"
	"github.com/DeedleFake/wdte/std/stream"
	"github.com/DeedleFake/wdte/std/strings"
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
// returns a function which returns the argument given minus that
// argument.
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
// returns a function which divides the original argument by its own
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
// returns a function which divides the original argument by its own
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
		return save(wdte.GoFunc(Greater), args...)
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

// LessEqual returns true if the first argument is less than or equal
// to the second.
//
// TODO: Document usage of wdte.Comparer.
func LessEqual(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(LessEqual), args...)
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
			return wdte.Bool(c <= 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c >= 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// GreaterEqual returns true if the first argument is greater than or
// equal to the second.
//
// TODO: Document usage of wdte.Comparer.
func GreaterEqual(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) <= 1 {
		return save(wdte.GoFunc(GreaterEqual), args...)
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
			return wdte.Bool(c >= 0)
		}
	}

	if cmp, ok := a2.(wdte.Comparer); ok {
		c, ord := cmp.Compare(a1)
		if ord {
			return wdte.Bool(c <= 0)
		}
	}

	return wdte.Error{
		Err:   fmt.Errorf("Unable to compare %v and %v", a1, a2),
		Frame: frame,
	}
}

// True returns a boolean true.
func True(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return wdte.Bool(true)
}

// False returns a boolean false. This is rarely necessary as most
// built-in functionality considers any value other than a boolean
// true to be false, but it's provided for completeness.
func False(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return wdte.Bool(false)
}

// And returns true if all of its arguments are true.
func And(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("&&")

	switch len(args) {
	case 0:
		return wdte.GoFunc(And)
	}

	for _, arg := range args {
		arg = arg.Call(frame)
		if arg != wdte.Bool(true) {
			return wdte.Bool(false)
		}
	}

	return wdte.Bool(true)
}

// Or returns true if any of its arguments are true.
func Or(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("||")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Or)
	}

	for _, arg := range args {
		arg = arg.Call(frame)
		if arg == wdte.Bool(true) {
			return wdte.Bool(true)
		}
	}

	return wdte.Bool(false)
}

// Not returns true if its argument isn't. Otherwise, it returns
// false.
func Not(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("!")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Not)
	}

	return wdte.Bool(args[0].Call(frame) != wdte.Bool(true))
}

// Module returns a module contiaining the functions in this package.
// It maps mathematical functions to the corresponding mathematical
// symbols. For example, Add() becomes `+`, Sub() becomes `-`, and so
// on. Comparisons get mapped to the cooresponding comparison symbols
// from C-style languages. For example, Equals() becomes `==`.
//
// This module is useful as a starting point for parsing. For example,
// if you want to parse a module and you want it to have access to
// these functions, you can use
//
//     m, err := std.Module().Parse(r, im)
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"+": wdte.GoFunc(Add),
			"-": wdte.GoFunc(Sub),
			"*": wdte.GoFunc(Mult),
			"/": wdte.GoFunc(Div),
			"%": wdte.GoFunc(Mod),

			"==":    wdte.GoFunc(Equals),
			"<":     wdte.GoFunc(Less),
			">":     wdte.GoFunc(Greater),
			"<=":    wdte.GoFunc(LessEqual),
			">=":    wdte.GoFunc(GreaterEqual),
			"true":  wdte.GoFunc(True),
			"false": wdte.GoFunc(False),
			"&&":    wdte.GoFunc(And),
			"||":    wdte.GoFunc(Or),
			"!":     wdte.GoFunc(Not),
		},
	}
}

// Import provides a simple importer that imports standard library
// modules.
var Import = wdte.ImportFunc(stdImporter)

func stdImporter(from string) (*wdte.Module, error) {
	switch from {
	case "stream":
		return stream.Module(), nil

	case "math":
		return stdmath.Module(), nil

	case "io":
		return stdio.Module(), nil
	case "io/file":
		return file.Module(), nil

	case "strings":
		return strings.Module(), nil
	case "arrays":
		return arrays.Module(), nil
	}

	return nil, fmt.Errorf("Unknown import: %q", from)
}
