package wdte_test

import (
	"strings"
	"testing"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/std/io"
)

type twriter struct {
	t *testing.T
}

func (w twriter) Write(data []byte) (int, error) {
	w.t.Logf("%s", data)
	return len(data), nil
}

func TestModule(t *testing.T) {
	const test = `
'stream' => s;
'io' => io;
'io/file' => file;

memo fib n => switch n {
	== 0 => 0;
	== 1 => 1;
	default => + (fib (- n 1)) (fib (- n 2));
};

memo fact n => switch n {
	<= 1 => 1;
	default => - n 1 -> fact -> * n;
};

main w r => (
	s.range 15
	-> s.map fib
	-> s.collect
	-> io.writeln w;

	s.new [5; 2; fib 7]
	-> s.map (+ 2)
	-> s.collect
	-> io.writeln w;

	fact 5 -> io.writeln w;

	w
	-> io.write 'This is a test.'
	-> io.writeln 'Or is it?';

	r
	-> io.lines
	-> s.map (io.writeln w)
	-> s.collect;

	file.open 'wdte_test.go'
	-> io.copy w
	-> io.close;

	#io.readString 'This is also a test.'
	#-> io.copy w;
);
`

	m, err := std.Module().Parse(strings.NewReader(test), std.Import)
	if err != nil {
		t.Fatal(err)
	}

	//t.Log("Imports:")
	//for i := range m.Imports {
	//	t.Logf("\t%q", i)
	//}

	//t.Log("Funcs:")
	//for f := range m.Funcs {
	//	t.Logf("\t%q", f)
	//}

	w := twriter{t: t}
	r := strings.NewReader(test)
	if err, ok := m.Funcs["main"].Call(wdte.F(), io.Writer{w}, io.Reader{r}).(error); ok {
		t.Fatal(err)
	}
}
