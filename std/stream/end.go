package stream

import "github.com/DeedleFake/wdte"

// Collect iterates over a stream, putting each of its elements into a
// new array, in the order they are yielded by the stream, before
// returning the array.
//
// If any of the values yielded by the stream are errors, then
// iteration stops and that error is returned instead.
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

// Drain drains the stream, iterating over it much like collect does,
// but discarding each value. When it's finished it returns the
// now-empty stream. If an error is yielded by the stream, then
// iteration stops and that error is returned instead.
//
// The primary purpose of this function is to allow map to be used as
// a type of foreach-style loop without the extra allocation that
// collect performs. For example:
//
//     s.range 5 -> s.map (io.writeln io.stdout) -> s.drain;
func Drain(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Drain)
	}

	frame = frame.WithID("drain")

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

// Any takes two arguments, a stream and a function. It iterates over
// the stream's values, calling the given function on each element. If
// any of the calls return true, than the whole function returns true.
// If it reaches the end of the stream, then it returns
// false.
//
// If given only one argument, it returns a function which checks its
// own argument, a stream, against the function it was originally
// given.
func Any(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Any)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Any(frame, append(next, args...)...)
		})
	}

	frame = frame.WithID("any")

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

// All takes two arguments, a stream and a function. It iterates over
// the stream's values, calling the given function on each element. If
// any of the calls don't return true, than the whole function returns
// false. If it reaches the end of the stream, then it returns
// true.
//
// If given only one argument, it returns a function which checks its
// own argument, a stream, against the function it was originally
// given.
func All(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(All)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return All(frame, append(next, args...)...)
		})
	}

	frame = frame.WithID("all")

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
