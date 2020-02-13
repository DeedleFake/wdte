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
		out  []scanner.Token
	}{
		{
			name: "Simple",
			in: `"test" => test;
f => + test.other 3 -5.2;
# This is a comment.
o => print "double\n" 'single\\';`,
			out: []scanner.Token{
				{Type: scanner.String, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.ID, Val: "f"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.ID, Val: "+"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "."},
				{Type: scanner.ID, Val: "other"},
				{Type: scanner.Number, Val: float64(3)},
				{Type: scanner.Number, Val: float64(-5.2)},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.ID, Val: "o"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.ID, Val: "print"},
				{Type: scanner.String, Val: "double\n"},
				{Type: scanner.String, Val: "single\\"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF, Val: nil},
			},
		},
		{
			name: "Switch",
			in:   `test => n {case => r;};`,
			out: []scanner.Token{
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.ID, Val: "n"},
				{Type: scanner.Keyword, Val: "{"},
				{Type: scanner.ID, Val: "case"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.ID, Val: "r"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: "}"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF, Val: nil},
			},
		},
		{
			name: "Compound",
			in:   `test => (a); test => (a;);`,
			out: []scanner.Token{
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.Keyword, Val: "("},
				{Type: scanner.ID, Val: "a"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: ")"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.Keyword, Val: "("},
				{Type: scanner.ID, Val: "a"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: ")"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF, Val: nil},
			},
		},
		{
			name: "Compound/Collect",
			in:   `test => (|a|); test => (|a;|);`,
			out: []scanner.Token{
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.Keyword, Val: "(|"},
				{Type: scanner.ID, Val: "a"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: "|)"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.Keyword, Val: "(|"},
				{Type: scanner.ID, Val: "a"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: "|)"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF, Val: nil},
			},
		},
		{
			name: "FuncMods",
			in:   `memo test => ();`,
			out: []scanner.Token{
				{Type: scanner.Keyword, Val: "memo"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.Keyword, Val: "("},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: ")"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF, Val: nil},
			},
		},
		{
			name: "MaybeNumber",
			in:   `-5 .3 . -`,
			out: []scanner.Token{
				{Type: scanner.Number, Val: float64(-5)},
				{Type: scanner.Number, Val: .3},
				{Type: scanner.Keyword, Val: "."},
				{Type: scanner.ID, Val: "-"},
				{Type: scanner.EOF, Val: nil},
			},
		},
		{
			name: "LongSymbol",
			in:   `(@`,
			out: []scanner.Token{
				{Type: scanner.Keyword, Val: "(@"},
				{Type: scanner.EOF},
			},
		},
		{
			name: "LetExpr",
			in:   `let add x y => + x y;`,
			out: []scanner.Token{
				{Type: scanner.Keyword, Val: "let"},
				{Type: scanner.ID, Val: "add"},
				{Type: scanner.ID, Val: "x"},
				{Type: scanner.ID, Val: "y"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.ID, Val: "+"},
				{Type: scanner.ID, Val: "x"},
				{Type: scanner.ID, Val: "y"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF},
			},
		},
		{
			name: "LetPattern",
			in:   `let [a b] => [3; 5];`,
			out: []scanner.Token{
				{Type: scanner.Keyword, Val: "let"},
				{Type: scanner.Keyword, Val: "["},
				{Type: scanner.ID, Val: "a"},
				{Type: scanner.ID, Val: "b"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: "]"},
				{Type: scanner.Keyword, Val: "=>"},
				{Type: scanner.Keyword, Val: "["},
				{Type: scanner.Number, Val: float64(3)},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Number, Val: float64(5)},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.Keyword, Val: "]"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF},
			},
		},
		{
			name: "Macro",
			in:   `@fmt[{q}, 'greetings'];`,
			out: []scanner.Token{
				{Type: scanner.Macro, Val: [2]string{"fmt", "{q}, 'greetings'"}},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.EOF},
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

				assertTokensEqual(t, test.out[i], s.Tok())

				t.Logf("Got %#v", s.Tok())
			}
			if err := s.Err(); err != nil {
				t.Errorf("Scanner error: %v", err)
			}
		})
	}
}

func assertTokensEqual(t *testing.T, ex scanner.Token, got scanner.Token) {
	if ex.Type != got.Type {
		t.Errorf("Unexpected token type:")
		t.Errorf("\tExpected %#v.", ex)
		t.Errorf("\tGot %#v.", got)
		return
	}

	if ex.Val != got.Val {
		t.Errorf("Unexpected token value:")
		t.Errorf("\tExpected %#v.", ex)
		t.Errorf("\tGot %#v.", got)
		return
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
