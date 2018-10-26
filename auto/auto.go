package auto

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/DeedleFake/wdte"
)

var (
	arrayType  = reflect.TypeOf(wdte.Array(nil))
	numberType = reflect.TypeOf(wdte.Number(0))
	stringType = reflect.TypeOf(wdte.String(""))
)

func fromWDTE(w wdte.Func, expected reflect.Type) reflect.Value {
	v := reflect.ValueOf(w)

	switch expected.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return v.Convert(expected)

	case reflect.Array:
		v := v.Convert(arrayType).Interface().(wdte.Array)
		if len(v) != expected.Len() {
			panic(errors.New("array length does not match"))
		}

		t := reflect.ArrayOf(expected.Len(), expected.Elem())
		r := reflect.New(t).Elem()
		for i := 0; i < r.Len(); i++ {
			r.Index(i).Set(fromWDTE(v[i], expected.Elem()))
		}
		return r

	case reflect.Func:
		panic(errors.New("func arguments are not yet supported"))

	case reflect.Map:
		panic(errors.New("map arguments are not yet supported"))

	case reflect.Slice:
		v := v.Convert(arrayType).Interface().(wdte.Array)

		t := reflect.SliceOf(expected.Elem())
		r := reflect.MakeSlice(t, 0, len(v))
		for _, e := range v {
			r = reflect.Append(r, fromWDTE(e, expected.Elem()))
		}
		return r

	case reflect.Struct:
		panic(errors.New("struct arguments are not yet supported"))
	}

	panic(fmt.Errorf("unsupported type: %v", expected))
}

func toWDTE(v reflect.Value) wdte.Func {
	switch v.Kind() {
	case reflect.Bool:
		return wdte.Bool(v.Bool())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		return v.Convert(numberType).Interface().(wdte.Number)

	case reflect.Array, reflect.Slice:
		r := make(wdte.Array, v.Len())
		for i := range r {
			r[i] = toWDTE(v.Index(i))
		}
		return r

	case reflect.Func:
		return Func("<auto>", v.Interface())

	case reflect.Map:
		panic(errors.New("maps are not yet supported"))

	case reflect.Ptr:
		return toWDTE(v.Elem())

	case reflect.String:
		return v.Convert(stringType).Interface().(wdte.String)

	case reflect.Struct:
		panic(errors.New("structs are not yet supported"))
	}

	panic(fmt.Errorf("unsupported type: %T", v))
}
