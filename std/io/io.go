// Package io contains WDTE functions for dealing with files and other
// types of data streams.
package io

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/DeedleFake/wdte"
)

type reader interface {
	wdte.Func
	io.Reader
}

// Reader wraps an io.Reader, allowing it to be used as a WDTE function.
type Reader struct {
	io.Reader
}

func (r Reader) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return r
}

type writer interface {
	wdte.Func
	io.Writer
}

// Writer wraps an io.Writer, allowing it to be used as a WDTE function.
type Writer struct {
	io.Writer
}

func (w Writer) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return w
}

func Writeln(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("writeln")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Writeln)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, more ...wdte.Func) wdte.Func {
			return Writeln(frame, append(more, args...)...)
		})
	}

	w := args[0].Call(frame).(writer)
	d := args[1].Call(frame)
	_, err := fmt.Fprintln(w, d)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return w
}

func String(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("string")

	if len(args) == 0 {
		return wdte.GoFunc(String)
	}

	r := args[0].Call(frame).(reader)

	var buf bytes.Buffer
	_, err := io.Copy(&buf, r)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return wdte.String(buf.String())
}

type scanner struct {
	s *bufio.Scanner
}

func (s scanner) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return s
}

func (s scanner) Next(frame wdte.Frame) (wdte.Func, bool) {
	ok := s.s.Scan()
	if !ok {
		err := s.s.Err()
		if err != nil {
			return wdte.Error{Err: err, Frame: frame}, false
		}
		return nil, false
	}

	return wdte.String(s.s.Text()), true
}

func Lines(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("lines")

	if len(args) == 0 {
		return wdte.GoFunc(Lines)
	}

	r := args[0].Call(frame).(reader)
	return scanner{s: bufio.NewScanner(r)}
}

// Module returns a module for easy importing into an actual script.
// The imported functions have the same names as the functions in this
// package, except that the first letter is lowercase.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"stdin":  Reader{Reader: os.Stdin},
			"stdout": Writer{Writer: os.Stdout},
			"stderr": Writer{Writer: os.Stderr},

			"writeln": wdte.GoFunc(Writeln),

			"string": wdte.GoFunc(String),
			"lines":  wdte.GoFunc(Lines),
		},
	}
}
