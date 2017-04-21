package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v <file.ebnf>\n", os.Args[0])
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %q: %v", flag.Arg(0), err)
		os.Exit(1)
	}
	defer file.Close()

	g, err := LoadGrammar(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading grammar: %v", err)
		os.Exit(1)
	}

	err = tmpl.Execute(os.Stderr, g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error printing code: %v", err)
		os.Exit(1)
	}
}
