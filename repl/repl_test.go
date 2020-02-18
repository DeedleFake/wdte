package repl_test

import (
	"bufio"
	"fmt"
	"strings"
	"testing"

	"github.com/DeedleFake/wdte/repl"
	"github.com/DeedleFake/wdte/std"
)

func TestREPL(t *testing.T) {
	src := bufio.NewReader(strings.NewReader(`let test [a b] => + a b;`))
	r := repl.New(func() ([]byte, error) {
		return src.ReadBytes('\n')
	}, std.Import, nil, nil)

	ret, err := r.Next()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%#v\n", ret)
}
