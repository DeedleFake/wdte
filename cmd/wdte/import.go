package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
	_ "github.com/DeedleFake/wdte/std/all"
)

func importScript(from string, im wdte.Importer) (*wdte.Scope, error) {
	file, err := os.Open(from + ".wdte")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	c, err := wdte.Parse(file, im)
	if err != nil {
		return nil, err
	}

	s, _ := c.Collect(std.F())
	return s, nil
}

func importer(wd string, blacklist []string, args []string) wdte.Importer {
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

		if strings.HasPrefix(from, ".") || strings.HasPrefix(from, "/") {
			im := importer(path.Dir(from), blacklist, args)

			s, err := importPlugin(from, im)
			if !os.IsNotExist(err) {
				return s, err
			}

			s, err = importScript(from, im)
			if !os.IsNotExist(err) {
				return s, err
			}
		}

		if from == "cli" {
			return cliScope, nil
		}

		return std.Import(from)
	})
}
