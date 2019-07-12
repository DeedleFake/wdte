package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

func wdteVersion() (string, error) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("failed to read build info")
	}

	for _, dep := range bi.Deps {
		if dep.Path != "github.com/DeedleFake/wdte" {
			continue
		}

		return dep.Version, nil
	}

	return bi.Main.Version, nil
}

func main() {
	blacklist := flag.String("blacklist", "", "Comma-separated list of modules that can't be imported.")
	eval := flag.String("e", "", "An expression to evaluate instead of reading from a file.")
	version := flag.Bool("version", false, "Print the Go and WDTE versions and then exit.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] [<file> | -] [arguments...]\n\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		g := runtime.Version()
		w, err := wdteVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to get WDTE version: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Go: %v\n", g)
		fmt.Printf("WDTE: %v\n", w)
		return
	}

	im := importer("", strings.Split(*blacklist, ","), flag.Args(), nil)

	if *eval != "" {
		file(im, strings.NewReader(*eval))
		return
	}

	inpath := flag.Arg(0)
	switch inpath {
	case "", "-":
		stdin(im, nil)

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
