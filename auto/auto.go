package auto

import (
	"fmt"
	"reflect"

	"github.com/DeedleFake/wdte"
)

func fromWDTE(w wdte.Func, expected reflect.Type) reflect.Value {
	v := reflect.ValueOf(w)

	switch expected.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		return v.Convert(expected)

	case reflect.Array, reflect.Slice:

	case reflect.Func:

	case reflect.Map:

	case reflect.Struct:
	}

	panic(fmt.Errorf("unexpected kind: %v", expected.Kind()))
}

func toWDTE(v reflect.Value) wdte.Func {
	panic("Not implemented.")
}
