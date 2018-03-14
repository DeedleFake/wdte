package main

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/arrays"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/strings"
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

func im(from string) (*wdte.Scope, error) {
	switch from {
	case "canvas":
		stdout.Get("style").Set("display", "none")
		canvas.Get("style").Set("display", "block")
		out = &elementWriter{stderr}
		return CanvasModule(), nil
	}

	return std.Import(from)
}

func saveToClipboard(str string) {
	document := js.Global.Get("document")

	textarea := document.Call("createElement", "textarea")
	textarea.Set("value", str)

	document.Get("body").Call("appendChild", textarea)
	textarea.Call("select")
	document.Call("execCommand", "copy")
	textarea.Call("blur")
	document.Get("body").Call("removeChild", textarea)
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

		m, err := wdte.Parse(in, wdte.ImportFunc(im))
		if err != nil {
			log.Printf("Failed to parse: %v", err)
			return nil
		}

		var funcs map[wdte.ID]wdte.Func
		funcs = map[wdte.ID]wdte.Func{
			"print": wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
				if len(args) == 0 {
					return funcs["print"]
				}

				frame = frame.Sub("print")

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
			}),
		}

		stdout.Set("innerHTML", "")
		if err, ok := m.Call(wdte.F().WithScope(std.Scope.Map(funcs))).(error); ok {
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

	js.Global.Set("save", func(args ...interface{}) interface{} {
		js.Global.Get("location").Set("hash", url.QueryEscape(stdin.Get("value").String()))
		saveToClipboard(js.Global.Get("location").Get("href").String())
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

		hash := js.Global.Get("location").Get("hash").String()
		if len(hash) != 0 {
			document.Call("querySelector", "#example").Set("value", "Custom")

			src, err := url.QueryUnescape(strings.TrimPrefix(hash, "#"))
			if err != nil {
				stderr.Set("innerHTML", fmt.Sprintf("Failed to unescape code: %v", err))
				return nil
			}

			stdin.Set("value", src)
		}

		return nil
	})
}
