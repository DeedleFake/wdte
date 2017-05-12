package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/std/math"
	"github.com/DeedleFake/wdte/std/stream"
	"github.com/gopherjs/gopherjs/js"
)

const initial = `# Welcome to the WDTE playground, a browser based evaluation
# environment for WDTE. This playground's features includes the
# standard function set as well as a number of importable modules.
#
# If you have never seen WDTE before and are completely confused at
# the moment, try reading the overview on the WDTE project's wiki:
# https://github.com/DeedleFake/wdte/wiki
#
# For documentation on the standard function set, see
# https://godoc.org/github.com/DeedleFake/wdte/std
#
# Importable modules:
# * 'math' (https://godoc.org/github.com/DeedleFake/wdte/std/math)
# * 'stream' (https://godoc.org/github.com/DeedleFake/wdte/std/stream)
#
# In addition, a print function is provided which uses the Go fmt
# package to create a string representation of its arguments. This
# string is printed to the output pane and then returned.

'math' => m;
'stream' => s;

memo fib n => switch n {
	== 0 => 0;
	== 1 => 1;
	default => + (fib (- n 1)) (fib (- n 2));
};

main => (
	fib 50 -> print;

	s.range (* m.pi -1) m.pi .5
	-> s.map m.sin
	-> s.collect
	-> print;
);`

type elementWriter struct {
	*js.Object
}

func (e elementWriter) Write(data []byte) (int, error) {
	return e.WriteString(string(data))
}

func (e elementWriter) WriteString(data string) (int, error) {
	t := e.Get("innerHTML").String()
	e.Set("innerHTML", t+data)
	e.Set("scrollTop", e.Get("scrollHeight"))

	return len(data), nil
}

func im(from string) (*wdte.Module, error) {
	switch from {
	case "math":
		return math.Module(), nil
	case "stream":
		return stream.Module(), nil
	}

	return nil, fmt.Errorf("Unknown import: %q", from)
}

func main() {
	document := js.Global.Get("document")
	stdin := document.Call("getElementById", "input")
	stdout := document.Call("getElementById", "output")
	stderr := document.Call("getElementById", "error")

	log.SetFlags(log.Ltime)
	log.SetOutput(&elementWriter{stderr})

	js.Global.Set("run", func(args ...interface{}) *js.Object {
		i := strings.NewReader(stdin.Get("value").String())
		o := &elementWriter{stdout}

		m, err := wdte.Parse(i, wdte.ImportFunc(im))
		if err != nil {
			log.Printf("Failed to parse: %v", err)
			return nil
		}
		std.Insert(m)

		m.Funcs["print"] = wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
			if len(args) == 0 {
				return m.Funcs["print"]
			}

			frame = frame.WithID("print")

			a := make([]interface{}, 0, len(args))
			for _, arg := range args {
				arg = arg.Call(frame)
				if _, ok := arg.(error); ok {
					return arg
				}
				a = append(a, arg)
			}

			str := fmt.Sprint(a...)
			fmt.Fprintln(o, str)
			return wdte.String(str)
		})

		main, ok := m.Funcs["main"]
		if !ok {
			log.Println("No main function found.")
			return nil
		}

		stdout.Set("innerHTML", "")
		if err, ok := main.Call(wdte.F()).(error); ok {
			log.Println(err)
		}

		return nil
	})

	js.Global.Set("clear", func(args ...interface{}) *js.Object {
		stderr.Set("innerHTML", "")
		return nil
	})

	js.Global.Set("main", func(args ...interface{}) *js.Object {
		stdin.Set("value", initial)
		return nil
	})
}
