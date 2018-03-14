package main

import (
	"os"

	"github.com/DeedleFake/wdte"
)

func importPlugin(from string, im wdte.Importer) (*wdte.Scope, error) {
	return nil, os.ErrNotExist
}
