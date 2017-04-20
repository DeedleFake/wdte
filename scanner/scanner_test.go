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
				scanner.String{Val: "test"},
				scanner.Keyword{Val: "->"},
				scanner.ID{Val: "test"},
				scanner.Keyword{Val: ";"},
				scanner.ID{Val: "f"},
				scanner.Keyword{Val: "="},
				scanner.ID{Val: "+"},
				scanner.ID{Val: "test"},
				scanner.Keyword{Val: "."},
				scanner.ID{Val: "other"},
				scanner.Keyword{Val: ";"},
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
			}
		})
	}
}

func assertTokensEqual(t *testing.T, ex scanner.Token, got scanner.Token) {
	if reflect.TypeOf(ex) != reflect.TypeOf(got) {
		t.Errorf("Expected %T. Got %T.", ex, got)
		return
	}

	if !reflect.DeepEqual(ex, got) {
		t.Errorf("Expected %#v. Got %#v.", ex, got)
	}
}
