package std

import (
	"fmt"

	"github.com/DeedleFake/wdte"
)

var (
	// Import provides a simple importer that imports registered
	// modules.
	Import = wdte.ImportFunc(stdImporter)

	modules = make(map[string]*wdte.Scope)
)

func stdImporter(from string) (*wdte.Scope, error) {
	if m, ok := modules[from]; ok {
		return m, nil
	}

	return nil, fmt.Errorf("Unknown import: %q", from)
}

// Register registers a module for importing by Import.
func Register(name string, module *wdte.Scope) {
	modules[name] = module
}
