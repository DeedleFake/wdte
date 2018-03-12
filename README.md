wdte
====

[![GoDoc](https://godoc.org/github.com/DeedleFake/wdte?status.svg)](https://godoc.org/github.com/DeedleFake/wdte)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeedleFake/wdte)](https://goreportcard.com/report/github.com/DeedleFake/wdte)

WDTE is a simple, functional-ish, embedded scripting language.

Why does this exist?
--------------------

Good question. In fact, I found myself asking the same thing, hence the name.

I had a number of design goals in mind when I started working on this project:

* Extremely simple. Entire grammar is less than 20-30 lines of specification.
* Grammar is LL(1) parseable.
* Functional-ish, but not particularly strict about it.
* Designed primarily for embedding. No command-line interpreter by default.
* Extremely easy to use from the binding side. In this case, that's primarily Go.

If you want to try the language yourself, feel free to take a look at [the playground][playground]. It shows not only some of the features of the language in terms of actually writing code in it, but also how embeddable it is. The playground runs entirely in the browser *on the client's end* thanks to [GopherJS][gopherjs].

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
let i => import 'some/import/path/or/another';

i.print 3;
+ 5 2 -> i.print;
`

func im(from string) (*wdte.Scope, error) {
	var print wdte.GoFunc
	print = wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		if len(args) < 1 {
			return print
		}

		a := args[0].Call(frame)
		fmt.Println(a)
		return a
	})

	return wdte.S().Map(map[wdte.ID]wdte.Func{
		"print": print,
	}), nil
}

func main() {
	m, err := wdte.Parse(strings.NewReader(src), wdte.ImportFunc(im))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing module: %v\n", err)
		os.Exit(1)
	}

	scope := wdte.S().Sub("+", wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		var sum wdte.Number
		for _, arg := range args {
			sum += arg.Call(frame).(wdte.Number)
		}
		return sum
	}))

	m.Call(wdte.F().WithScope(scope))
}
```

##### Output

```
3
7
```

Documentation
-------------

For an overview of the language's design and features, see [the GitHub wiki][wiki].

Status
------

WDTE is in a pre-alpha state. It is filled with bugs and large amounts of stuff are subject to change without warning. That being said, if you're interested in anything, feel free to submit a pull request and get things fixed and/or implemented faster.

[playground]: https://deedlefake.github.io/wdte
[gopherjs]: https://github.com/gopherjs/gopherjs
[wiki]: https://github.com/DeedleFake/wdte/wiki
