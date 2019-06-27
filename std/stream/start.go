package stream

import (
	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/wdteutil"
)

// New is a WDTE function with the following signature:
//
//    new initial next
//    (new next) initial
//
// It returns a new Stream that calls next in order to get the next
// element in the stream, passing it first initial and then the
// previous value on each call. The Stream ends when next returns end.
func New(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("new")

	if len(args) < 2 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(New), args...)
	}

	prev := args[0]
	next := args[1].Call(frame)

	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		if _, ok := prev.(end); ok {
			return nil, false
		}

		prev = next.Call(frame, prev)
		if _, ok := prev.(end); ok {
			return nil, false
		}

		return prev, true
	})
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
	frame = frame.Sub("range")

	if len(args) == 0 {
		return wdte.GoFunc(Range)
	}

	// Current index, minimum/maximum value, and step.
	var i, m wdte.Number
	s := wdte.Number(1)

	switch len(args) {
	case 1:
		m = args[0].Call(frame).(wdte.Number)

	case 2:
		i = args[0].Call(frame).(wdte.Number)
		m = args[1].Call(frame).(wdte.Number)

		if i > m {
			s = -1
		}

	default:
		i = args[0].Call(frame).(wdte.Number)
		m = args[1].Call(frame).(wdte.Number)
		s = args[2].Call(frame).(wdte.Number)
	}

	return NextFunc(func(frame wdte.Frame) (wdte.Func, bool) {
		if (s >= 0) && (i >= m) {
			return nil, false
		}
		if (s < 0) && (i <= m) {
			return nil, false
		}

		n := i
		i += s
		return n, true
	})
}

// Concat is a WDTE function with the following signatures:
//
//    concat s ...
//    (concat s) ...
//
// It returns a new Stream that yields the values of all of its
// argument Streams in the order that they were given.
func Concat(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("concat")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Concat), args...)
	}

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
