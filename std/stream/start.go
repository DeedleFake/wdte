package stream

import "github.com/DeedleFake/wdte"

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
