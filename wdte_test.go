package wdte_test

import (
	"bytes"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/std/io"
)

type test struct {
	name string

	script string
	im     wdte.Importer

	args []wdte.Func
	ret  wdte.Func

	in  string
	out string
	err string
}

func runTests(t *testing.T, tests []test) {
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer

			im := test.im
			if im == nil {
				im = wdte.ImportFunc(func(from string) (*wdte.Module, error) {
					switch from {
					case "io":
						m := io.Module()
						m.Funcs["stdin"] = io.Reader{strings.NewReader(test.in)}
						m.Funcs["stdout"] = io.Writer{&stdout}
						m.Funcs["stderr"] = io.Writer{&stderr}
						return m, nil
					}

					return std.Import(from)
				})
			}

			m, err := std.Module().Parse(strings.NewReader(test.script), im)
			if err != nil {
				t.Fatalf("Failed to parse script: %v", err)
			}

			main, ok := m.Funcs["main"]
			if !ok {
				t.Fatal("No main function.")
			}

			ret := main.Call(wdte.F(), test.args...)
			switch test.ret {
			case nil:
				if err, ok := ret.(error); ok {
					t.Errorf("Return: Got an error: %v", err)
				}

			default:
				switch ret := ret.(type) {
				case wdte.Comparer:
					if c, _ := ret.Compare(test.ret); c != 0 {
						t.Errorf("Return:\n\tExpected %#v\n\tGot %#v\n\t\t%v", test.ret, ret, ret)
					}

				default:
					if !reflect.DeepEqual(ret, test.ret) {
						t.Errorf("Return:\n\tExpected %#v\n\tGot %#v\n\t\t%v", test.ret, ret, ret)
					}
				}
			}

			if out := stdout.String(); out != test.out {
				t.Errorf("Stdout:\n\tExpected %q\n\tGot %q", test.out, out)
			}
			if err := stderr.String(); err != test.err {
				t.Errorf("Stderr:\n\tExpected %q\n\tGot %q", test.err, err)
			}
		})
	}
}

type frameFunc struct {
	wdte.Frame
}

func (f frameFunc) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return f
}

func (f frameFunc) Compare(other wdte.Func) (int, bool) {
	o := other.(frameFunc)

	if f.ID() != o.ID() {
		return 1, false
	}
	if !reflect.DeepEqual(f.Args(), o.Args()) {
		return 1, false
	}

	fp := frameFunc{f.Parent()}
	op := frameFunc{o.Parent()}

	if ((fp.ID() == "") || (op.ID() == "")) && (fp.ID() != op.ID()) {
		return 1, false
	}

	return fp.Compare(op)
}

func TestBasics(t *testing.T) {
	//imFrame := wdte.ImportFunc(func(from string) (*wdte.Module, error) {
	//	return &wdte.Module{
	//		Funcs: map[wdte.ID]wdte.Func{
	//			"get": wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	//				return frameFunc{frame.WithID("get")}
	//			}),
	//		},
	//	}, nil
	//})

	//frame := wdte.CustomFrame("unknown function, maybe Go", []wdte.Func{}, nil)
	//frame = wdte.CustomFrame("main", []wdte.Func{}, &frame)
	//frame = wdte.CustomFrame("test", []wdte.Func{}, &frame).Pos(1, 45)
	//frame = wdte.CustomFrame("get", []wdte.Func{}, &frame).Pos(1, 26)

	runTests(t, []test{
		{
			name:   "Fib",
			script: `main n => switch n { <= 1 => n; default => + (main (- n 2)) (main (- n 1)); };`,
			args:   []wdte.Func{wdte.Number(12)},
			ret:    wdte.Number(144),
		},
		{
			// Wonder why memo exists? Try removing the keyword from this
			// test script and see what happens.
			name:   "Fib/Memo",
			script: `memo main n => switch n { <= 1 => n; default => + (main (- n 2)) (main (- n 1)); };`,
			args:   []wdte.Func{wdte.Number(38)},
			ret:    wdte.Number(39088169),
		},
		{
			name:   "PassModule",
			script: `'somemodule' => m; test im => im.num; main => test m;`,
			im: wdte.ImportFunc(func(from string) (*wdte.Module, error) {
				return &wdte.Module{
					Funcs: map[wdte.ID]wdte.Func{
						"num": wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
							return wdte.Number(3)
						}),
					},
				}, nil
			}),
			ret: wdte.Number(3),
		},
		//{
		//	name:   "Frame",
		//	script: `'frame' => frame; test => frame.get; main => test;`,
		//	im:     imFrame,
		//	ret:    frameFunc{frame},
		//},
	})
}

func TestMath(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Abs",
			script: `'math' => m; main n => m.abs n;`,
			args:   []wdte.Func{wdte.Number(-3)},
			ret:    wdte.Number(3),
		},
		{
			name:   "Ceil",
			script: `'math' => m; main n => m.ceil n;`,
			args:   []wdte.Func{wdte.Number(1.1)},
			ret:    wdte.Number(2),
		},
		{
			name:   "Floor",
			script: `'math' => m; main n => m.floor n;`,
			args:   []wdte.Func{wdte.Number(1.1)},
			ret:    wdte.Number(1),
		},
		{
			name:   "Sin",
			script: `'math' => m; main n => m.sin n;`,
			args:   []wdte.Func{wdte.Number(3)},
			ret:    wdte.Number(math.Sin(3)),
		},
		{
			name:   "Cos",
			script: `'math' => m; main n => m.cos n;`,
			args:   []wdte.Func{wdte.Number(3)},
			ret:    wdte.Number(math.Cos(3)),
		},
		{
			name:   "Tan",
			script: `'math' => m; main n => m.tan n;`,
			args:   []wdte.Func{wdte.Number(3)},
			ret:    wdte.Number(math.Tan(3)),
		},
		{
			name:   "π",
			script: `'math' => m; main => m.pi;`,
			ret:    wdte.Number(math.Pi),
		},
	})
}

func TestStream(t *testing.T) {
	runTests(t, []test{
		{
			name:   "New/Args",
			script: `'stream' => s; main a b c => s.new a b c -> s.collect;`,
			args:   []wdte.Func{wdte.Number(3), wdte.Number(6), wdte.Number(9)},
			ret:    wdte.Array{wdte.Number(3), wdte.Number(6), wdte.Number(9)},
		},
		{
			name:   "Range",
			script: `'stream' => s; main start end step => s.range start end step -> s.collect;`,
			args:   []wdte.Func{wdte.Number(3), wdte.Number(12), wdte.Number(3)},
			ret:    wdte.Array{wdte.Number(3), wdte.Number(6), wdte.Number(9)},
		},
		{
			name:   "Concat",
			script: `'stream' => s; main => s.concat (s.range 2) (s.range 3) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(1), wdte.Number(0), wdte.Number(1), wdte.Number(2)},
		},
		{
			name:   "Map",
			script: `'stream' => s; main => s.range 3 -> s.map (* 5) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(5), wdte.Number(10)},
		},
		{
			name:   "Filter",
			script: `'stream' => s; main => s.range 5 -> s.filter (< 3) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(1), wdte.Number(2)},
		},
		{
			name:   "Reduce",
			script: `'stream' => s; main => s.range 1 6 -> s.reduce 1 *;`,
			ret:    wdte.Number(120),
		},
		{
			name:   "Any/True",
			script: `'stream' => s; main => s.range 5 -> s.any (== 3);`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "Any/False",
			script: `'stream' => s; main => s.range 3 -> s.any (== 3);`,
			ret:    wdte.Bool(false),
		},
	})
}

func TestIO(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Write",
			script: `'io' => io; main => 'test' -> io.write io.stdout;`,
			out:    "test",
		},
		{
			name:   "Writeln",
			script: `'io' => io; main => 'test' -> io.writeln io.stdout;`,
			out:    "test\n",
		},
		{
			name:   "Lines",
			script: `'io' => io; 'stream' => s; main str => io.readString str -> io.lines -> s.collect;`,
			args:   []wdte.Func{wdte.String("Line 1\nLine 2\nLine 3")},
			ret:    wdte.Array{wdte.String("Line 1"), wdte.String("Line 2"), wdte.String("Line 3")},
		},
		{
			name:   "Scan",
			script: `'io' => io; 'stream' => s; main str => io.readString str -> io.scan '|||' -> s.collect;`,
			args:   []wdte.Func{wdte.String("Part 1|||Part 2|||Part 3")},
			ret:    wdte.Array{wdte.String("Part 1"), wdte.String("Part 2"), wdte.String("Part 3")},
		},
	})
}
