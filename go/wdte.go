package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/arrays"
	wdteio "github.com/DeedleFake/wdte/std/io"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/strings"
)

type errorIO string

func (err errorIO) Read([]byte) (int, error) {
	return 0, errors.New(string(err))
}

func (err errorIO) Write([]byte) (int, error) {
	return 0, errors.New(string(err))
}

func (err errorIO) Error() string {
	return string(err)
}

func main() {
	var (
		bufPool = &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		}
	)

	js.Global().Set("WDTE", map[string]interface{}{
		"run": js.NewCallback(func(args []js.Value) {
			input := args[0].String()
			output := args[1]

			buf := bufPool.Get().(*bytes.Buffer)
			defer func() {
				buf.Reset()
				bufPool.Put(buf)
			}()

			c, err := wdte.Parse(strings.NewReader(input), wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
				switch from {
				case "io":
					return wdteio.Scope.Map(map[wdte.ID]wdte.Func{
						"stdin": wdteio.Reader{
							Reader: errorIO("stdin is not supported in the playground"),
						},

						"stdout": wdteio.Writer{
							Writer: buf,
						},

						"stderr": wdteio.Writer{
							Writer: buf,
						},
					}), nil
				}

				return std.Import(from)
			}))
			if err != nil {
				output.Invoke(fmt.Sprintf("Error: Failed to parse input: %v", err))
				return
			}

			frame := std.F()

			ctx, cancel := context.WithTimeout(frame.Context(), 5*time.Second)
			defer cancel()

			r := c.Call(frame.WithContext(ctx))
			if err, ok := r.(error); ok {
				fmt.Fprintf(buf, "\n\nError: %v", err)
			}

			output.Invoke(js.Null(), buf.String())
		}),
	})

	select {}
}
