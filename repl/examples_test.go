package repl_test

import (
	"fmt"
	"strings"

	"github.com/DeedleFake/wdte/repl"
)

func ExamplePartial() {
	stack, partial := repl.Partial(strings.NewReader("let io =>"), nil, nil)
	fmt.Println(partial)

	stack, partial = repl.Partial(strings.NewReader("import 'io'"), stack, nil)
	fmt.Println(partial)

	_, partial = repl.Partial(strings.NewReader(";"), stack, nil)
	fmt.Println(partial)
	// Output: true
	// true
	// false
}
