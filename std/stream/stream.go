// Package stream provides WDTE functions for manipulating streams of
// data.
package stream

import "github.com/DeedleFake/wdte"

// A Stream is a type of function that can yield successive values.
type Stream interface {
	wdte.Func

	// Next returns the next value and true, or an undefined value and
	// false if the stream is empty.
	Next(frame wdte.Frame) (wdte.Func, bool)
}

// A NextFunc wraps a Go function, making it possible to use it as a
// Stream. When called as a WDTE function, the function simply returns
// itself.
type NextFunc func(frame wdte.Frame) (wdte.Func, bool)

func (n NextFunc) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return n
}

func (n NextFunc) Next(frame wdte.Frame) (wdte.Func, bool) { // nolint
	return n(frame)
}

// Module returns a module for easy importing into an actual script.
// The imported functions have the same names as the functions in this
// package, except that the first letter is lowercase.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"new":    wdte.GoFunc(New),
			"range":  wdte.GoFunc(Range),
			"concat": wdte.GoFunc(Concat),

			"map":       wdte.GoFunc(Map),
			"filter":    wdte.GoFunc(Filter),
			"flatMap":   wdte.GoFunc(FlatMap),
			"enumerate": wdte.GoFunc(Enumerate),

			"collect": wdte.GoFunc(Collect),
			"drain":   wdte.GoFunc(Drain),
			"reduce":  wdte.GoFunc(Reduce),
			//"chain":   wdte.GoFunc(Chain),
			"any": wdte.GoFunc(Any),
			"all": wdte.GoFunc(All),
		},
	}
}
