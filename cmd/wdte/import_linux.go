package main

import (
	"errors"
	"os"
	"plugin"

	"github.com/DeedleFake/wdte"
)

func importPlugin(from string, im wdte.Importer) (*wdte.Scope, error) {
	_, err := os.Stat(from + ".so")
	if err != nil {
		return nil, err
	}

	p, err := plugin.Open(from + ".so")
	if err != nil {
		return nil, err
	}

	init, err := p.Lookup("S")
	if err != nil {
		return nil, err
	}

	s, ok := init.(func() *wdte.Scope)
	if !ok {
		return nil, errors.New("S symbol has wrong type")
	}

	return s(), nil
}
