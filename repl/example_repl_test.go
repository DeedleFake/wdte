package repl_test

import (
	"fmt"
	"os"
	"reflect"

	"github.com/DeedleFake/wdte/repl"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/all"
	"github.com/peterh/liner"
)

var (
	mode = ">>> "
)

func next(lr *liner.State) repl.NextFunc {
	return func() ([]byte, error) {
		line, err := lr.Prompt(mode)
		return []byte(line + "\n"), err
	}
}

func Example() {
	lr := liner.NewLiner()
	lr.SetCtrlCAborts(true)
	defer lr.Close()

	r := repl.New(next(lr), std.Import, std.S())

	for {
		ret, err := r.Next()
		if err != nil {
			if err == repl.ErrIncomplete {
				mode = "... "
				continue
			}

			if err == liner.ErrPromptAborted {
				r.Cancel()
				mode = ">>> "
				continue
			}

			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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

		mode = ">>> "
	}
}
