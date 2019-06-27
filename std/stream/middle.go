package stream

import (
	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/wdteutil"
)

// Map is a WDTE function with the following signature:
//
//    (map f) s
//
// It returns a Stream which calls f on each element yielded by the
// Stream s, yielding the return values of f in their place.
func Map(frame wdte.Frame, args ...wdte.Func) (mapper wdte.Func) {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Map)
	}

	f := args[0].Call(frame)

	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		switch len(args) {
		case 0:
			return mapper
		}

		s := args[0].Call(frame).(Stream)

		return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
			n, ok := s.Next(frame)
			if !ok {
				return nil, false
			}

			return f.Call(frame.Sub("map"), n), true
		})
	})
}

// Filter is a WDTE function with the following signature:
//
//    (filter f) s
//
// It returns a Stream which yields only those values yielded by the
// Stream s that (f value) results in true for.
func Filter(frame wdte.Frame, args ...wdte.Func) (filter wdte.Func) {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Filter)
	}

	f := args[0].Call(frame)

	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		switch len(args) {
		case 0:
			return filter
		}

		s := args[0].Call(frame).(Stream)

		return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
			for {
				n, ok := s.Next(frame)
				if !ok {
					return nil, false
				}

				if f.Call(frame.Sub("filter"), n) == wdte.Bool(true) {
					return n, true
				}
			}
		})
	})
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
func FlatMap(frame wdte.Frame, args ...wdte.Func) (mapper wdte.Func) {
	switch len(args) {
	case 0:
		return wdte.GoFunc(FlatMap)
	}

	f := args[0].Call(frame)

	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		switch len(args) {
		case 0:
			return mapper
		}

		s := args[0].Call(frame).(Stream)

		var cur Stream
		return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
			for {
				if cur != nil {
					n, ok := cur.Next(frame)
					if ok {
						return n, true
					}

					cur = nil
				}

				n, ok := s.Next(frame)
				if !ok {
					return nil, false
				}

				r := f.Call(frame, n).Call(frame)
				if r, ok := r.(Stream); ok {
					cur = r
					continue
				}

				return r, true
			}
		})
	})
}

// Enumerate is a WDTE function with the following signature:
//
//    enumerate s
//
// It returns a Stream which yields values of the form [i; v] where i
// is the zero-based index of the element v that was yielded by the
// Stream s.
func Enumerate(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("enumerate")

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

// Repeat is a WDTE function with the following signature:
//
//    repeat s
//
// Repeat returns a Stream that buffers the elements from the Stream
// s. Once s has ended, the Stream starts repeating from the begnning
// of the buffer, looping infinitely. In other words,
//
//    range 3 -> repeat
//
// will yield the sequence (0, 1, 2) repeatedly with no end.
//
// Repeat is most useful used with Limit. When combining the two, note
// that Limit limits individual elements, not repetitions, so the
// number passed to Limit should be multiplied properly if the client
// wants to limit to a specific number of repetitions.
func Repeat(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("repeat")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Repeat)
	}

	s := args[0].Call(frame).(Stream)

	var buf []wdte.Func
	loop := -1
	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		frame = frame.Sub("repeat")

	exit:
		if loop >= 0 {
			n := buf[loop]

			loop++
			loop %= len(buf)

			return n, true
		}

		n, ok := s.Next(frame)
		if !ok {
			loop = 0
			goto exit
		}

		buf = append(buf, n)
		return n, true
	})
}

// Limit is a WDTE function with the following signature:
//
//    (limit n) s
//
// Limit returns a Stream that stops after a maximum of n elements
// from s have been yielded.
func Limit(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("limit")

	if len(args) == 0 {
		return wdte.GoFunc(Limit)
	}

	n := args[0].Call(frame).(wdte.Number)

	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		frame = frame.Sub("limit")

		s := args[0].Call(frame).(Stream)

		return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
			frame = frame.Sub("limit")

			if n <= 0 {
				return nil, false
			}
			n--

			return s.Next(frame)
		})
	})
}

// Zip is a WDTE function with the following signatures:
//
//    (zip s1) ...
//    zip ...
//
// Zip returns a Stream which yields the streams that it is given
// simultaneuously as arrays. In other words,
//
//    zip (a.stream [1; 2; 3]) (a.stream ['a'; 'b'; 'c'])
//
// will yield
//
//    [1; 'a']
//    [2; 'b']
//    [3; 'c']
//
// The order of the yielded arrays matches the order that the streams
// are given in. If one of the streams ends before the other ones, End
// will be yielded for that stream after that point.
func Zip(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("zip")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Zip), args...)
	}

	streams := make([]Stream, 0, len(args))
	for _, arg := range args {
		streams = append(streams, arg.Call(frame).(Stream))
	}

	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		frame = frame.Sub("zip")

		r := make(wdte.Array, 0, len(streams))
		var more bool
		for _, stream := range streams {
			n, ok := stream.Next(frame)
			if !ok {
				r = append(r, end{})
				continue
			}

			r = append(r, n)
			more = true
		}

		return r, more
	})
}
