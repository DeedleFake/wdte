package wdte_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/std/math"
	"github.com/DeedleFake/wdte/std/stream"
)

func TestModule(t *testing.T) {
	const test = `
'stream' => s;

memo fib n => switch n {
	== 0 => 0;
	== 1 => 1;
	default => + (fib (- n 1)) (fib (- n 2));
};

memo fact n => switch n {
	<= 1 => 1;
	default => - n 1 -> fact -> * n;
};

main => (
	s.range 15
  -> s.map fib
	-> s.collect
	-> print;

	s.new [5; 2; fib 7]
	-> s.map (+ 2)
	-> s.collect
	-> print;

	fact 5 -> print;
);
`

	m, err := new(wdte.Module).Insert(std.Module()).Parse(
		strings.NewReader(test),
		wdte.ImportFunc(func(from string) (*wdte.Module, error) {
			switch from {
			case "stream":
				return stream.Module(), nil
			case "math":
				return math.Module(), nil
			}

			return nil, fmt.Errorf("Unknown import: %q", from)
		}),
	)
	if err != nil {
		t.Fatal(err)
	}

	m.Funcs["print"] = wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		frame = frame.WithID("print")

		if len(args) < 1 {
			return m.Funcs["print"]
		}

		a := args[0].Call(frame)
		if _, ok := a.(error); !ok {
			t.Logf("%v", a)
		}
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

	if err, ok := m.Funcs["main"].Call(wdte.F()).(error); ok {
		t.Fatal(err)
	}
}
