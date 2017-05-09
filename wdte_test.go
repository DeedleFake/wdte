package wdte_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
)

func TestModule(t *testing.T) {
	const test = `
'test' => test;

add x y => print (+ x y;);

main => add 3 -5;
`

	m, err := wdte.Parse(strings.NewReader(test), nil)
	if err != nil {
		t.Fatal(err)
	}

	m.Funcs["+"] = wdte.GoFunc(func(args ...wdte.Func) wdte.Func {
		a1 := args[0].Call()
		a2 := args[1].Call()

		return a1.(wdte.Number) + a2.(wdte.Number)
	})

	m.Funcs["print"] = wdte.GoFunc(func(args ...wdte.Func) wdte.Func {
		a := args[0].Call()
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

	m.Funcs["main"].Call()
}
