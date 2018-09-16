package wdte_test

import (
	"bytes"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/arrays"
	"github.com/DeedleFake/wdte/std/io"
	_ "github.com/DeedleFake/wdte/std/math"
	_ "github.com/DeedleFake/wdte/std/stream"
	_ "github.com/DeedleFake/wdte/std/strings"
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
				im = wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
					scope, err := std.Import(from)
					if err != nil {
						return nil, err
					}

					switch from {
					case "io":
						return scope.Map(map[wdte.ID]wdte.Func{
							"stdin":  io.Reader{strings.NewReader(test.in)},
							"stdout": io.Writer{&stdout},
							"stderr": io.Writer{&stderr},
						}), nil
					}

					return scope, nil
				})
			}

			m, err := wdte.Parse(strings.NewReader(test.script), im)
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

		known := scope.Known()
		if len(known) != 3 {
			t.Errorf("Expected to find 3 variables in scope. Found %v", len(known))
		}

		var found int
		for _, id := range known {
			switch id {
			case "x", "test", "q":
				found++
			}
		}
		if found != 3 {
			t.Errorf("Expected to find %q, %q, and %q in scope.\nFound %v",
				"x",
				"test",
				"q",
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
			name:   "Simple/Args",
			script: `let test n => + n 3; test 5;`,
			ret:    wdte.Number(8),
		},
		{
			name:   "Simple/Infix",
			script: `{3 + 2};`,
			ret:    wdte.Number(5),
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
			name:   "Chain/Ignored",
			script: `1 -> + 2 -- + 5 -> - 1;`,
			ret:    wdte.Number(2),
		},
		{
			name:   "Fib",
			script: `let main n => switch n { <= 1 => n; true => + (main (- n 2)) (main (- n 1)); }; main 12;`,
			ret:    wdte.Number(144),
		},
		{
			// Wonder why memo exists? Try removing the keyword from this
			// test script and see what happens.
			name:   "Fib/Memo",
			script: `let memo main n => switch n { <= 1 => n; true => + (main (- n 2)) (main (- n 1)); }; main 38;`,
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
			script: `let test a => a 10; test (@ t n => switch n { <= 1 => n; true => + (t (- n 2)) (t (- n 1)); };);`,
			ret:    wdte.Number(55),
		},
		{
			name:   "Lambda/Fib/Memo",
			script: `let test a => a 38; test (@ memo t n => switch n { <= 1 => n; true => + (t (- n 2)) (t (- n 1)); };);`,
			ret:    wdte.Number(39088169),
		},
		{
			name:   "Lambda/Compound",
			script: `let io => import 'io'; let test a => a 3; test (@ t n => io.write io.stdout 'Test'; + n 2);`,
			ret:    wdte.Number(5),
			out:    "Test",
		},
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
			name:   "Collect",
			script: `let t => collect (let test => 3); t.test;`,
			ret:    wdte.Number(3),
		},
		{
			name:   "Known",
			script: `let t => collect (let test => 3; let other => 5); known t;`,
			ret:    wdte.Array{wdte.String("other"), wdte.String("test")},
		},
		{
			// TODO: Move this and some others into a separate test, such as
			// `TestStd` or something.
			name:   "Sub",
			script: `let t => collect (let test => 3); let t => sub t 'test2' 5; t.test2;`,
			ret:    wdte.Number(5),
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
			script: `let s => import 'stream'; let main a b c => s.new a b c -> s.collect;`,
			args:   []wdte.Func{wdte.Number(3), wdte.Number(6), wdte.Number(9)},
			ret:    wdte.Array{wdte.Number(3), wdte.Number(6), wdte.Number(9)},
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
			script: `let s => import 'stream'; let test a => s.new a (+ a 1); let main => s.range 3 -> s.flatMap test -> s.collect;`,
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
			script: `let s => import 'stream'; let main => s.new 'a' 'b' 'c' -> s.enumerate -> s.collect;`,
			ret: wdte.Array{
				wdte.Array{wdte.Number(0), wdte.String("a")},
				wdte.Array{wdte.Number(1), wdte.String("b")},
				wdte.Array{wdte.Number(2), wdte.String("c")},
			},
		},
		{
			name:   "Drain",
			script: `let s => import 'stream'; let io => import 'io'; let main => s.range 5 -> s.map (io.writeln io.stdout) -> s.drain;`,
			out:    "0\n1\n2\n3\n4\n",
		},
		{
			name:   "Reduce",
			script: `let s => import 'stream'; let main => s.range 1 6 -> s.reduce 1 *;`,
			ret:    wdte.Number(120),
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
			name:   "Lines",
			script: `let io => import 'io'; let s => import 'stream'; let main str => io.readString str -> io.lines -> s.collect;`,
			args:   []wdte.Func{wdte.String("Line 1\nLine 2\nLine 3")},
			ret:    wdte.Array{wdte.String("Line 1"), wdte.String("Line 2"), wdte.String("Line 3")},
		},
		{
			name:   "Scan",
			script: `let io => import 'io'; let s => import 'stream'; let main str => io.readString str -> io.scan '|||' -> s.collect;`,
			args:   []wdte.Func{wdte.String("Part 1|||Part 2|||Part 3")},
			ret:    wdte.Array{wdte.String("Part 1"), wdte.String("Part 2"), wdte.String("Part 3")},
		},
	})
}

func TestStrings(t *testing.T) {
	runTests(t, []test{
		{
			name:   "Contains",
			script: `let s => import 'stream'; let str => import 'strings'; let main => s.new "this" "is" "a" "test" -> s.filter (str.contains "t") -> s.collect;`,
			ret:    wdte.Array{wdte.String("this"), wdte.String("test")},
		},
		{
			name:   "Prefix",
			script: `let s => import 'stream'; let str => import 'strings'; let main => s.new "this" "is" "a" "test" -> s.filter (str.prefix "i") -> s.collect;`,
			ret:    wdte.Array{wdte.String("is")},
		},
		{
			name:   "Suffix",
			script: `let s => import 'stream'; let str => import 'strings'; let main => s.new "this" "is" "a" "test" -> s.filter (str.suffix "t") -> s.collect;`,
			ret:    wdte.Array{wdte.String("test")},
		},
		{
			name:   "Index",
			script: `let s => import 'stream'; let str => import 'strings'; let main => s.new 'abcde' 'bcdef' 'cdefg' 'defgh' 'efghi' -> s.map (str.index 'cd') -> s.collect;`,
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
			name: "SortStable",
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
