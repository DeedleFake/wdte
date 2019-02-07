package main

import (
	"fmt"
	"io"
	"os"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

func file(im wdte.Importer, file io.Reader) {
	m, err := wdte.Parse(file, im, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse script: %v", err)
		os.Exit(1)
	}

	ret := m.Call(std.F())
	if err, ok := ret.(error); ok {
		fmt.Fprintf(os.Stderr, "Script returned an error: %v", err)
		os.Exit(3)
	}
}
