package main

import (
	"fmt"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/all"
)

// TODO: This should be able to import other scripts.
func importer(blacklist []string) wdte.Importer {
	return wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
		for _, m := range blacklist {
			if from == m {
				return nil, fmt.Errorf("%q is blacklisted", from)
			}
		}

		return std.Import(from)
	})
}
