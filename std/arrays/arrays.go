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

// Stream returns a stream that iterates over a given array.
func Stream(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("stream")

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

// S returns a top-level scope containing the various functions in
// this package.
func S() *wdte.Scope {
	return wdte.S().Map(map[wdte.ID]wdte.Func{
		"stream": wdte.GoFunc(Stream),
	})
}

func init() {
	std.Register("arrays", S())
}
