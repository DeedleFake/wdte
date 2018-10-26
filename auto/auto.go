package auto

import (
	"fmt"
	"reflect"

	"github.com/DeedleFake/wdte"
)

var (
	arrayType = reflect.TypeOf(wdte.Array(nil))
)

func fromWDTE(w wdte.Func, expected reflect.Type) reflect.Value {
	v := reflect.ValueOf(w)

	switch expected.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return v.Convert(expected)

	case reflect.Array:
		v := v.Convert(arrayType).Interface().(wdte.Array)
		if len(v) != expected.Len() {
			panic("array length does not match")
		}

		t := reflect.ArrayOf(expected.Len(), expected.Elem())
		r := reflect.New(t).Elem()
		for i := 0; i < r.Len(); i++ {
			r.Index(i).Set(fromWDTE(v[i], expected.Elem()))
		}
		return r

	case reflect.Func:
		panic("func arguments are not yet supported")

	case reflect.Map:
		panic("map arguments are not yet supported")

	case reflect.Slice:
		v := v.Convert(arrayType).Interface().(wdte.Array)

		t := reflect.SliceOf(expected.Elem())
		r := reflect.MakeSlice(t, 0, len(v))
		for _, e := range v {
			r = reflect.Append(r, fromWDTE(e, expected.Elem()))
		}
		return r

	case reflect.Struct:
		panic("struct arguments are not yet supported")
	}

	panic(fmt.Errorf("unexpected kind: %v", expected.Kind()))
}

func toWDTE(v reflect.Value) wdte.Func {
	panic("Not implemented.")
}
