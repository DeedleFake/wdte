package stream

import (
	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/auto"
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
	for {
		n, ok := s.Next(frame)
		if !ok {
			return s
		}
		if _, ok := n.(error); ok {
			return n
		}
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
		return auto.SaveArgsReverse(wdte.GoFunc(Reduce), args...)
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
		return auto.SaveArgsReverse(wdte.GoFunc(Fold), args...)
	}

	s := args[0].Call(frame).(Stream)
	cur, ok := s.Next(frame)
	if !ok {
		return End()
	}

	r := args[1]

	return Reduce(frame, s, cur, r)
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
		return auto.SaveArgsReverse(wdte.GoFunc(Any), args...)
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
		return auto.SaveArgsReverse(wdte.GoFunc(All), args...)
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
