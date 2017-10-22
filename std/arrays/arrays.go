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

func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"at": wdte.GoFunc(At),
		},
	}
}
