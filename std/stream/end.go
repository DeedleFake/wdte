package stream

import (
	"sort"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/wdteutil"
)

// end is a special value returned by the function provided to new.
type end struct{}

func (end) Call(wdte.Frame, ...wdte.Func) wdte.Func {
	return end{}
}

func (end) Reflect(name string) bool { // nolint
	return name == "End"
}

func (end) String() string {
	return "<end>"
}

// End returns a special value that is returned by the next function
// provided to new when it wants to end the stream.
func End() wdte.Func {
	return end{}
}

// Collect is a WDTE function with the following signature:
//
//    collect s
//
// Iterates through the Stream s, collecting the yielded elements into
// an array. When the Stream ends, it returns the collected array.
func Collect(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("collect")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Collect)
	}

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

// Drain is a WDTE function with the following signature:
//
//    drain s
//
// Drain is the same as Collect, but it simply discards elements as
// they are yielded by the Stream, returning the empty Stream when
// it's done. The main purpose of this function is to allow Map to be
// used as a foreach-style loop without the allocation that Collect
// performs.
func Drain(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("drain")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Drain)
	}

	s := args[0].Call(frame).(Stream)

	last := End()
	for {
		n, ok := s.Next(frame)
		if !ok {
			return last
		}
		if _, ok := n.(error); ok {
			return n
		}

		last = n
	}
}

// Reduce is a WDTE function with the following signatures:
//
//    reduce s i r
//    (reduce r) s i
//    (reduce i r) s
//
// Reduce performs a reduction on the Stream s, resulting in a single
// value, which is returned. i is the initial value for the reduction,
// and r is the reducer. r is expected to have the following
// signature:
//
//    r acc n
//
// r is passed the accumulated value as acc, starting with i, and the
// latest value yielded by the Stream as n. Whatever value r returns
// is used as the next value of acc until the Stream is empty, at
// which point the last value of acc is returned. For example,
//
//    range 5 -> reduce 0 +
//
// returns a summation of the range [0,5).
func Reduce(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("reduce")

	if len(args) < 3 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(Reduce), args...)
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

// Fold is a WDTE function with the following signatures:
//
//    fold s r
//    (fold r) s
//
// Fold is exactly like Reduce, but is uses the first element of the
// Stream s as its initial element, rather than taking an explicit
// one. If there is no first element, it returns End.
func Fold(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("fold")

	if len(args) < 2 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(Fold), args...)
	}

	s := args[0].Call(frame).(Stream)
	cur, ok := s.Next(frame)
	if !ok {
		return End()
	}

	r := args[1]

	return Reduce(frame, s, cur, r)
}

// Extent is a WDTE function with the following signatures:
//
//    extent s n less
//    (extent n less) s
//    (extent less) s n
//    ((extent less) n) s
//
// It drains the Stream s, building up a list of up to n elements
// yielded for which less returns true compared to other elements in
// the list, sorted such that the first element of the list is the
// most less of them. In other words, it returns the n most minimum
// elements using less to perform the compartison. For example,
//
//    range 10 -> extent 3 >
//
// will return [9; 8; 7].
//
// If n is less than 0, there is no limit on the length of the list
// built, meaning that it will contain every element that the Stream
// yields, essentially acting like a sorting variant of collect.
func Extent(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("extent")

	if len(args) < 3 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(Extent), args...)
	}

	s := args[0].Call(frame).(Stream)
	length := args[1].Call(frame).(wdte.Number)
	less := args[2].Call(frame)

	var c func(wdte.Array, int, wdte.Func) wdte.Array
	c = func(extent wdte.Array, i int, f wdte.Func) wdte.Array {
		extent = append(extent[:i], append(wdte.Array{f}, extent[i:]...)...)

		if (length >= 0) && (len(extent) >= int(length)) {
			extent = extent[:int(length)]

			c = func(extent wdte.Array, i int, f wdte.Func) wdte.Array {
				copy(extent[i+1:], extent[i:])
				extent[i] = f
				return extent
			}
		}

		return extent
	}

	sc := int(length)
	if sc < 0 {
		sc = 0
	}

	extent := make(wdte.Array, 0, sc)
	for {
		n, ok := s.Next(frame)
		if !ok {
			return extent
		}
		n = n.Call(frame)

		i := sort.Search(len(extent), func(i int) bool {
			return less.Call(frame, n, extent[i]) == wdte.Bool(true)
		})
		if i < len(extent) {
			extent = c(extent, i, n)
			continue
		}

		if (length < 0) || (len(extent) < int(length)) {
			extent = append(extent, n)
		}
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
//	frame = frame.Sub("call")
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

// Any is a WDTE function with the following signatures:
//
//    any s f
//    (any f) s
//
// It iterates over the Stream s, passing each yielded element to f in
// turn. If any of those calls returns true, then the entire function
// returns true. Otherwise it returns false. It is short-circuiting.
func Any(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("any")

	if len(args) < 2 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(Any), args...)
	}

	s := args[0].Call(frame).(Stream)
	f := args[1].Call(frame)

	for {
		n, ok := s.Next(frame)
		if !ok {
			return wdte.Bool(false)
		}

		if b, ok := f.Call(frame, n).(wdte.Bool); bool(b) && ok {
			return wdte.Bool(true)
		}
	}
}

// All is a WDTE function with the following signatures:
//
//    all s f
//    (all f) s
//
// It iterates over the Stream s, passing each yielded element to f in
// turn. If all of those calls return true, then the entire function
// returns true. Otherwise it returns false. It is short-circuiting.
func All(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("all")

	if len(args) < 2 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(All), args...)
	}

	s := args[0].Call(frame).(Stream)
	f := args[1].Call(frame)

	for {
		n, ok := s.Next(frame)
		if !ok {
			return wdte.Bool(true)
		}

		if b, ok := f.Call(frame, n).(wdte.Bool); bool(!b) || !ok {
			return wdte.Bool(false)
		}
	}
}
