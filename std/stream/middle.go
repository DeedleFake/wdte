package stream

import "github.com/DeedleFake/wdte"

// A mapper is a function that creates stream wrappers that call a
// function on each value yielded by the underlying stream, returning
// the result.
type mapper struct {
	m wdte.Func
}

// Map is a WDTE function with the following signature:
//
//    (map f) s
//
// It returns a Stream which calls f on each element yielded by the
// Stream s, yielding the return values of f in their place.
func Map(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Map)
	}

	return &mapper{m: args[0].Call(frame)}
}

func (m *mapper) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
				n, ok := a.Next(frame)
				if !ok {
					return nil, false
				}

				frame = frame.Sub("map")
				return m.m.Call(frame, n), true
			})
		}
	}

	return m
}

type filter struct {
	f wdte.Func
}

// Filter is a WDTE function with the following signature:
//
//    (filter f) s
//
// It returns a Stream which yields only those values yielded by the
// Stream s that (f value) results in true for.
func Filter(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Filter)
	}

	return filter{f: args[0].Call(frame)}
}

func (f filter) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
				for {
					n, ok := a.Next(frame)
					if !ok {
						return nil, false
					}

					frame = frame.Sub("filter")
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

// FlatMap is a WDTE function with the following signature:
//
//    (flatMap f) s
//
// It's identical to Map with one caveat: If a call to f yields a
// Stream, the elements of that Stream are yielded in turn before
// continuing the iteration of s. In other words,
//
//    range 3 -> flatMap (new 0 1) -> collect
//
// returns
//
//    [0; 1; 0; 1; 0; 1]
func FlatMap(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Map)
	}

	return &flatMapper{m: args[0].Call(frame)}
}

func (m *flatMapper) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	switch len(args) {
	case 1:
		if a, ok := args[0].Call(frame).(Stream); ok {
			var cur Stream
			return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
				frame = frame.Sub("flatMap")

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

// Enumerate is a WDTE function with the following signature:
//
//    enumerate s
//
// It returns a Stream which yields values of the form [i; v] where i
// is the zero-based index of the element v that was yielded by the
// Stream s.
func Enumerate(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Enumerate)
	}

	s := args[0].Call(frame).(Stream)

	var i wdte.Number
	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		frame = frame.Sub("enumerate")

		n, ok := s.Next(frame)
		if !ok {
			return nil, false
		}

		r := wdte.Array{i, n}
		i++
		return r, true
	})
}
