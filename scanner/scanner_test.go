package scanner_test

import (
	"io"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte/scanner"
)

func TestScanner(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  []scanner.TokenValue
	}{
		{
			name: "Simple",
			in: `"test" => test;
f => + test.other 3 -5.2;
# This is a comment.
o => print "double\n" 'single\\';`,
			out: []scanner.TokenValue{
				scanner.String("test"),
				scanner.Keyword("=>"),
				scanner.ID("test"),
				scanner.Keyword(";"),
				scanner.ID("f"),
				scanner.Keyword("=>"),
				scanner.ID("+"),
				scanner.ID("test"),
				scanner.Keyword("."),
				scanner.ID("other"),
				scanner.Number(float64(3)),
				scanner.Number(float64(-5.2)),
				scanner.Keyword(";"),
				scanner.ID("o"),
				scanner.Keyword("=>"),
				scanner.ID("print"),
				scanner.String("double\n"),
				scanner.String("single\\"),
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
		{
			name: "Switch",
			in:   `test => n {case => r;};`,
			out: []scanner.TokenValue{
				scanner.ID("test"),
				scanner.Keyword("=>"),
				scanner.ID("n"),
				scanner.Keyword("{"),
				scanner.ID("case"),
				scanner.Keyword("=>"),
				scanner.ID("r"),
				scanner.Keyword(";"),
				scanner.Keyword("}"),
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
		{
			name: "Compound",
			in:   `test => (a); test => (a;);`,
			out: []scanner.TokenValue{
				scanner.ID("test"),
				scanner.Keyword("=>"),
				scanner.Keyword("("),
				scanner.ID("a"),
				scanner.Keyword(";"),
				scanner.Keyword(")"),
				scanner.Keyword(";"),
				scanner.ID("test"),
				scanner.Keyword("=>"),
				scanner.Keyword("("),
				scanner.ID("a"),
				scanner.Keyword(";"),
				scanner.Keyword(")"),
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
		{
			name: "FuncMods",
			in:   `memo test => ();`,
			out: []scanner.TokenValue{
				scanner.Keyword("memo"),
				scanner.ID("test"),
				scanner.Keyword("=>"),
				scanner.Keyword("("),
				scanner.Keyword(";"),
				scanner.Keyword(")"),
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
		{
			name: "MaybeNumber",
			in:   `-5 .3 . -`,
			out: []scanner.TokenValue{
				scanner.Number(float64(-5)),
				scanner.Number(.3),
				scanner.Keyword("."),
				scanner.ID("-"),
				scanner.EOF{},
			},
		},
		{
			name: "LongSymbol",
			in:   `(@`,
			out: []scanner.TokenValue{
				scanner.Keyword("(@"),
				scanner.EOF{},
			},
		},
		{
			name: "LetExpr",
			in:   `let add x y => + x y;`,
			out: []scanner.TokenValue{
				scanner.Keyword("let"),
				scanner.ID("add"),
				scanner.ID("x"),
				scanner.ID("y"),
				scanner.Keyword("=>"),
				scanner.ID("+"),
				scanner.ID("x"),
				scanner.ID("y"),
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
		{
			name: "LetPattern",
			in:   `let [a b] => [3; 5];`,
			out: []scanner.TokenValue{
				scanner.Keyword("let"),
				scanner.Keyword("["),
				scanner.ID("a"),
				scanner.ID("b"),
				scanner.Keyword(";"),
				scanner.Keyword("]"),
				scanner.Keyword("=>"),
				scanner.Keyword("["),
				scanner.Number(float64(3)),
				scanner.Keyword(";"),
				scanner.Number(float64(5)),
				scanner.Keyword(";"),
				scanner.Keyword("]"),
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
		{
			name: "Macro",
			in:   `@fmt[{q}, 'greetings'];`,
			out: []scanner.TokenValue{
				scanner.Macro{"fmt", "{q}, 'greetings'"},
				scanner.Keyword(";"),
				scanner.EOF{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scanner.New(strings.NewReader(test.in), nil)
			for i := 0; s.Scan(); i++ {
				if i >= len(test.out) {
					t.Fatalf("Extra token: %#v", s.Tok())
				}

				if test.out[i] != s.Tok().Val {
					t.Errorf("Unexpected token value:")
					t.Errorf("\tExpected %#v.", test.out[i])
					t.Errorf("\tGot %#v.", s.Tok().Val)
				}

				t.Logf("Got %#v", s.Tok())
			}
			if err := s.Err(); err != nil {
				t.Errorf("Scanner error: %v", err)
			}
		})
	}
}

var (
	// These are for examples to use.
	r io.Reader
)

func ExampleScanner() {
	s := scanner.New(r, nil)
	for s.Scan() {
		/* Do something with s.Tok(). */
	}
	if err := s.Err(); err != nil {
		panic(err)
	}
}
