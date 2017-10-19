package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/gopherjs/gopherjs/js"
)

var (
	document *js.Object

	stdin  *js.Object
	stdout *js.Object
	stderr *js.Object

	canvas    *js.Object
	canvasCtx *js.Object

	example *js.Object
)

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

var (
	in  io.Reader
	out io.Writer
)

func im(from string) (*wdte.Module, error) {
	switch from {
	case "canvas":
		stdout.Get("style").Set("display", "none")
		canvas.Get("style").Set("display", "block")
		out = &elementWriter{stderr}
		return CanvasModule(), nil

	case "io", "io/file":
		return nil, fmt.Errorf("%q is disabled in the playground", from)
	}

	return std.Import(from)
}

func main() {
	document = js.Global.Get("document")

	stdin = document.Call("getElementById", "input")
	stdout = document.Call("getElementById", "output")
	stderr = document.Call("getElementById", "error")

	canvas = document.Call("getElementById", "canvas")
	canvasCtx = canvas.Call("getContext", "2d")

	example = document.Call("getElementById", "example")

	log.SetFlags(log.Ltime)
	log.SetOutput(&elementWriter{stderr})

	js.Global.Set("run", func(args ...interface{}) interface{} {
		stdout.Get("style").Set("display", "block")
		canvas.Get("style").Set("display", "none")

		in = strings.NewReader(stdin.Get("value").String())
		out = &elementWriter{stdout}

		canvasCtx.Set("fillStyle", "white")
		canvasCtx.Call("fillRect", 0, 0, 640, 480)

		m, err := new(wdte.Module).Insert(std.Module()).Parse(in, wdte.ImportFunc(im))
		if err != nil {
			log.Printf("Failed to parse: %v", err)
			return nil
		}

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
			fmt.Fprintln(out, str)
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

	js.Global.Set("clear", func(args ...interface{}) interface{} {
		stderr.Set("innerHTML", "")
		return nil
	})

	js.Global.Set("changeExample", func(args ...interface{}) interface{} {
		stdin.Set("value", examples[example.Get("value").String()])
		return nil
	})

	js.Global.Set("main", func(args ...interface{}) interface{} {
		stdin.Set("value", examples["fib"])

		document.Call("querySelector", ".tab").Call(
			"addEventListener",
			"keydown",
			js.MakeFunc(func(this *js.Object, args []*js.Object) interface{} {
				if args[0].Get("keyCode").Int() != 9 {
					return nil
				}

				start := this.Get("selectionStart").Int()
				end := this.Get("selectionEnd").Int()

				target := args[0].Get("target")
				target.Set("value", target.Get("value").String()[:start]+"\t"+target.Get("value").String()[end:])

				this.Set("selectionStart", start+1)
				this.Set("selectionEnd", start+1)

				args[0].Call("preventDefault")

				return nil
			}),
		)

		return nil
	})
}
