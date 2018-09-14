package stream

import (
	"github.com/DeedleFake/wdte"
)

// An array is a stream that iterates over an array.
type array struct {
	a wdte.Array
	i int
}

// New is a WDTE function with the following signature:
//
//    new ...
//
// It returns a Stream that iterates over its arguments.
func New(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(New)
	}

	frame = frame.Sub("new")

	return &array{
		a: wdte.Array(args),
	}
}

func (a *array) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return a
}

func (a *array) Next(frame wdte.Frame) (wdte.Func, bool) { // nolint
	if a.i >= len(a.a) {
		return nil, false
	}

	r := a.a[a.i]
	a.i++
	return r, true
}

func (a *array) String() string { // nolint
	return "<stream>"
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

// Range is a WDTE function with the following signatures:
//
//    range end
//    range start end
//    range start end step
//
// It returns a new Stream which iterates from start to end, stepping
// by step each time. In other words, it's similar to the following
// pseudo Go code
//
//    for i := start; i < end; i += step {
//      yield i
//    }
//
// but with the difference that if step is negative, then the loop
// condition is inverted.
//
// If start is not specified, it is assumed to be 0. If step is not
// specified it is assumed to be 1 if start is greater than or equal
// to end, and -1 if start is less then end.
func Range(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Range)
	}

	frame = frame.Sub("range")

	switch len(args) {
	case 1:
		return &rng{
			m: args[0].Call(frame).(wdte.Number),
			s: 1,
		}

	case 2:
		start := args[0].Call(frame).(wdte.Number)
		end := args[1].Call(frame).(wdte.Number)

		s := wdte.Number(1)
		if start > end {
			s = -1
		}

		return &rng{
			i: start,
			m: end,
			s: s,
		}
	}

	return &rng{
		i: args[0].Call(frame).(wdte.Number),
		m: args[1].Call(frame).(wdte.Number),
		s: args[2].Call(frame).(wdte.Number),
	}
}

func (r *rng) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return r
}

func (r *rng) Next(frame wdte.Frame) (wdte.Func, bool) { // nolint
	if (r.s >= 0) && (r.i >= r.m) {
		return nil, false
	}
	if (r.s < 0) && (r.i <= r.m) {
		return nil, false
	}

	n := r.i
	r.i += r.s
	return n, true
}

func (r *rng) String() string { // nolint
	return "<stream>"
}

// Concat is a WDTE function with the following signatures:
//
//    concat s ...
//    (concat s) ...
//
// It returns a new Stream that yields the values of all of its
// argument Streams in the order that they were given.
func Concat(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Concat)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Concat(frame, append(args, next...)...)
		})
	}

	frame = frame.Sub("concat")

	var i int
	cur := args[0].Call(frame).(Stream)
	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		if i >= len(args) {
			return nil, false
		}

		for {
			n, ok := cur.Next(frame)
			if !ok {
				i++
				if i >= len(args) {
					return nil, false
				}

				cur = args[i].Call(frame).(Stream)
				continue
			}
			return n, ok
		}
	})
}
