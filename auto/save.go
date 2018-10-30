package auto

import "github.com/DeedleFake/wdte"

// SaveArgs returns a function which, when called, prepends args to
// its argument list.
func SaveArgs(f wdte.Func, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return f
	}

	return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
		return f.Call(frame, append(args, next...)...)
	})
}

// SaveArgsReverse returns a function which, when called, appends args
// to its argument list.
func SaveArgsReverse(f wdte.Func, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return f
	}

	return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
		return f.Call(frame, append(next, args...)...)
	})
}
