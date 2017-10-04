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
						t.Errorf("Return:\n\tExpected %#v\n\tGot %#v", test.ret, ret)
					}

				default:
					if !reflect.DeepEqual(ret, test.ret) {
						t.Errorf("Return:\n\tExpected %#v\n\tGot %#v", test.ret, ret)
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
			name:   "Fib/Memo",
			script: `memo main n => switch n { <= 1 => n; default => + (main (- n 2)) (main (- n 1)); };`,
			args:   []wdte.Func{wdte.Number(38)},
			ret:    wdte.Number(39088169),
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
			name:   "Map",
			script: `'stream' => s; main => s.range 3 -> s.map (* 5) -> s.collect;`,
			ret:    wdte.Array{wdte.Number(0), wdte.Number(5), wdte.Number(10)},
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
	})
}

// Wonder why memo exists? I disabled it for this run.
//
//=== RUN   TestBasics
//=== RUN   TestBasics/Fib
//=== RUN   TestBasics/Fib/Memo
//SIGQUIT: quit
//PC=0x42dcc4 m=0 sigcode=0
//
//goroutine 0 [idle]:
//runtime.scang(0xc420001200, 0xc420020560)
//        /usr/lib/go/src/runtime/proc.go:893 +0xb4
//runtime.markroot.func1()
//        /usr/lib/go/src/runtime/mgcmark.go:264 +0x6d
//runtime.systemstack(0x7fff8cdd8020)
//        /usr/lib/go/src/runtime/asm_amd64.s:360 +0xab
//runtime.markroot(0xc420020560, 0xffffffff0000000b)
//        /usr/lib/go/src/runtime/mgcmark.go:245 +0x308
//runtime.gcDrain(0xc420020560, 0xd)
//        /usr/lib/go/src/runtime/mgcmark.go:912 +0xd4
//runtime.gcBgMarkWorker.func2()
//        /usr/lib/go/src/runtime/mgc.go:1796 +0x18f
//runtime.systemstack(0x626500)
//        /usr/lib/go/src/runtime/asm_amd64.s:344 +0x79
//runtime.mstart()
//        /usr/lib/go/src/runtime/proc.go:1125
//
//goroutine 9 [GC worker (idle)]:
//runtime.systemstack_switch()
//        /usr/lib/go/src/runtime/asm_amd64.s:298 fp=0xc42002a748 sp=0xc42002a740 pc=0x454b30
//runtime.gcBgMarkWorker(0xc42001f300)
//        /usr/lib/go/src/runtime/mgc.go:1760 +0x202 fp=0xc42002a7d8 sp=0xc42002a748 pc=0x4187d2
//runtime.goexit()
//        /usr/lib/go/src/runtime/asm_amd64.s:2337 +0x1 fp=0xc42002a7e0 sp=0xc42002a7d8 pc=0x4576b1
//created by runtime.gcBgMarkStartWorkers
//        /usr/lib/go/src/runtime/mgc.go:1655 +0x7e
//
//goroutine 1 [chan receive, 9 minutes]:
//testing.(*T).Run(0xc42008a000, 0x55feed, 0xa, 0x567ea0, 0x469126)
//        /usr/lib/go/src/testing/testing.go:790 +0x2fc
//testing.runTests.func1(0xc42008a000)
//        /usr/lib/go/src/testing/testing.go:1004 +0x64
//testing.tRunner(0xc42008a000, 0xc420083de0)
//        /usr/lib/go/src/testing/testing.go:746 +0xd0
//testing.runTests(0xc42000b1e0, 0x61ff00, 0x1, 0x1, 0x4)
//        /usr/lib/go/src/testing/testing.go:1002 +0x2d8
//testing.(*M).Run(0xc42003df18, 0xc420083f70)
//        /usr/lib/go/src/testing/testing.go:921 +0x111
//main.main()
//        github.com/DeedleFake/wdte/_test/_testmain.go:46 +0xdb
//
//goroutine 5 [chan receive, 9 minutes]:
//testing.(*T).Run(0xc42008a0f0, 0x55f7e8, 0x8, 0xc42005c240, 0x4b5401)
//        /usr/lib/go/src/testing/testing.go:790 +0x2fc
//github.com/DeedleFake/wdte_test.TestBasics(0xc42008a0f0)
//        $HOME/devel/go/src/github.com/DeedleFake/wdte/wdte_test.go:47 +0x17a
//testing.tRunner(0xc42008a0f0, 0x567ea0)
//        /usr/lib/go/src/testing/testing.go:746 +0xd0
//created by testing.(*T).Run
//        /usr/lib/go/src/testing/testing.go:789 +0x2de
//
//goroutine 7 [GC assist wait]:
//github.com/DeedleFake/wdte.wdte.Frame.New(...)
//        $HOME/devel/go/src/github.com/DeedleFake/wdte/wdte.go:158
//github.com/DeedleFake/wdte.wdte.Frame.WithID(...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:158
//github.com/DeedleFake/wdte/std.Sub(0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, 0xc4204f1710, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:49 +0x99
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, 0xc4204f1710, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, 0xc4204f1710, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, 0xc4204f1710, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111c20, 0xc420111c60, 0x2, 0x2, 0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a300, 0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, 0xc4204f1710, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a300, 0x55ee75, 0x1, 0xc420367af0, 0x1, 0x1, 0xc4204f1710, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc4200a5b80, 0x55ee77, 0x1, 0xc420047440, 0x1, 0x1, 0xc420590fc0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc420047440, 0x1, 0x1, 0xc420590fc0, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc420047440, 0x1, 0x1, 0xc420590fc0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc420047440, 0x1, 0x1, 0xc420590fc0, 0xc420111be0, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc420047440, 0x1, 0x1, 0xc4203e9920, 0xc420111be0, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eb6, 0x1, 0x55ee75, 0x1, 0xc420047440, 0x1, 0x1, 0xc4203e9920, 0xc420111be0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111ba0, 0x55ee75, 0x1, 0xc420047440, 0x1, 0x1, 0xc4203e9920, 0xc420111be0, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111ba0, 0xc420111be0, 0x2, 0x2, 0x55ee75, 0x1, 0xc420047440, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a2a0, 0x55ee75, 0x1, 0xc420047440, 0x1, 0x1, 0xc4203e9920, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a2a0, 0x55ee75, 0x1, 0xc420047440, 0x1, 0x1, 0xc4203e9920, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc4200a5bc0, 0x55ee77, 0x1, 0xc420047460, 0x1, 0x1, 0xc420590f90, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc420047460, 0x1, 0x1, 0xc420590f90, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc420047460, 0x1, 0x1, 0xc420590f90, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc420047460, 0x1, 0x1, 0xc420590f90, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc420047460, 0x1, 0x1, 0xc4203e9bf0, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc420047460, 0x1, 0x1, 0xc4203e9bf0, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc420047460, 0x1, 0x1, 0xc4203e9bf0, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111c20, 0xc420111c60, 0x2, 0x2, 0x55ee75, 0x1, 0xc420047460, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a300, 0x55ee75, 0x1, 0xc420047460, 0x1, 0x1, 0xc4203e9bf0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a300, 0x55ee75, 0x1, 0xc420047460, 0x1, 0x1, 0xc4203e9bf0, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc420229140, 0x55ee77, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc420590f60, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc420590f60, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc420590f60, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc420590f60, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc4204cd6b0, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc4204cd6b0, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc4204cd6b0, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111c20, 0xc420111c60, 0x2, 0x2, 0x55ee75, 0x1, 0xc4200e9250, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a300, 0x55ee75, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc4204cd6b0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a300, 0x55ee75, 0x1, 0xc4200e9250, 0x1, 0x1, 0xc4204cd6b0, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc420538e80, 0x55ee77, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc420590f30, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc420590f30, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc420590f30, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc420590f30, 0xc420111be0, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc42042dce0, 0xc420111be0, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eb6, 0x1, 0x55ee75, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc42042dce0, 0xc420111be0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111ba0, 0x55ee75, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc42042dce0, 0xc420111be0, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111ba0, 0xc420111be0, 0x2, 0x2, 0x55ee75, 0x1, 0xc42059e6d0, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a2a0, 0x55ee75, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc42042dce0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a2a0, 0x55ee75, 0x1, 0xc42059e6d0, 0x1, 0x1, 0xc42042dce0, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc420538ec0, 0x55ee77, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc420590f00, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc420590f00, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc420590f00, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc420590f00, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc4204cc090, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc4204cc090, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc4204cc090, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111c20, 0xc420111c60, 0x2, 0x2, 0x55ee75, 0x1, 0xc42059e6f0, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a300, 0x55ee75, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc4204cc090, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a300, 0x55ee75, 0x1, 0xc42059e6f0, 0x1, 0x1, 0xc4204cc090, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc4200a5ec0, 0x55ee77, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420590ed0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420590ed0, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420590ed0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420590ed0, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420451fb0, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420451fb0, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420451fb0, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111c20, 0xc420111c60, 0x2, 0x2, 0x55ee75, 0x1, 0xc420047af0, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a300, 0x55ee75, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420451fb0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a300, 0x55ee75, 0x1, 0xc420047af0, 0x1, 0x1, 0xc420451fb0, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc420137280, 0x55ee77, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc420590ea0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc420590ea0, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc420590ea0, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc420590ea0, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc4201b5200, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc4201b5200, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc4201b5200, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111c20, 0xc420111c60, 0x2, 0x2, 0x55ee75, 0x1, 0xc4200e9350, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a300, 0x55ee75, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc4201b5200, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a300, 0x55ee75, 0x1, 0xc4200e9350, 0x1, 0x1, 0xc4201b5200, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc420539e40, 0x55ee77, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420590e70, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420590e70, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420590e70, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420590e70, 0xc420111be0, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420494e40, 0xc420111be0, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eb6, 0x1, 0x55ee75, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420494e40, 0xc420111be0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111ba0, 0x55ee75, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420494e40, 0xc420111be0, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//github.com/DeedleFake/wdte.Expr.Call(0x6105c0, 0xc420111ba0, 0xc420111be0, 0x2, 0x2, 0x55ee75, 0x1, 0xc4200e99b0, 0x1, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:309 +0x7c
//github.com/DeedleFake/wdte.(*Expr).Call(0xc42011a2a0, 0x55ee75, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420494e40, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xbc
//github.com/DeedleFake/wdte.FramedFunc.Call(0x610500, 0xc42011a2a0, 0x55ee75, 0x1, 0xc4200e99b0, 0x1, 0x1, 0xc420494e40, 0x55ee77, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:499 +0x7c
//github.com/DeedleFake/wdte.(*FramedFunc).Call(0xc420539e80, 0x55ee77, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420590e40, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0xb2
//github.com/DeedleFake/wdte.Arg.Call(0x0, 0x55ee77, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420590e40, 0x0, 0x0, 0x0, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:485 +0x108
//github.com/DeedleFake/wdte.(*Arg).Call(0x644fa0, 0x55ee77, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420590e40, 0x0, 0x0, 0x0, ...)
//        <autogenerated>:1 +0x88
//github.com/DeedleFake/wdte/std.Sub(0x55ee77, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420590e40, 0xc420111c60, 0x2, 0x2, 0x1, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/std/std.go:51 +0x253
//github.com/DeedleFake/wdte.GoFunc.Call(0x567e98, 0x55ee75, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420495260, 0xc420111c60, 0x2, 0x2, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:247 +0xcb
//github.com/DeedleFake/wdte.Local.Call(0xc4200e98f0, 0xc420115eee, 0x1, 0x55ee75, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420495260, 0xc420111c60, ...)
//        $GOPATH/src/github.com/DeedleFake/wdte/wdte.go:395 +0xd0
//github.com/DeedleFake/wdte.(*Local).Call(0xc420111c20, 0x55ee75, 0x1, 0xc4200e99d0, 0x1, 0x1, 0xc420495260, 0xc420111c60, 0x2, 0x2, ...)
//        <autogenerated>:1 +0xa4
//created by testing.(*T).Run
//        /usr/lib/go/src/testing/testing.go:789 +0x2de
//
//rax    0x393b8466e87
//rbx    0xfffffffffffffade
//rcx    0x17
//rdx    0x3acfd687
//rdi    0x1
//rsi    0x7fff8cdd7f08
//rbp    0x7fff8cdd7f68
//rsp    0x7fff8cdd7f28
//r8     0x10
//r9     0xc420050670
//r10    0x7fff8cdd7f08
//r11    0x1
//r12    0x0
//r13    0x0
//r14    0x4563e0
//r15    0x0
//rip    0x42dcc4
//rflags 0x212
//cs     0x33
//fs     0x0
//gs     0x0
//*** Test killed with quit: ran too long (10m0s).
//FAIL    github.com/DeedleFake/wdte      600.016s
