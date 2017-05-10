wdte
====

[![GoDoc](https://godoc.org/github.com/DeedleFake/wdte?status.svg)](https://godoc.org/github.com/DeedleFake/wdte)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeedleFake/wdte)](https://goreportcard.com/report/github.com/DeedleFake/wdte)

WDTE ('why the') is a simple, functional-ish, embedded scripting language.

Why does this exist?
--------------------

Good question. In fact, I found myself asking the same thing, hence the name.

I had a number of design goals in mind when I started working on this project:

* Extremely simple. Entire grammar is less than 20-30 lines of EBNF specification.
* Grammar is LL(1) parseable.
* Functional-ish, but not particularly strict about it.
* Designed primarily for embedding. No command-line interpreter by default.
* Extremely easy to use from the binding side. In this case, that's primarily Go.

Example
-------

```
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DeedleFake/wdte"
)

const src = `
'import' => i;

main => (
	i.print 3;
	+ 5 2 -> i.print;
);
`

func im(from string) (*wdte.Module, error) {
	var print wdte.GoFunc
	print = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		if len(args) < 1 {
			return print
		}

		a := args[0].Call(frame)
		fmt.Println(a)
		return a
	})

	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"print": print,
		},
	}, nil
}

func main() {
	m, err := wdte.Parse(strings.NewReader(src), wdte.ImportFunc(im))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing module: %v\n", err)
		os.Exit(1)
	}
	m.Funcs["+"] = wdte.GoFunc(func(frame []wdte.Func, args ...wdte.Func) wdte.Func {
		var sum wdte.Number
		for _, arg := range args {
			sum += arg.Call(frame).(wdte.Number)
		}
		return sum
	})

	m.Funcs["main"].Call(nil)
}
```

##### Output

```
5
7
```

Documentation
-------------

For an overview of the languages design and features, see [the GitHub wiki](https://github.com/DeedleFake/wdte/wiki/Overview).

Status
------

WDTE is in a very, very pre-alpha state. It is filled with bugs, parts of it are simply not implemented yet, and large amounts of stuff are subject to change without warning. That being said, if you're interested in anything, feel free to submit a pull request and get things fixed and/or implemented faster.
