package scanner_test

import (
	"reflect"
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
			in: `"test" -> test;
f = + test.other 3 -5.2;
o = print "double\n" 'single\\';`,
			out: []scanner.Token{
				{Type: scanner.String, Val: "test"},
				{Type: scanner.Keyword, Val: "->"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.ID, Val: "f"},
				{Type: scanner.Keyword, Val: "="},
				{Type: scanner.ID, Val: "+"},
				{Type: scanner.ID, Val: "test"},
				{Type: scanner.Keyword, Val: "."},
				{Type: scanner.ID, Val: "other"},
				{Type: scanner.Number, Val: float64(3)},
				{Type: scanner.Number, Val: float64(-5.2)},
				{Type: scanner.Keyword, Val: ";"},
				{Type: scanner.ID, Val: "o"},
				{Type: scanner.Keyword, Val: "="},
				{Type: scanner.ID, Val: "print"},
				{Type: scanner.String, Val: "double\n"},
				{Type: scanner.String, Val: "single\\"},
				{Type: scanner.Keyword, Val: ";"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scanner.New(strings.NewReader(test.in))
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
	if reflect.TypeOf(ex) != reflect.TypeOf(got) {
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
