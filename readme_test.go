package wdte_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/wdteutil"
)

const src = `
let i => import 'some/import/path/or/another';

i.print 3;
+ 5 2 -> i.print;
7 -> + 5 -> i.print;
`

func im(from string) (*wdte.Scope, error) {
	return wdte.S().Map(map[wdte.ID]wdte.Func{
		"print": wdteutil.Func("print", func(v interface{}) interface{} {
			fmt.Println(v)
			return v
		}),
	}), nil
}

func Sum(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("+")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Sum), args...)
	}

	var sum wdte.Number
	for _, arg := range args {
		sum += arg.(wdte.Number)
	}
	return sum
}

func Example() {
	m, err := wdte.Parse(strings.NewReader(src), wdte.ImportFunc(im), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing script: %v\n", err)
		os.Exit(1)
	}

	scope := wdte.S().Add("+", wdte.GoFunc(Sum))

	r := m.Call(wdte.F().WithScope(scope))
	if err, ok := r.(error); ok {
		fmt.Fprintf(os.Stderr, "Error running script: %v\n", err)
		os.Exit(1)
	}
}
