package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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

func main() {
	blacklist := flag.String("blacklist", "", "Comma-separated list of modules that can't be imported.")
	eval := flag.String("e", "", "An expression to evaluate instead of reading from a file.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <file> | -\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	im := importer(strings.Split(*blacklist, ","))

	if *eval != "" {
		file(im, strings.NewReader(*eval))
		return
	}

	inpath := flag.Arg(0)
	switch inpath {
	case "", "-":
		stdin(im)

	default:
		f, err := os.Open(inpath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open %q: %v", inpath, err)
			os.Exit(1)
		}
		defer f.Close()

		file(im, f)
	}
}
