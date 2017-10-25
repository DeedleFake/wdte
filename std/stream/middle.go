package stream

import "github.com/DeedleFake/wdte"

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

type flatMapper struct {
	m wdte.Func
}

// FlatMap works similarly to map, but if the mapping function returns
// an stream, the contents of that stream are substituted for the
// values of the stream, rather than the stream itself being yielded.
// For example,
//
//     s.range 3 -> s.flatMap (s.new 0 1) -> s.collect
//
// returns
//
//     [0; 1; 0; 1; 0; 1]
func FlatMap(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Map)
	}

	return &flatMapper{m: args[0].Call(frame)}
}

func (m *flatMapper) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			var cur Stream
			return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
				frame = frame.WithID("flatMap")

				for {
					if cur != nil {
						n, ok := cur.Next(frame)
						if ok {
							return n, true
						}

						cur = nil
					}

					n, ok := a.Next(frame)
					if !ok {
						return nil, false
					}

					r := m.m.Call(frame, n).Call(frame)
					if r, ok := r.(Stream); ok {
						cur = r
						continue
					}

					return r, true
				}
			})
		}
	}

	return m
}

// Enumerate takes a stream and returns a new stream that yields
// values from the original stream wrapped in an array such that each
// value yielded is in the format
//
//     [i; v]
//
// where i is the index of the value yielded and v is the original
// value.
func Enumerate(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Enumerate)
	}

	s := args[0].Call(frame).(Stream)

	var i wdte.Number
	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		frame = frame.WithID("enumerate")

		n, ok := s.Next(frame)
		if !ok {
			return nil, false
		}

		r := wdte.Array{i, n}
		i++
		return r, true
	})
}
