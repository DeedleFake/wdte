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
			in:   `"test" -> test; f = + test.other;`,
			out: []scanner.Token{
				&scanner.String{Val: "test"},
				&scanner.Keyword{Val: "->"},
				&scanner.ID{Val: "test"},
				&scanner.Keyword{Val: ";"},
				&scanner.ID{Val: "f"},
				&scanner.Keyword{Val: "="},
				&scanner.ID{Val: "+"},
				&scanner.ID{Val: "test"},
				&scanner.Keyword{Val: "."},
				&scanner.ID{Val: "other"},
				&scanner.Keyword{Val: ";"},
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

func tokenVal(t scanner.Token) (val interface{}) {
	switch t := t.(type) {
	case *scanner.Number:
		val = t.Val
	case *scanner.String:
		val = t.Val
	case *scanner.Nil:
	case *scanner.ID:
		val = t.Val
	case *scanner.Keyword:
		val = t.Val
	}

	return
}

func assertTokensEqual(t *testing.T, ex scanner.Token, got scanner.Token) {
	if reflect.TypeOf(ex) != reflect.TypeOf(got) {
		t.Errorf("Unexpected token type:")
		t.Errorf("\tExpected %#v.", ex)
		t.Errorf("\tGot %#v.", got)
		return
	}

	if tokenVal(ex) != tokenVal(got) {
		t.Errorf("Unexpected token value:")
		t.Errorf("\tExpected %#v.", ex)
		t.Errorf("\tGot %#v.", got)
		return
	}
}
