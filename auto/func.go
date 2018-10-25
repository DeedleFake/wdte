package auto

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/DeedleFake/wdte"
)

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
