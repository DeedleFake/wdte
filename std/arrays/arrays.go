package arrays

import "github.com/DeedleFake/wdte"

// At returns the element at the index of the first argument specified
// by the second argument. In other words,
//
//     at a i
//
// is the equivalent of
//
//     a[i]
//
// in a C-style language.
//
// If only given one argument, it returns a function which returns the
// element at that index of its own argument.
func At(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("at")

	switch len(args) {
	case 0:
		return wdte.GoFunc(At)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return At(frame, append(next, args...)...)
		})
	}

	a := args[0].Call(frame).(wdte.Array)
	i := args[1].Call(frame).(wdte.Number)

	return a[int(i)]
}

// Len returns the length of an array.
func Len(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("len")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Len)
	}

	a := args[0].Call(frame).(wdte.Array)

	return wdte.Number(len(a))
}

// A streamer is a stream that iterates over an array.
type streamer struct {
	a wdte.Array
	i int
}

// Stream returns a stream that iterates over a given array.
func Stream(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("stream")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Stream)
	}

	return &streamer{a: args[0].Call(frame).(wdte.Array)}
}

func (a *streamer) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return a
}

func (a *streamer) Next(frame wdte.Frame) (wdte.Func, bool) {
	if a.i >= len(a.a) {
		return nil, false
	}

	r := a.a[a.i]
	a.i++
	return r, true
}

// Module returns a module for easy importing into an actual script.
// The imported functions have the same names as the functions in this
// package, except that the first letter is lowercase.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"at":  wdte.GoFunc(At),
			"len": wdte.GoFunc(Len),

			"stream": wdte.GoFunc(Stream),
		},
	}
}
