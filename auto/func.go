package auto

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/DeedleFake/wdte"
)

// Func returns a wdte.Func that wraps the given function f. The
// returned function automatically translates any supported types to
// and from the appropriate types for both arguments and return
// values. The given name is used to name the frame used inside the
// function.
//
// Note that currently this is limited to functions that have a single
// return value, and not all types are supported. Currently supported
// types are as follows:
//
// Arguments:
// * Any number type, including uintptrs.
// * Arrays and slices of supported types. Note that the passed WDTE
//   array's length must match the expected length of the array in the
//   Go function's arguments.
// * Booleans.
// * Strings.
//
// Return types:
// * Any number type, including uintptrs.
// * Arrays and slices of supported types.
// * Booleans.
// * Strings.
// * Pointers to supported types.
// * Functions that are supported by this function. The functions will
//   use a frame with the name "<auto>".
func Func(name string, f interface{}) wdte.Func {
	v := reflect.ValueOf(f)

	t := v.Type()
	if t.Kind() != reflect.Func {
		panic(errors.New("f is not a function"))
	}
	if t.NumOut() != 1 {
		panic(fmt.Errorf("invalid number of returns: %v", t.NumOut()))
	}

	var r wdte.GoFunc
	r = wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		frame = frame.Sub(wdte.ID(name))

		if len(args) < t.NumIn() {
			return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
				return r(frame, append(args, next...)...)
			})
		}

		in := make([]reflect.Value, t.NumIn())
		for i := range in {
			in[i] = fromWDTE(args[i].Call(frame), t.In(i))
		}

		out := v.Call(in)

		return toWDTE(out[0])
	})
	return r
}
