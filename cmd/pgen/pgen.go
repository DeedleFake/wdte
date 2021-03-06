package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <file.cfg>\n", os.Args[0])
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}
	output := flag.String("out", "", "File to output to, or stdout if blank.")
	detectAmbiguity := flag.Bool("detectAmbiguity", true, "Give warnings for grammar ambiguities.")
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

	g, err := LoadGrammar(file, *detectAmbiguity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading grammar: %v", err)
		os.Exit(1)
	}

	out := &formatter{w: os.Stdout}
	if *output != "" {
		file, err := os.Create(*output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %q: %v", *output, err)
			os.Exit(1)
		}
		defer file.Close()

		out.w = file
	}

	err = tmpl.Execute(out, g)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error printing code: %v", err)
		os.Exit(1)
	}

	err = out.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting code: %v", err)
		os.Exit(1)
	}
}
