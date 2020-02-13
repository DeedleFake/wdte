package wdte_test

import (
	"bytes"
	"io"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/scanner"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/arrays"
	wdteio "github.com/DeedleFake/wdte/std/io"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/rand"
	"github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/strings"
)

type test struct {
	disabled bool

	name string

	script string
	im     wdte.Importer
	macros scanner.MacroMap

	args []wdte.Func
	ret  wdte.Func

	in  string
	out string
	err string
}

func runTests(t *testing.T, tests []test) {
	t.Helper()
	t.Parallel()

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Helper()
			t.Parallel()

			if test.disabled {
				t.SkipNow()
			}

			var stdout, stderr bytes.Buffer

			im := test.im
			if im == nil {
				im = wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
					scope, err := std.Import(from)
					if err != nil {
						return nil, err
					}

					switch from {
					case "io":
						return scope.Map(map[wdte.ID]wdte.Func{
							"stdin":  wdteio.Reader{Reader: strings.NewReader(test.in)},
							"stdout": wdteio.Writer{Writer: &stdout},
							"stderr": wdteio.Writer{Writer: &stderr},
						}), nil
					}

					return scope, nil
				})
			}

			m, err := wdte.Parse(strings.NewReader(test.script), im, test.macros)
			if err != nil {
				t.Fatalf("Failed to parse script: %v", err)
			}

			ret := m.Call(std.F(), test.args...)

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
	t.Run("Scope/Known", func(t *testing.T) {
		scope := wdte.S()
		scope = scope.Add("x", wdte.Number(3))
		scope = scope.Add("test", wdte.String("This is a test."))
		scope = scope.Map(map[wdte.ID]wdte.Func{
			"q":    wdte.String("Other"),
			"test": wdte.String("Or is it?"),
		})
		scope = scope.Custom(func(id wdte.ID) wdte.Func {
			if id == "a" {
				return wdte.Bool(true)
			}

			return nil
		}, func(known map[wdte.ID]struct{}) {
			known["a"] = struct{}{}
		})

		known := scope.Known()
		if len(known) != 4 {
			t.Errorf("Expected to find 4 variables in scope. Found %v", len(known))
		}

		var found int
		for _, id := range known {
			switch id {
			case "x", "test", "q", "a":
				found++
			}
		}
		if found != 4 {
			t.Errorf("Expected to find %q, %q, %q, and %q in scope.\nFound %v",
				"x",
				"test",
				"q",
				"a",
				known,
			)
		}
	})

	runTests(t, []test{
		{
			name:   "Simple",
			script: `3;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Simple/Func",
			script: `let main => 3; main;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Simple/Pattern",
			script: `let [a b] => [3; 5]; [b; a];`,
			ret:    wdte.Array{wdte.Number(5), wdte.Number(3)},
		},
		{
			name:   "Simple/Args",
			script: `let test n => + n 3; test 5;`,
			ret:    wdte.Number(8),
		},
		{
			name:   "Simple/Compound/Args",
			script: `let test x y => + x y; (test 3) 5;`,
			ret:    wdte.Number(8),
		},
		{
			name:   "Simple/Memo",
			script: `let memo test n => + n 3; let main => (test 5; test 5); main;`,
			ret:    wdte.Number(8),
		},
		{
			name:   "Simple/VariableArgs",
			script: `let test => +; test 3 5;`,
			ret:    wdte.Number(8),
		},
		{
			name:   "Simple/Number/Compare",
			script: `[< .2 .5; < .6 .5; > .2 .5; > .6 .5];`,
			ret:    wdte.Array{wdte.Bool(true), wdte.Bool(false), wdte.Bool(false), wdte.Bool(true)},
		},
		{
			name:   "Simple/String/Compare",
			script: `[< 'a' 'b'; > 'a' 'b'];`,
			ret:    wdte.Array{wdte.Bool(true), wdte.Bool(false)},
		},
		{
			name:   "Chain",
			script: `1 -> + 2 -> - 3;`,
			ret:    wdte.Number(0),
		},
		{
			name:   "Chain/Slot",
			script: `1 : a -> + 2 : b -> - (* a 3) -> + b;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Chain/Slot/Pattern",
			script: `[1; 2] : a -> a : [a b] -> [b; a];`,
			ret:    wdte.Array{wdte.Number(2), wdte.Number(1)},
		},
		{
			name:   "Chain/Ignored",
			script: `1 -> + 2 -- + 5 -> - 1;`,
			ret:    wdte.Number(2),
		},
		{
			name:   "Chain/Error",
			script: `1 -| 'Not broken yet.' -> a -> + 3 -| 'It broke.' -| 'Or did it?';`,
			ret:    wdte.String("It broke."),
		},
		{
			name:   "Chain/Error/Ignored",
			script: `1 -- a -> + 3 -| 'It broke.';`,
			ret:    wdte.String("It broke."),
		},
		{
			name:   "Fib",
			script: `let main n => n { <= 1 => n; true => + (main (- n 2)) (main (- n 1)); }; main 12;`,
			ret:    wdte.Number(144),
		},
		{
			// Wonder why memo exists? Try removing the keyword from this
			// test script and see what happens.
			name:   "Fib/Memo",
			script: `let memo main n => n { <= 1 => n; true => + (main (- n 2)) (main (- n 1)); }; main 38;`,
			ret:    wdte.Number(39088169),
		},
		{
			name:   "PassModule",
			script: `let m => import 'somemodule'; let test im => im.num; test m;`,
			im: wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
				return wdte.S().Add("num", wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
					return wdte.Number(3)
				}),
				), nil
			}),
			ret: wdte.Number(3),
		},
		{
			name:   "Sub/Multiple",
			script: `let m => import 'somemodule'; (m).sub.num;`,
			im: wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
				return wdte.S().Add("sub", wdte.S().Add("num", wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
					return wdte.Number(3)
				}),
				)), nil
			}),
			ret: wdte.Number(3),
		},
		{
			name:   "Sub/Expr",
			script: `let m => import 'math'; m.(cos 3);`,
			ret:    wdte.Number(math.Cos(3)),
		},
		{
			name:   "Sub/Closure",
			script: `let m => import 'math'; let test val => m.(cos val); test 3;`,
			ret:    wdte.Number(math.Cos(3)),
		},
		{
			name:   "Array/Empty",
			script: `[];`,
			ret:    wdte.Array{},
		},
		{
			name:   "Array/Args",
			script: `let io => import 'io'; let a => import 'arrays'; let s => import 'stream'; let test a => [a]; a.stream (test 3) -> s.map (io.writeln io.stdout) -> s.drain;`,
			out:    "3\n",
		},
		{
			name:   "Lambda",
			script: `let test a => a 3; test (@ t n => * n 2);`,
			ret:    wdte.Number(6),
		},
		{
			name:   "Lambda/Closure",
			script: `let test a => a 3; let q => 2; test (@ t n => * n q);`,
			ret:    wdte.Number(6),
		},
		{
			name:   "Lambda/Fib",
			script: `let test a => a 10; test (@ t n => n { <= 1 => n; true => + (t (- n 2)) (t (- n 1)); };);`,
			ret:    wdte.Number(55),
		},
		{
			name:   "Lambda/Fib/Memo",
			script: `let test a => a 38; test (@ memo t n => n { <= 1 => n; true => + (t (- n 2)) (t (- n 1)); };);`,
			ret:    wdte.Number(39088169),
		},
		{
			name:   "Lambda/Compound",
			script: `let io => import 'io'; let test a => a 3; test (@ t n => io.write io.stdout 'Test'; + n 2);`,
			ret:    wdte.Number(5),
			out:    "Test",
		},
		{
			name:   "Let/Inner",
			script: `let x => 3; (let x => 5); x;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Let/Return",
			script: `let x => 3;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Let/NoArgs",
			script: `let x => 3; let x => + x 5; x;`,
			ret:    wdte.Number(8),
		},
		{
			name:   "ScopeCaching",
			script: `let io => import 'io'; let test => import 'test'; let x => test.reader -> io.string; [x; x];`,
			ret:    wdte.Array{wdte.String("one"), wdte.String("one")},
			im: wdte.ImportFunc(func(im string) (*wdte.Scope, error) {
				if im == "test" {
					return wdte.S().Map(map[wdte.ID]wdte.Func{
						"reader": wdteio.Reader{Reader: &pieceReader{
							pieces: []io.Reader{
								strings.NewReader("one"),
								strings.NewReader("two"),
							},
						}},
					}), nil
				}

				return std.Import(im)
			}),
		},
		{
			name:   "Macro/ROT13",
			script: `@rot13[test];`,
			ret:    wdte.String("grfg"),
			macros: scanner.MacroMap{
				"rot13": func(input string) ([]scanner.Token, error) {
					r := make([]rune, 0, len(input))
					for _, c := range input {
						switch {
						case (c >= 'a') && (c <= 'z'):
							r = append(r, (c-'a'+13)%26+'a')
						case (c >= 'A') && (c <= 'Z'):
							r = append(r, (c-'A'+13)%26+'A')
						default:
							r = append(r, c)
						}
					}
					return []scanner.Token{
						{Type: scanner.String, Val: string(r)},
					}, nil
				},
			},
		},
	})
}

func TestStd(t *testing.T) {
	runTests(t, []test{
		{
			name:   "True",
			script: `true;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "False",
			script: `false;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "And/True",
			script: `&& true true;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "And/False",
			script: `&& true false;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "Or/True",
			script: `|| false true;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "Or/False",
			script: `|| false false;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "Not/True",
			script: `! false;`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "Not/False",
			script: `! true;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "Len/String",
			script: `len 'test';`,
			ret:    wdte.Number(4),
		},
		{
			name:   "Len/Array",
			script: `len [3; 5; 1];`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Len/Other",
			script: `len 5;`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "At/String",
			script: `at 'test' 2;`,
			ret:    wdte.String("s"),
		},
		{
			name:   "At/Array",
			script: `at [3; 5; 1] 0;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "At/Scope",
			script: `let m => import 'math'; at m 'pi';`,
			ret:    wdte.Number(math.Pi),
		},
		{
			disabled: true,
			name:     "Set/Scope",
			script:   `let t => collect (let test => 3); let t => set t 'test2' 5; t.test2;`,
			ret:      wdte.Number(5),
		},
		{
			name:   "Set/Array",
			script: `let t => [1; 2; 3]; set t 1 5;`,
			ret:    wdte.Array{wdte.Number(1), wdte.Number(5), wdte.Number(3)},
		},
		{
			disabled: true,
			name:     "Collect",
			script:   `let t => collect (let test => 3); t.test;`,
			ret:      wdte.Number(3),
		},
		{
			disabled: true,
			name:     "Known",
			script:   `let t => collect (let test => 3; let other => 5); known t;`,
			ret:      wdte.Array{wdte.String("other"), wdte.String("test")},
		},
		{
			name:   "Reflect",
			script: `[reflect 'string' 'String'; 'string' {reflect 'String' => 'test'}];`,
			ret:    wdte.Array{wdte.Bool(true), wdte.String("test")},
		},
	})
}

func TestMath(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Abs",
			script: `let m => import 'math'; m.abs -3;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Ceil",
			script: `let m => import 'math'; m.ceil 1.1;`,
			ret:    wdte.Number(2),
		},
		{
			name:   "Floor",
			script: `let m => import 'math'; m.floor 1.1;`,
			ret:    wdte.Number(1),
		},
		{
			name:   "Sin",
			script: `let m => import 'math'; m.sin 3;`,
			ret:    wdte.Number(math.Sin(3)),
		},
		{
			name:   "Cos",
			script: `let m => import 'math'; m.cos 3;`,
			ret:    wdte.Number(math.Cos(3)),
		},
		{
			name:   "Tan",
			script: `let m => import 'math'; m.tan 3;`,
			ret:    wdte.Number(math.Tan(3)),
		},
		{
			name:   "π",
			script: `let m => import 'math'; m.pi;`,
			ret:    wdte.Number(math.Pi),
		},
	})
}

func TestStream(t *testing.T) {
	runTests(t, []test{
		{
			name:   "New",
			script: `let s => import 'stream'; s.new 0 (@ f n => + n 1 {> 5 => s.end}) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(1), wdte.Number(2), wdte.Number(3), wdte.Number(4), wdte.Number(5)},
		},
		{
			name:   "Range/1",
			script: `let s => import 'stream'; let main start end step => s.range start end step -> s.collect;`,
			args:   []wdte.Func{wdte.Number(3), wdte.Number(12), wdte.Number(3)},
			ret:    wdte.Array{wdte.Number(3), wdte.Number(6), wdte.Number(9)},
		},
		{
			name:   "Range/2",
			script: `let s => import 'stream'; s.range 1 3 -> s.collect;`,
			ret:    wdte.Array{wdte.Number(1), wdte.Number(2)},
		},
		{
			name:   "Range/3",
			script: `let s => import 'stream'; s.range 1 6 2 -> s.collect;`,
			ret:    wdte.Array{wdte.Number(1), wdte.Number(3), wdte.Number(5)},
		},
		{
			name:   "Concat",
			script: `let s => import 'stream'; let main => s.concat (s.range 2) (s.range 3) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(1), wdte.Number(0), wdte.Number(1), wdte.Number(2)},
		},
		{
			name:   "Map",
			script: `let s => import 'stream'; let main => s.range 3 -> s.map (* 5) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(5), wdte.Number(10)},
		},
		{
			name:   "Filter",
			script: `let s => import 'stream'; let main => s.range 5 -> s.filter (< 3) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(1), wdte.Number(2)},
		},
		{
			name:   "FlatMap",
			script: `let a => import 'arrays'; let s => import 'stream'; let test n => a.stream [n; + n 1]; let main => s.range 3 -> s.flatMap test -> s.collect;`,
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
			script: `let a => import 'arrays'; let s => import 'stream'; let main => a.stream ['a'; 'b'; 'c'] -> s.enumerate -> s.collect;`,
			ret: wdte.Array{
				wdte.Array{wdte.Number(0), wdte.String("a")},
				wdte.Array{wdte.Number(1), wdte.String("b")},
				wdte.Array{wdte.Number(2), wdte.String("c")},
			},
		},
		{
			name:   "RepeatAndLimit",
			script: `let s => import 'stream'; s.range 3 -> s.repeat -> s.limit 9 -> s.collect;`,
			ret: wdte.Array{
				wdte.Number(0), wdte.Number(1), wdte.Number(2),
				wdte.Number(0), wdte.Number(1), wdte.Number(2),
				wdte.Number(0), wdte.Number(1), wdte.Number(2),
			},
		},
		{
			name:   "Skip",
			script: `let s => import 'stream'; s.range 3 -> s.skip 2 -> s.collect;`,
			ret:    wdte.Array{wdte.Number(2)},
		},
		{
			name:   "Zip",
			script: `let s => import 'stream'; s.zip (s.range 2) (s.range 1 2) -> s.collect;`,
			ret: wdte.Array{
				wdte.Array{wdte.Number(0), wdte.Number(1)},
				wdte.Array{wdte.Number(1), stream.End()},
			},
		},
		{
			name: "Drain",
			script: `
				let s => import 'stream';
				let io => import 'io';
				let main =>
					s.range 5
					-> s.map (@ p v => io.writeln io.stdout v; v)
					-> s.drain
					;
				`,
			ret: wdte.Number(4),
			out: "0\n1\n2\n3\n4\n",
		},
		{
			name:   "Reduce",
			script: `let s => import 'stream'; let main => s.range 1 6 -> s.reduce 1 *;`,
			ret:    wdte.Number(120),
		},
		{
			name:   "Fold",
			script: `let s => import 'stream'; s.range 5 -> s.fold +;`,
			ret:    wdte.Number(10),
		},
		{
			name:   "Extent/Min",
			script: `let s => import 'stream'; let a => import 'arrays'; a.stream [5; 3; 6; 2; 7; 1] -> s.extent 3 <;`,
			ret:    wdte.Array{wdte.Number(1), wdte.Number(2), wdte.Number(3)},
		},
		{
			name:   "Extent/Max",
			script: `let s => import 'stream'; let a => import 'arrays'; a.stream [5; 3; 6; 2; 7; 1] -> s.extent 3 >;`,
			ret:    wdte.Array{wdte.Number(7), wdte.Number(6), wdte.Number(5)},
		},
		{
			name:   "Extent/Range",
			script: `let s => import 'stream'; s.range 10 -> s.extent 3 >;`,
			ret:    wdte.Array{wdte.Number(9), wdte.Number(8), wdte.Number(7)},
		},
		{
			name:   "Extent/Sort",
			script: `let s => import 'stream'; s.concat (s.range 5) (s.range 2 10 3) -> s.extent -1 <;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(1), wdte.Number(2), wdte.Number(2), wdte.Number(3), wdte.Number(4), wdte.Number(5), wdte.Number(8)},
		},
		{
			name:   "Any/True",
			script: `let s => import 'stream'; let main => s.range 5 -> s.any (== 3);`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "Any/False",
			script: `let s => import 'stream'; let main => s.range 3 -> s.any (== 3);`,
			ret:    wdte.Bool(false),
		},
		{
			name:   "All/True",
			script: `let s => import 'stream'; let main => s.range 5 -> s.all (< 5);`,
			ret:    wdte.Bool(true),
		},
		{
			name:   "All/False",
			script: `let s => import 'stream'; let main => s.range 5 -> s.all (< 3);`,
			ret:    wdte.Bool(false),
		},
		{
			name: "FibLike",
			script: `
				let s => import 'stream';
				let a => import 'arrays';

				let pop array => a.stream array -> s.skip 1 -> s.collect;

				let windows size stream => s.new
						(stream -> s.limit size -> s.collect)
						(@ next prev =>
							let n => stream -> s.limit 1 -> s.collect;
							len n {
								== 0 => s.end;
								true => a.concat (pop prev) n;
							};
						)
					;

				let isFibLike stream => stream
					-> windows 3
					-> s.all (@ all v => == (at v 2) (+ (at v 1) (at v 0)))
					;

				[
					[1; 1; 2; 3; 5; 8] -> a.stream -> isFibLike;
					[2; 2; 4; 6; 10] -> a.stream -> isFibLike;
					[1; 1; 2; 3; 5; 7] -> a.stream -> isFibLike;
				];
			`,
			ret: wdte.Array{wdte.Bool(true), wdte.Bool(true), wdte.Bool(false)},
		},
	})
}

func TestIO(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Write",
			script: `let io => import 'io'; let main => 'test' -> io.write io.stdout;`,
			out:    "test",
		},
		{
			name:   "Writeln",
			script: `let io => import 'io'; let main => 'test' -> io.writeln io.stdout;`,
			out:    "test\n",
		},
		{
			name:   "Panic",
			script: `let io => import 'io'; + a b -| io.panic io.stderr 'Failed to add a and b' -| 3;`,
			err:    `Failed to add a and b: "a" is not in scope` + "\n",
		},
		{
			name:   "Lines",
			script: `let io => import 'io'; let s => import 'stream'; let str => import 'strings'; let main v => str.read v -> io.lines -> s.collect;`,
			args:   []wdte.Func{wdte.String("Line 1\nLine 2\nLine 3")},
			ret:    wdte.Array{wdte.String("Line 1"), wdte.String("Line 2"), wdte.String("Line 3")},
		},
		{
			name:   "Scan",
			script: `let io => import 'io'; let s => import 'stream'; let str => import 'strings'; let main v => str.read v -> io.scan '|||' -> s.collect;`,
			args:   []wdte.Func{wdte.String("Part 1|||Part 2|||Part 3")},
			ret:    wdte.Array{wdte.String("Part 1"), wdte.String("Part 2"), wdte.String("Part 3")},
		},
	})
}

func TestStrings(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Contains",
			script: `let a => import 'arrays'; let s => import 'stream'; let str => import 'strings'; let main => a.stream ["this"; "is"; "a"; "test"] -> s.filter (str.contains "t") -> s.collect;`,
			ret:    wdte.Array{wdte.String("this"), wdte.String("test")},
		},
		{
			name:   "Prefix",
			script: `let a => import 'arrays'; let s => import 'stream'; let str => import 'strings'; let main => a.stream ["this"; "is"; "a"; "test"] -> s.filter (str.prefix "i") -> s.collect;`,
			ret:    wdte.Array{wdte.String("is")},
		},
		{
			name:   "Suffix",
			script: `let a => import 'arrays'; let s => import 'stream'; let str => import 'strings'; let main => a.stream ["this"; "is"; "a"; "test"] -> s.filter (str.suffix "t") -> s.collect;`,
			ret:    wdte.Array{wdte.String("test")},
		},
		{
			name:   "Index",
			script: `let a => import 'arrays'; let s => import 'stream'; let str => import 'strings'; let main => a.stream ['abcde'; 'bcdef'; 'cdefg'; 'defgh'; 'efghi'] -> s.map (str.index 'cd') -> s.collect;`,
			ret:    wdte.Array{wdte.Number(2), wdte.Number(1), wdte.Number(0), wdte.Number(-1), wdte.Number(-1)},
		},
		{
			name:   "Upper",
			script: `let str => import 'strings'; let main => str.upper 'QwErTy';`,
			ret:    wdte.String("QWERTY"),
		},
		{
			name:   "Lower",
			script: `let str => import 'strings'; let main => str.lower 'QwErTy';`,
			ret:    wdte.String("qwerty"),
		},
		{
			name:   "Repeat",
			script: `let str => import 'strings'; str.repeat 'test' 3;`,
			ret:    wdte.String("testtesttest"),
		},
		{
			name:   "Split",
			script: `let str => import 'strings'; [str.split 'a test' ' '; str.split 'this is a test' ' ' 2; (str.split ' ') 'this is also a test' 3; (str.split ' ' 2) 'or is it'];`,
			ret: wdte.Array{
				wdte.Array{wdte.String("a"), wdte.String("test")},
				wdte.Array{wdte.String("this"), wdte.String("is a test")},
				wdte.Array{wdte.String("this"), wdte.String("is"), wdte.String("also a test")},
				wdte.Array{wdte.String("or"), wdte.String("is it")},
			},
		},
		{
			name:   "Join",
			script: `let str => import 'strings'; str.join ['this'; 'is'; 'a'; 'test'] ' ';`,
			ret:    wdte.String("this is a test"),
		},
		{
			name:   "Format",
			script: `let str => import 'strings'; let main => str.format '{#2}{#0}{}' 3 6 9;`,
			ret:    wdte.String("936"),
		},
		{
			name:   "Format/Type",
			script: `let str => import 'strings'; let main => str.format '{?}' 3;`,
			ret:    wdte.String("wdte.Number(3)"),
		},
		{
			name:   "Format/Quote",
			script: `let str => import 'strings'; let main => str.format '{q}' 'It is as if the socialists were to accuse us of not wanting persons to eat because we do not want the state to raise grain.';`,
			ret:    wdte.String(`"It is as if the socialists were to accuse us of not wanting persons to eat because we do not want the state to raise grain."`),
		},
		{
			name:   "Format/Array",
			script: `let str => import 'strings'; str.format '{}' [3; 5; 2];`,
			ret:    wdte.String(`[3; 5; 2]`),
		},
		{
			disabled: true,
			name:     "Format/Scope",
			script:   `let str => import 'strings'; let test => collect (let x => 3; let a => 2); str.format '{}' test;`,
			ret:      wdte.String(`scope(a: 2; x: 3)`),
		},
		{
			name:   "Format/Lambda",
			script: `let str => import 'strings'; str.format '{}' (@ s n => + n 2);`,
			ret:    wdte.String(`(@ s n => ...)`),
		},
		{
			name:   "Format/Partial",
			script: `let str => import 'strings'; let t => str.format '{} + {}: {}' 3 2; + 3 2 -> t;`,
			ret:    wdte.String("3 + 2: 5"),
		},
	})
}

func TestArrays(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Concat",
			script: `let a => import 'arrays'; a.concat [2; 5] [3; 6] [7];`,
			ret:    wdte.Array{wdte.Number(2), wdte.Number(5), wdte.Number(3), wdte.Number(6), wdte.Number(7)},
		},
		{
			name:   "Sort",
			script: `let a => import 'arrays'; a.sort [5; 3; 7] <;`,
			ret:    wdte.Array{wdte.Number(3), wdte.Number(5), wdte.Number(7)},
		},
		{
			disabled: true,
			name:     "SortStable",
			script: `
				let a => import 'arrays';
				let s => import 'stream';
				let str => import 'strings';

				a.sortStable [
						collect (
							let cat => 'two';
							let val => 5;
						);
						collect (
							let cat => 'one';
							let val => 7;
						);
						collect (
							let cat => 'two';
							let val => 3;
						);
					]
					(@ s e1 e2 => < e1.cat e2.cat)
				-> a.stream
				-> s.map (@ s e => str.format '{q} {}' e.cat e.val)
				-> s.collect
				;
			`,
			ret: wdte.Array{wdte.String(`"one" 7`), wdte.String(`"two" 5`), wdte.String(`"two" 3`)},
		},
		{
			name:   "Stream",
			script: `let a => import 'arrays'; let s => import 'stream'; let main => a.stream ['this'; 'is'; 'a'; 'test'] -> s.collect;`,
			ret:    wdte.Array{wdte.String("this"), wdte.String("is"), wdte.String("a"), wdte.String("test")},
		},
	})
}

func TestRand(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Simple",
			script: `let m => import 'math'; let rand => import 'rand'; rand.gen 1 -> rand.next -> * 100 -> m.floor;`,
			ret:    wdte.Number(math.Floor(rand.New(rand.NewSource(1)).Float64() * 100)),
		},
		{
			name:   "Stream",
			script: `let s => import 'stream'; let m => import 'math'; let rand => import 'rand'; rand.gen 1 -> rand.stream 3 -> s.map (* 100) -> s.map m.floor -> s.collect;`,
			ret:    wdte.Array{wdte.Number(60), wdte.Number(94), wdte.Number(66)},
		},
	})
}

type pieceReader struct {
	pieces []io.Reader
	i      int
}

func (r *pieceReader) Read(buf []byte) (int, error) {
	if r.i >= len(r.pieces) {
		return 0, io.EOF
	}

	n, err := r.pieces[r.i].Read(buf)
	if err == io.EOF {
		r.i++
	}
	return n, err
}
