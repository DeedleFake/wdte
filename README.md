wdte
====

[![GoDoc](https://godoc.org/github.com/DeedleFake/wdte?status.svg)](https://godoc.org/github.com/DeedleFake/wdte)
[![Go Report Card](https://goreportcard.com/badge/github.com/DeedleFake/wdte)](https://goreportcard.com/report/github.com/DeedleFake/wdte)
[![cover.run](https://cover.run/go/github.com/DeedleFake/wdte.svg?style=flat&tag=golang-1.10)](https://cover.run/go?tag=golang-1.10&repo=github.com/DeedleFake/wdte)

WDTE is a simple, functional-ish, embedded scripting language.

Why does this exist?
--------------------

Good question. In fact, I found myself asking the same thing, hence the name.

I had a number of design goals in mind when I started working on this project:

* Extremely simple. Entire grammar is less than 20-30 lines of specification.
* Grammar is LL(1) parseable.
* Functional-ish, but not particularly strict about it.
* Designed primarily for embedding.
* Extremely easy to use from the binding side. In this case, that's primarily Go.

If you want to try the language yourself, feel free to take a look at [the playground][playground]. It shows not only some of the features of the language in terms of actually writing code in it, but also how embeddable it is. The playground runs entirely in the browser *on the client's end* thanks to WebAssembly.

Example
-------

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/wdteutil"
)

const src = `
let i => import 'some/import/path/or/another';

i.print 3;
+ 5 2 -> i.print;
7 -> + 5 -> i.print;
`

func im(from string) (*wdte.Scope, error) {
	return wdte.S().Map(map[wdte.ID]wdte.Func{
		"print": wdteutil.Func("print", func(v interface{}) interface{} {
		fmt.Println(v)
		return v
	}),
	}), nil
}

func Sum(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("+")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Sum), args...)
	}

	var sum wdte.Number
	for _, arg := range args {
		sum += arg.(wdte.Number)
	}
	return sum
}

func main() {
	m, err := wdte.Parse(strings.NewReader(src), wdte.ImportFunc(im), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing script: %v\n", err)
		os.Exit(1)
	}

	scope := wdte.S().Add("+", wdte.GoFunc(Sum))

	r := m.Call(wdte.F().WithScope(scope))
	if err, ok := r.(error); ok {
		fmt.Fprintf(os.Stderr, "Error running script: %v\n", err)
		os.Exit(1)
	}
}
```

##### Output

```
3
7
12
```

Documentation
-------------

For an overview of the language's design and features, see [the GitHub wiki][wiki].

Status
------

WDTE is in a pre-alpha state. It is filled with bugs and large amounts of stuff are subject to change without warning. That being said, if you're interested in anything, feel free to submit a pull request and get things fixed and/or implemented faster.

[playground]: https://deedlefake.github.io/wdte
[wiki]: https://github.com/DeedleFake/wdte/wiki
