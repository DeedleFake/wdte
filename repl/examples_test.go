package repl_test

import (
	"fmt"
	"strings"

	"github.com/DeedleFake/wdte/repl"
)

func ExamplePartial() {
	stack, partial := repl.Partial(strings.NewReader("let io =>"), nil)
	fmt.Println(partial)

	stack, partial = repl.Partial(strings.NewReader("import 'io'"), stack)
	fmt.Println(partial)

	stack, partial = repl.Partial(strings.NewReader(";"), stack)
	fmt.Println(partial)
	// Output: true
	// true
	// false
}
