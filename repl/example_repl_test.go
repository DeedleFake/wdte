package repl_test

import (
	"fmt"
	"os"
	"reflect"

	"github.com/DeedleFake/wdte/repl"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/all"
)

func Example() {
	r := repl.New(os.Stdin, std.Import, std.S())

	for {
		fmt.Printf("> ")
		ret, err := r.Next()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		if ret == nil {
			break
		}

		switch reflect.Indirect(reflect.ValueOf(ret)).Kind() {
		case reflect.Struct:
			fmt.Printf(": complex value\n")

		default:
			fmt.Printf(": %v\n", ret)
		}
	}
}
