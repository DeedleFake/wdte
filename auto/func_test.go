package auto_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/auto"
	"github.com/DeedleFake/wdte/std/stream"
)

func testFunc(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(testFunc)
	}

	return args[0].(wdte.Number) + 1
}

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
			f: func(f func(v int) int) func() int {
				var v int
				return func() int {
					p := v
					v = f(v)
					return p
				}
			},
			args:  []wdte.Func{wdte.GoFunc(testFunc)},
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
		{
			name: "Stream",
			f: func(s stream.Stream) int {
				s.Next(wdte.F())
				v, _ := s.Next(wdte.F())
				return int(v.(wdte.Number))
			},
			args: []wdte.Func{stream.Range(wdte.F(), wdte.Number(10))},
			ret:  wdte.Number(1),
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

func ExampleFunc() {
	c, err := wdte.Parse(strings.NewReader(`add 3 2;`), nil, nil)
	if err != nil {
		panic(err)
	}

	scope := wdte.S().Add("add",
		auto.Func("add", func(n1, n2 int) int {
			return n1 + n2
		}),
	)

	fmt.Println(c.Call(wdte.F().WithScope(scope)))

	// Output: 5
}
