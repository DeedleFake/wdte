package wdte_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
)

func TestModule(t *testing.T) {
	const test = `
'test' => test;

fib n => switch n {
	0 => 0;
	default => + (fib (- n 1)) (fib (- n 2));
};

main => print (fib 5);
`

	m, err := wdte.Parse(strings.NewReader(test), nil)
	if err != nil {
		t.Fatal(err)
	}

	m.Funcs["+"] = wdte.GoFunc(func(scope []wdte.Func, args ...wdte.Func) wdte.Func {
		a1 := args[0].Call(scope)
		a2 := args[1].Call(scope)

		return a1.(wdte.Number) + a2.(wdte.Number)
	})

	m.Funcs["-"] = wdte.GoFunc(func(scope []wdte.Func, args ...wdte.Func) wdte.Func {
		a1 := args[0].Call(scope)
		a2 := args[1].Call(scope)

		return a1.(wdte.Number) - a2.(wdte.Number)
	})

	m.Funcs["print"] = wdte.GoFunc(func(scope []wdte.Func, args ...wdte.Func) wdte.Func {
		a := args[0].Call(scope)
		t.Logf("%v", a)
		return a
	})

	t.Log("Imports:")
	for i := range m.Imports {
		t.Logf("\t%q", i)
	}

	t.Log("Funcs:")
	for f := range m.Funcs {
		t.Logf("\t%q", f)
	}

	m.Funcs["main"].Call(nil)
}
