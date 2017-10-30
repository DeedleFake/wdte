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
	disabled bool

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
			if test.disabled {
				t.SkipNow()
			}

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

func TestBasics(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Simple",
			script: `main => 3;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Simple/Memo",
			script: `memo test n => + n 3; main => (test 5; test 5);`,
			ret:    wdte.Number(8),
		},
		{
			name:   "Chain",
			script: `main => 1 -> + 2 -> - 3;`,
			ret:    wdte.Number(0),
		},
		{
			name:   "Chain/Slot",
			script: `main => 1 : a -> + 2 : b -> - (* a 3) -> + b;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Chain/Ignored",
			script: `main => 1 -> + 2 -- + 5 -> - 1;`,
			ret:    wdte.Number(2),
		},
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
		{
			name:   "Array/Args",
			script: `'io' => io; 'arrays' => a; 'stream' => s; test a => [a]; main => a.stream (test 3) -> s.map (io.writeln io.stdout) -> s.drain;`,
			out:    "3\n",
		},
		{
			name:   "Lambda",
			script: `test a => a 3; main => test (@ t n => * n 2);`,
			ret:    wdte.Number(6),
		},
		{
			name:   "Lambda/Closure",
			script: `test a => a 3; main q => test (@ t n => * n q);`,
			args:   []wdte.Func{wdte.Number(2)},
			ret:    wdte.Number(6),
		},
		{
			name:   "Lambda/Fib",
			script: `test a => a 10; main => test (@ t n => switch n { <= 1 => n; default => + (t (- n 2)) (t (- n 1)); };);`,
			ret:    wdte.Number(55),
		},
		{
			name:   "Lambda/Fib/Memo",
			script: `test a => a 38; main => test (@ memo t n => switch n { <= 1 => n; default => + (t (- n 2)) (t (- n 1)); };);`,
			ret:    wdte.Number(39088169),
		},
		{
			name:   "True",
			script: `main => true;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "False",
			script: `main => false;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "And/True",
			script: `main => && true true;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "And/False",
			script: `main => && true false;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "Or/True",
			script: `main => || false true;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "Or/False",
			script: `main => || false false;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "Not/True",
			script: `main => ! false;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "Not/False",
			script: `main => ! true;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "ReturnFunc",
			script: `test => +; main => test 2 3;`,
			ret:    wdte.Number(5),
		},
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
			name:   "FlatMap",
			script: `'stream' => s; test a => [a; + a 1]; main => s.range 3 -> s.flatMap test -> s.collect;`,
			ret: wdte.Array{
				wdte.Number(0),
				wdte.Number(1),
				wdte.Number(1),
				wdte.Number(2),
				wdte.Number(2),
				wdte.Number(3),
			},
		},
		{
			name:   "Enumerate",
			script: `'stream' => s; main => s.new 'a' 'b' 'c' -> s.enumerate -> s.collect;`,
			ret: wdte.Array{
				wdte.Array{wdte.Number(0), wdte.String("a")},
				wdte.Array{wdte.Number(1), wdte.String("b")},
				wdte.Array{wdte.Number(2), wdte.String("c")},
			},
		},
		{
			name:   "Drain",
			script: `'stream' => s; 'io' => io; main => s.range 5 -> s.map (io.writeln io.stdout) -> s.drain;`,
			out:    "0\n1\n2\n3\n4\n",
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
		{
			name:   "All/True",
			script: `'stream' => s; main => s.range 5 -> s.all (< 5);`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "All/False",
			script: `'stream' => s; main => s.range 5 -> s.all (< 3);`,
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

func TestStrings(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Contains",
			script: `'stream' => s; 'strings' => str; main => s.new "this" "is" "a" "test" -> s.filter (str.contains "t") -> s.collect;`,
			ret:    wdte.Array{wdte.String("this"), wdte.String("test")},
		},
		{
			name:   "Prefix",
			script: `'stream' => s; 'strings' => str; main => s.new "this" "is" "a" "test" -> s.filter (str.prefix "i") -> s.collect;`,
			ret:    wdte.Array{wdte.String("is")},
		},
		{
			name:   "Suffix",
			script: `'stream' => s; 'strings' => str; main => s.new "this" "is" "a" "test" -> s.filter (str.suffix "t") -> s.collect;`,
			ret:    wdte.Array{wdte.String("test")},
		},
		{
			name:   "Index",
			script: `'stream' => s; 'strings' => str; main => s.new 'abcde' 'bcdef' 'cdefg' 'defgh' 'efghi' -> s.map (str.index 'cd') -> s.collect;`,
			ret:    wdte.Array{wdte.Number(2), wdte.Number(1), wdte.Number(0), wdte.Number(-1), wdte.Number(-1)},
		},
		{
			name:   "Len",
			script: `'strings' => str; main => str.len 'abc';`,
			ret:    wdte.Number(3),
		},
		{
			name:   "At",
			script: `'strings' => str; main => str.at 'test' 2;`,
			ret:    wdte.String('s'),
		},
		{
			name:   "Upper",
			script: `'strings' => str; main => str.upper 'QwErTy';`,
			ret:    wdte.String("QWERTY"),
		},
		{
			name:   "Lower",
			script: `'strings' => str; main => str.lower 'QwErTy';`,
			ret:    wdte.String("qwerty"),
		},
		{
			name:   "Format",
			script: `'strings' => str; main => str.format '{#2}{#0}{}' 3 6 9;`,
			ret:    wdte.String("936"),
		},
		{
			name:   "Format/Type",
			script: `'strings' => str; main => str.format '{?}' 3;`,
			ret:    wdte.String("wdte.Number(3)"),
		},
		{
			name:   "Format/Quote",
			script: `'strings' => str; main => str.format '{q}' 'It is as if the socialists were to accuse us of not wanting persons to eat because we do not want the state to raise grain.';`,
			ret:    wdte.String(`"It is as if the socialists were to accuse us of not wanting persons to eat because we do not want the state to raise grain."`),
		},
	})
}

func TestArrays(t *testing.T) {
	runTests(t, []test{
		{
			name:   "At",
			script: `'arrays' => a; main => [3; 6; 9] -> a.at 1;`,
			ret:    wdte.Number(6),
		},
		{
			name:   "Stream",
			script: `'arrays' => a; 'stream' => s; main => a.stream ['this'; 'is'; 'a'; 'test'] -> s.collect;`,
			ret:    wdte.Array{wdte.String("this"), wdte.String("is"), wdte.String("a"), wdte.String("test")},
		},
	})
}

//func ExampleModule_Eval() {
//	const src = `
//	'math' => m;
//
//	npi a => * m.pi a;
//`
//
//	m, err := std.Module().Parse(strings.NewReader(src), std.Import)
//	if err != nil {
//		log.Fatalf("Failed to parse module: %v", err)
//	}
//
//	r, err := m.Eval(strings.NewReader("npi 5"))
//	if err != nil {
//		log.Fatalf("Failed to evaluate: %v", err)
//	}
//
//	fmt.Println(r.Call(wdte.F()))
//	// Output: 15.707963267948966
//}
