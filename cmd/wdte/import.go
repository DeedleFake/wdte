package main

import (
	"fmt"
	"os"
	"path/filepath"
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

	std.Register("cli", cliScope)

	return wdte.ImportFunc(func(from string) (*wdte.Scope, error) {
		for _, m := range blacklist {
			if from == m {
				return nil, fmt.Errorf("%q is blacklisted", from)
			}
		}

		if strings.HasPrefix(from, ".") || strings.HasPrefix(from, "/") {
			path := filepath.FromSlash(from)
			im := importer(filepath.Dir(path), blacklist, args)

			s, err := importPlugin(filepath.Join(wd, path), im)
			if !os.IsNotExist(err) {
				return s, err
			}

			s, err = importScript(filepath.Join(wd, path), im)
			if !os.IsNotExist(err) {
				return s, err
			}
		}

		return std.Import(from)
	})
}
