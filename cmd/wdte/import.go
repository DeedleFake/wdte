package main

import (
	"fmt"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/all"
)

// TODO: This should be able to import other scripts.
func importer(blacklist []string, args []string) wdte.Importer {
	wargs := make(wdte.Array, 0, len(args))
	for _, arg := range args {
		wargs = append(wargs, wdte.String(arg))
	}

	cliScope := wdte.S().Map(map[wdte.ID]wdte.Func{
		"args": wargs,
	})

	return wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
		for _, m := range blacklist {
			if from == m {
				return nil, fmt.Errorf("%q is blacklisted", from)
			}
		}

		if from == "cli" {
			return cliScope, nil
		}

		return std.Import(from)
	})
}
