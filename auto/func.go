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
// return value.
//
// Unrecognized types are passed through with an attempted conversion,
// allowing a function to, for example, take a stream.Stream as an
// argument. Similarly, if an expected type of an argument is the
// exact type of the value passed, the value is passed through
// directly. If a return value's type implements wdte.Func, it is also
// passed directly. Types with special handling are as follows:
//
// Arguments:
//    * Arrays and slices. Note that the passed WDTE array's length
//      must match the expected length of the array in the Go
//      function's arguments.
//
// Return types:
//    * Arrays and slices.
//    * Pointers.
//    * Functions that are supported by this function. The functions
//      will use a frame with the name "<auto>".
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
			return SaveArgs(r, args...)
		}

		in := make([]reflect.Value, t.NumIn())
		for i := range in {
			in[i] = fromWDTE(frame, args[i].Call(frame), t.In(i))
		}

		out := v.Call(in)

		return toWDTE(out[0])
	})
	return r
}

// FromFunc does the opposite of Func, returning a Go function with
// the signature given by expected that, when called, calls w with
// frame. Type conversions are handled the same as in Func, but in
// reverse, such that the return type stipulations in Func apply to
// the arguments to w, and vice versa for the return value of w. The
// requested function type must return exactly zero or one values.
func FromFunc(frame wdte.Frame, w wdte.Func, expected reflect.Type) reflect.Value {
	if expected.Kind() != reflect.Func {
		panic(errors.New("expected is not a function type"))
	}
	if expected.NumOut() > 1 {
		panic(fmt.Errorf("invalid number of returns: %v", expected.NumOut()))
	}

	return reflect.MakeFunc(expected, func(args []reflect.Value) []reflect.Value {
		wargs := make([]wdte.Func, 0, len(args))
		for _, arg := range args {
			wargs = append(wargs, toWDTE(arg))
		}

		fmt.Println(w)
		r := w.Call(frame, wargs...)
		if expected.NumOut() == 0 {
			// TODO: Should zero return values be allowed?
			return nil
		}
		return []reflect.Value{fromWDTE(frame, r, expected.Out(0))}
	})
}
