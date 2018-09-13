package arrays

import (
	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

// A streamer is a stream that iterates over an array.
type streamer struct {
	a wdte.Array
	i int
}

// Stream is a WDTE function with the following signature:
//
//    stream a
//
// Returns a stream.Stream that iterates over the array a.
func Stream(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	switch len(args) {
	case 0:
		return wdte.GoFunc(Stream)
	}

	return &streamer{a: args[0].Call(frame).(wdte.Array)}
}

func (a *streamer) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return a
}

func (a *streamer) Next(frame wdte.Frame) (wdte.Func, bool) { // nolint
	if a.i >= len(a.a) {
		return nil, false
	}

	r := a.a[a.i]
	a.i++
	return r, true
}

// Scope is a scope containing the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"stream": wdte.GoFunc(Stream),
})

func init() {
	std.Register("arrays", Scope)
}
