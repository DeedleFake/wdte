// Package arrays contains functions for manipulating arrays.
package arrays

import (
	"sort"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/wdteutil"
)

// Concat is a WDTE function with the following signatures:
//
//    concat array ...
//    (concat array) ...
//
// Returns an array containing the concatonation of all of its arguments, all of which should be arrays, in the order that they were passed to it. For example,
//
//    concat [3; 6] [2; 5]
//
// returns
//
//    [3; 6; 2; 5]
func Concat(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("concat")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Concat), args...)
	}

	array := args[0].(wdte.Array)
	for _, arg := range args[1:] {
		array = append(array[:len(array):len(array)], arg.(wdte.Array)...)
	}
	return array
}

func sorter(sortFunc func(interface{}, func(int, int) bool)) (f wdte.GoFunc) {
	return func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		if len(args) < 2 {
			return wdteutil.SaveArgs(f, args...)
		}

		var array wdte.Array
		var less wdte.Func
		switch a := args[0].(type) {
		case wdte.Array:
			array = a
			less = args[1]
		default:
			less = a
			array = args[1].(wdte.Array)
		}

		type errorFunc interface {
			wdte.Func
			error
		}

		var err errorFunc
		array = append(wdte.Array{}, array...)
		sortFunc(array, func(i1, i2 int) bool {
			if err != nil {
				return false
			}

			r := less.Call(frame, array[i1], array[i2])
			if e, ok := r.(errorFunc); ok {
				err = e
				return false
			}

			return r == wdte.Bool(true)
		})
		if err != nil {
			return err
		}

		return array
	}
}

// Sort is a WDTE function with the following signatures:
//
//    sort array less
//    (sort array) less
//    sort less array
//    (sort less) array
//
// Returns a sorted copy of the given array sorted using the given
// less function. The less function should take two arguments and
// return true if the first argument should be sorted earlier in the
// array then the second. Unlike sortStable, the relative positions of
// equal elements are undefined in the new array.
func Sort(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("sort")

	return sorter(sort.Slice).Call(frame, args...)
}

// SortStable is a WDTE function with the following signatures:
//
//    sortStable array less
//    (sortStable array) less
//    sortStable less array
//    (sortStable less) array
//
// Returns a sorted copy of the given array sorted using the given
// less function. The less function should take two arguments and
// return true if the first argument should be sorted earlier in the
// array then the second. Unlike sort, the relative positions of equal
// elements are preserved.
func SortStable(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("sortStable")

	return sorter(sort.SliceStable).Call(frame, args...)
}

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
	frame = frame.Sub("stream")

	if len(args) == 0 {
		return wdte.GoFunc(Stream)
	}

	return &streamer{a: args[0].(wdte.Array)}
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

func (a *streamer) String() string {
	return "<stream>"
}

func (a *streamer) Reflect(name string) bool {
	return name == "Stream"
}

// Scope is a scope containing the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"concat":     wdte.GoFunc(Concat),
	"sort":       wdte.GoFunc(Sort),
	"sortStable": wdte.GoFunc(SortStable),
	"stream":     wdte.GoFunc(Stream),
})

func init() {
	std.Register("arrays", Scope)
}
