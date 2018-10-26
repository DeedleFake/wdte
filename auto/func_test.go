package auto_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/auto"
)

func TestFunc(t *testing.T) {
	tests := []struct {
		name string
		f    interface{}
		args []wdte.Func
		ret  wdte.Func
	}{
		{
			name: "Number",
			f: func(v int) int {
				return v + 3
			},
			args: []wdte.Func{wdte.Number(2)},
			ret:  wdte.Number(5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := auto.Func(test.name, test.f)
			r := f.Call(wdte.F(), test.args...)
			if !reflect.DeepEqual(r, test.ret) {
				t.Errorf("Got %#v", r)
				t.Errorf("Expected %#v", test.ret)
			}
		})
	}
}

func TestFuncPartial(t *testing.T) {
	f := auto.Func("test", func(v1, v2 int) string {
		return fmt.Sprintf("(%v, %v)", v1, v2)
	})

	r := f.Call(wdte.F(), wdte.Number(2))
	if n, ok := r.(wdte.String); ok {
		t.Fatalf("Function returned string (%q) early", n)
	}

	r = r.Call(wdte.F(), wdte.Number(1))
	if r != wdte.String("(2, 1)") {
		t.Errorf("Got %#v", r)
		t.Errorf("Expected %q", wdte.String("(2, 1)"))
	}
}
