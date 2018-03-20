package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/repl"
	"github.com/DeedleFake/wdte/std"
	"github.com/peterh/liner"
	"golang.org/x/crypto/ssh/terminal"
)

func printRet(ret wdte.Func) {
	switch ret := ret.(type) {
	case error, fmt.Stringer:
		fmt.Printf(": %v\n", ret)
		return

	case wdte.GoFunc:
		fmt.Println(": complex value (GoFunc)")
		return
	}

	switch k := reflect.Indirect(reflect.ValueOf(ret)).Kind(); k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Struct, reflect.UnsafePointer:
		fmt.Printf(": complex value (%v)\n", k)

	default:
		fmt.Printf(": %v\n", ret)
	}
}

func stdin(im wdte.Importer) {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		file(im, os.Stdin)
		return
	}

	lr := liner.NewLiner()
	lr.SetCtrlCAborts(true)
	defer lr.Close()

	const (
		modeTop = ">>> "
		modeSub = "... "
	)

	mode := modeTop
	next := func() ([]byte, error) {
		line, err := lr.Prompt(mode)
		if err == nil {
			lr.AppendHistory(line)
		}
		return []byte(line + "\n"), err
	}

	r := repl.New(next, im, std.Scope)

	for {
		ret, err := r.Next()
		if err != nil {
			if err == repl.ErrIncomplete {
				mode = modeSub
				continue
			}

			if err == liner.ErrPromptAborted {
				r.Cancel()
				mode = modeTop
				continue
			}

			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue
		}
		if ret == nil {
			break
		}

		printRet(ret.Call(wdte.F().WithScope(r.Scope)))

		mode = modeTop
	}
}
