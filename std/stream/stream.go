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

// An array is a stream that iterates over an array.
type array struct {
	a wdte.Array
	i int
}

// New returns a new stream. If given one argument and that argument
// is an array, it iterates over the values of the array. If given
// more than one argument or the first argument is not an array, it
// iterates over its arguments.
func New(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(New)

	case 1:
		if a1, ok := args[0].Call(frame).(wdte.Array); ok {
			frame = frame.WithID("new")
			return &array{a: a1}
		}
	}

	frame = frame.WithID("new")
	return &array{a: args}
}

func (a *array) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return a
}

func (a *array) Next(frame wdte.Frame) (wdte.Func, bool) {
	if a.i >= len(a.a) {
		return nil, false
	}

	r := a.a[a.i]
	a.i++
	return r, true
}

// An rng (range) is a stream that yields successive numbers.
type rng struct {
	// i is the next number to return.
	i wdte.Number

	// m is the number to stop at, exclusive.
	m wdte.Number

	// s is the amount to increment every time Next() is called.
	s wdte.Number
}

// Range returns a stream that yields successive numbers. If given a
// single argument, the range yielded is [0, args[0]). If given two
// arguments, the range is [args[0], args[1]). If given three
// arguments, the range is the same as with two arguments, but the
// step in between numbers yielded is the third argument, rather than
// 1.
func Range(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Range)
	}

	frame = frame.WithID("range")

	switch len(args) {
	case 1:
		return &rng{
			m: args[0].Call(frame).(wdte.Number),
			s: 1,
		}

	case 2:
		return &rng{
			i: args[0].Call(frame).(wdte.Number),
			m: args[1].Call(frame).(wdte.Number),
			s: 1,
		}
	}

	return &rng{
		i: args[0].Call(frame).(wdte.Number),
		m: args[1].Call(frame).(wdte.Number),
		s: args[2].Call(frame).(wdte.Number),
	}
}

func (r *rng) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return r
}

func (r *rng) Next(frame wdte.Frame) (wdte.Func, bool) {
	if r.i >= r.m {
		return nil, false
	}

	n := r.i
	r.i += r.s
	return n, true
}

// A mapper is a function that creates stream wrappers that call a
// function on each value yielded by the underlying stream, returning
// the result.
type mapper struct {
	m wdte.Func
}

// Map returns a function that takes a stream and wraps the stream in
// a new stream that calls the function originally given to Map on
// each element before passing it on.
//
// Wow, that sound horrible. It's not too bad, though. It works like
// this. Call Map with some function `f` and get a new function. Then,
// call that returned function on a stream to get a new stream that
// calls `f` on each element when Next is called.
func Map(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Map)
	}

	return &mapper{m: args[0].Call(frame)}
}

func (m *mapper) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
				n, ok := a.Next(frame)
				if !ok {
					return nil, false
				}

				frame = frame.WithID("map")
				return m.m.Call(frame, n), true
			})
		}
	}

	return m
}

type filter struct {
	f wdte.Func
}

// Filter returns a stream which yields values from a previous stream
// that are passed through a filter. For example,
//
//     s.range 5 -> s.filter (<= 2) -> s.collect
//
// returns [0; 1; 2].
func Filter(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Filter)
	}

	return filter{f: args[0].Call(frame)}
}

func (f filter) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
				for {
					n, ok := a.Next(frame)
					if !ok {
						return nil, false
					}

					frame = frame.WithID("filter")
					if f.f.Call(frame, n) == wdte.Bool(true) {
						return n, true
					}
				}
			})
		}
	}

	return f
}

// Collect converts a stream into an array.
func Collect(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Collect)
	}

	frame = frame.WithID("collect")

	a, ok := args[0].Call(frame).(Stream)
	if !ok {
		return args[0]
	}

	r := wdte.Array{}
	for {
		n, ok := a.Next(frame)
		if !ok {
			break
		}
		if _, ok := n.(error); ok {
			return n
		}

		r = append(r, n.Call(frame))
	}

	return r
}

// Reduce reduces a stream to a single value using a reduction
// function. For example,
//
//     s.range 5 -> s.reduce 0 +
//
// will yield a summation of the number 0 through 4, inclusive.
//
// Reduce takes three arguments: A stream, an initial value, and a
// function. It first calls the reduction function with the initial
// value given and the first value from the stream. It then continues
// iterating over every value in the stream, passing both the previous
// output of the reduction function and the value the stream yielded
// until the stream is empty, at which point it returns the most
// recent output from the reduction function.
func Reduce(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("reduce")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Reduce)
	case 1, 2:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Reduce(frame, append(next, args...)...)
		})
	}

	s := args[0].Call(frame).(Stream)
	cur := args[1].Call(frame)
	r := args[2].Call(frame)

	for {
		n, ok := s.Next(frame)
		if !ok {
			return cur
		}

		cur = r.Call(frame, cur, n)
	}
}

// TODO: Implement this. It should be able to essentially insert its
// own output into another chain, so that
//
//    start -> (s.range 5 -> s.map f -> s.chain) -> end
//
// is possible.
//func Chain(frame wdte.Frame, args ...wdte.Func) wdte.Func {
//	switch len(args) {
//	case 0:
//		return wdte.GoFunc(Chain)
//	}
//
//	frame = frame.WithID("call")
//
//	a, ok := args[0].Call(frame).(Stream)
//	if !ok {
//		return args[0]
//	}
//
//	var prev wdte.Func
//	for {
//		n, ok := a.Next(frame)
//		if !ok {
//			break
//		}
//		if _, ok := n.(error); ok {
//			return n
//		}
//
//		if prev != nil {
//			n = n.Call(frame).Call(frame, prev.Call(frame))
//		}
//
//		prev = n
//	}
//
//	return prev
//}

// Module returns a module for easy importing into an actual script.
// The imported functions have the same names as the functions in this
// package, except that the first letter is lowercase.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"new":   wdte.GoFunc(New),
			"range": wdte.GoFunc(Range),

			"map":    wdte.GoFunc(Map),
			"filter": wdte.GoFunc(Filter),

			"collect": wdte.GoFunc(Collect),
			"reduce":  wdte.GoFunc(Reduce),
			//"chain":   wdte.GoFunc(Chain),
		},
	}
}
