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
		name  string
		f     interface{}
		args  []wdte.Func
		ret   wdte.Func
		calls []wdte.Func
	}{
		{
			name: "Number",
			f: func(v int) int {
				return v + 3
			},
			args: []wdte.Func{wdte.Number(2)},
			ret:  wdte.Number(5),
		},
		{
			name: "Array",
			f: func(a [3]int) [2]int {
				return [...]int{a[0] + a[1], a[1] + a[2]}
			},
			args: []wdte.Func{wdte.Array{wdte.Number(1), wdte.Number(2), wdte.Number(3)}},
			ret:  wdte.Array{wdte.Number(3), wdte.Number(5)},
		},
		{
			name: "Slice",
			f: func(s []int) []int {
				var sum int
				prod := 1
				for _, v := range s {
					sum += v
					prod *= v
				}
				return []int{sum, prod}
			},
			args: []wdte.Func{wdte.Array{wdte.Number(1), wdte.Number(2)}},
			ret:  wdte.Array{wdte.Number(3), wdte.Number(2)},
		},
		{
			name: "Func",
			f: func(v int) func() int {
				return func() int {
					v++
					return v - 1
				}
			},
			args:  []wdte.Func{wdte.Number(0)},
			calls: []wdte.Func{wdte.Number(0), wdte.Number(1), wdte.Number(2)},
		},
		{
			name: "Direct",
			f: func(a wdte.Array) wdte.String {
				return wdte.String(fmt.Sprint(a))
			},
			args: []wdte.Func{wdte.Array{wdte.String("a"), wdte.String("test")}},
			ret:  wdte.String("[a; test]"),
		},
		{
			name: "Bool",
			f: func(v bool) bool {
				return !v
			},
			args: []wdte.Func{wdte.Bool(false)},
			ret:  wdte.Bool(true),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := auto.Func(test.name, test.f)
			r := f.Call(wdte.F(), test.args...)

			if (test.ret != nil) && !reflect.DeepEqual(r, test.ret) {
				t.Errorf("Got %#v", r)
				t.Errorf("Expected %#v", test.ret)
			}

			if len(test.calls) > 0 {
				for _, ret := range test.calls {
					n := r.Call(wdte.F())
					if !reflect.DeepEqual(n, ret) {
						t.Errorf("Got %#v", n)
						t.Errorf("Expected %#v", ret)
					}
				}
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
