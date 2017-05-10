package wdte_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
)

func TestModule(t *testing.T) {
	const test = `
fib n => switch n {
	0 => 0;
	1 => 1;
	default => + (fib (- n 1)) (fib (- n 2));
};

main => (
	range 15
    -> map fib
	  -> print
	;
	[5; 2; (fib 7)]
	  -> map (+ 2)
		-> print
	;
);
`

	m, err := wdte.Parse(strings.NewReader(test), nil)
	if err != nil {
		t.Fatal(err)
	}

	// This could probably be done cleaner, but it works.
	m.Funcs["+"] = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		if len(args) == 1 {
			s := args[0].Call(frame).(wdte.Number)
			return wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
				a := args[0].Call(frame).(wdte.Number)
				return s + a
			})
		}

		a1 := args[0].Call(frame).(wdte.Number)
		a2 := args[1].Call(frame).(wdte.Number)

		return a1 + a2
	})

	m.Funcs["-"] = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		a1 := args[0].Call(frame).(wdte.Number)
		a2 := args[1].Call(frame).(wdte.Number)

		return a1 - a2
	})

	m.Funcs["range"] = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		a := args[0].Call(frame).(wdte.Number)

		r := make(wdte.Array, int(a))
		for i := range r {
			r[i] = wdte.Number(i)
		}

		return r
	})

	m.Funcs["map"] = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		m := args[0].Call(frame)
		return wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
			a := args[0].Call(frame).(wdte.Array)

			r := make(wdte.Array, len(a))
			for i := range r {
				r[i] = m.Call(frame, a[i].Call(frame))
			}
			return r
		})
	})

	m.Funcs["print"] = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		if len(args) < 1 {
			return m.Funcs["print"]
		}

		a := args[0].Call(frame)
		t.Logf("%v", a)
		return a
	})

	//t.Log("Imports:")
	//for i := range m.Imports {
	//	t.Logf("\t%q", i)
	//}

	//t.Log("Funcs:")
	//for f := range m.Funcs {
	//	t.Logf("\t%q", f)
	//}

	m.Funcs["main"].Call(nil)
}
