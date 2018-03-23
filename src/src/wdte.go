package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/arrays"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/strings"
	"github.com/gopherjs/gopherjs/js"
)

var (
	exports = js.Module.Get("exports")
)

func PrintFunc(w io.Writer) (print wdte.Func) {
	print = wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		if len(args) == 0 {
			return print
		}

		frame = frame.Sub("print")

		arg := args[0].Call(frame)
		fmt.Fprintln(w, arg)
		return arg
	})
	return print
}

func main() {
	var (
		bufPool = &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		}
	)

	exports.Set("eval", func(input string) string {
		buf := bufPool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufPool.Put(buf)
		}()

		c, err := wdte.Parse(strings.NewReader(input), std.Import)
		if err != nil {
			return fmt.Sprintf("Error: Failed to parse input:\n\n%v", err)
		}

		scope := std.Scope.Add("print", PrintFunc(buf))
		r := c.Call(wdte.F().WithScope(scope))
		if err, ok := r.(error); ok {
			fmt.Fprintf(buf, "\n\nError: %v", err)
		}

		return buf.String()
	})
}
