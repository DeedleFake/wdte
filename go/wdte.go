package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"syscall/js"
	"time"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/scanner"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/arrays"
	wdteio "github.com/DeedleFake/wdte/std/io"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/rand"
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

type stdout struct {
	io.Writer
}

func (w stdout) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return w
}

func (stdout) String() string {
	return "<writer(stdout)>"
}

type stderr struct {
	io.Writer
}

func (w stderr) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return w
}

func (stderr) String() string {
	return "<writer(stderr)>"
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
		"run": js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			go func() {
				input := args[0].String()
				output := args[1]

				buf := bufPool.Get().(*bytes.Buffer)
				defer func() {
					buf.Reset()
					bufPool.Put(buf)
				}()

				wdteio.Stderr = stderr{buf}
				iomod := wdteio.Scope.Map(map[wdte.ID]wdte.Func{
					"stdin": wdteio.Reader{
						Reader: errorIO("stdin is not supported in the playground"),
					},

					"stdout": stdout{buf},

					"stderr": stderr{buf},
				})

				importer := wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
					switch from {
					case "io":
						return iomod, nil

					case "playground":
						return wdte.S().Map(map[wdte.ID]wdte.Func{
							"wdteVersion": wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
								bi, ok := debug.ReadBuildInfo()
								if !ok {
									return wdte.Error{
										Frame: frame,
										Err:   errors.New("Failed to read build info"),
									}
								}

								for _, dep := range bi.Deps {
									if dep.Path != "github.com/DeedleFake/wdte" {
										continue
									}

									return wdte.String(dep.Version)
								}

								return wdte.Error{
									Frame: frame,
									Err:   errors.New("WDTE's version could not be determined"),
								}
							}),

							"goVersion": wdte.String(runtime.Version()),
						}), nil
					}

					return std.Import(from)
				})

				macros := scanner.MacroMap{
					"raw": func(text string) ([]scanner.Token, error) {
						return []scanner.Token{
							{
								Type: scanner.String,
								Val:  text,
							},
						}, nil
					},
				}

				c, err := wdte.Parse(strings.NewReader(input), importer, macros)
				if err != nil {
					output.Invoke(fmt.Sprintf("Error: Failed to parse input: %v", err))
					return
				}

				frame := std.F()
				frame = frame.WithScope(frame.Scope().Add("io", iomod))

				ctx, cancel := context.WithTimeout(frame.Context(), 5*time.Second)
				defer cancel()

				r := c.Call(frame.WithContext(ctx))
				if err, ok := r.(error); ok {
					fmt.Fprintf(buf, "\n\nError: %v", err)
				}

				output.Invoke(js.Null(), buf.String())
			}()

			return nil
		}),
	})

	select {}
}