// Package io contains WDTE functions for dealing with files and other
// types of data streams.
package io

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

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

func (r Reader) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
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

func (w Writer) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return w
}

// Seek seeks an io.Seeker and then returns it.
func Seek(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("seek")

	if len(args) < 3 {
		return wdte.GoFunc(func(frame wdte.Frame, more ...wdte.Func) wdte.Func {
			return Seek(frame, append(more, args...)...)
		})
	}

	s := args[0].Call(frame).(io.Seeker)
	off := int64(args[1].Call(frame).(wdte.Number))
	rel := args[2].Call(frame).(wdte.Number)

	var w int
	switch {
	case rel < 0:
		w = io.SeekEnd
	case rel == 0:
		w = io.SeekCurrent
	case rel > 0:
		w = io.SeekStart
	}

	_, err := s.Seek(off, w)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return s.(wdte.Func)
}

// Close closes a closer. This includes files opened with other
// functions in this module.
func Close(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("close")

	if len(args) == 0 {
		return wdte.GoFunc(Close)
	}

	c := args[0].Call(frame).(io.Closer)
	err := c.Close()
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return c.(wdte.Func)
}

// Combine combines multiple readers or multiple writers. If the
// arguments passed are readers, it uses Go's io.MultiReader to
// concatenate them. If the arguments passed are writers, it uses Go's
// io.MultiWriter to combine them. If only one argument is given, it
// returns a function which combines its arguments with the argument
// originally given.
func Combine(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("combine")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Combine)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, more ...wdte.Func) wdte.Func {
			return Combine(frame, append(args, more...)...)
		})
	}

	switch a0 := args[0].Call(frame).(type) {
	case reader:
		r := make([]io.Reader, 1, len(args))
		r[0] = a0
		for _, a := range args[1:] {
			r = append(r, a.Call(frame).(reader))
		}
		return Reader{Reader: io.MultiReader(r...)}

	case writer:
		w := make([]io.Writer, 1, len(args))
		w[0] = a0
		for _, a := range args[1:] {
			w = append(w, a.Call(frame).(writer))
		}
		return Writer{Writer: io.MultiWriter(w...)}

	default:
		panic(fmt.Errorf("Unexpected argument type: %T", a0))
	}
}

// Copy copies from a reader to a writer until it hits EOF using Go's
// io.Copy. It takes its arguments in either order, and, if given one
// only argument, returns a function which performs the copy using
// that argument and a single argument that it is given. In other
// words:
//
//     io.stdout -> io.copy io.stdin;
//
// and
//
//     io.stdin -> io.copy io.stdout;
//
// are mostly equivalent. The only difference is in the return value.
// Copy returns the second argument it was given to allow for easier
// chaining. For example, in the first example above it returns
// io.stdout, while in the second it returns io.stdin.
func Copy(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("copy")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Copy)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, more ...wdte.Func) wdte.Func {
			return Copy(frame, append(args, more...)...)
		})
	}

	var w writer
	var r reader
	var a1 wdte.Func
	switch a0 := args[0].Call(frame).(type) {
	case writer:
		w = a0
		r = args[1].Call(frame).(reader)
		a1 = r

	case reader:
		w = args[1].Call(frame).(writer)
		r = a0
		a1 = w
	}

	_, err := io.Copy(w, r)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return a1
}

type stringReader struct {
	*strings.Reader
}

func (r stringReader) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return r
}

// ReadString returns a reader that reads from a string.
func ReadString(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("readString")

	if len(args) == 0 {
		return wdte.GoFunc(ReadString)
	}

	s := args[0].Call(frame).(wdte.String)
	return stringReader{Reader: strings.NewReader(string(s))}
}

// String reads the entirety of a reader into a string and returns it.
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

// scanner is a simple wrapper that allows a bufio.Scanner to be used
// as a stream.Stream.
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

// Lines returns a stream that yields, as strings, successive lines
// read from a reader.
func Lines(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("lines")

	if len(args) == 0 {
		return wdte.GoFunc(Lines)
	}

	r := args[0].Call(frame).(reader)
	return scanner{s: bufio.NewScanner(r)}
}

// Words returns a stream that yields, as strings, successive words
// read from a reader.
func Words(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("words")

	if len(args) == 0 {
		return wdte.GoFunc(Words)
	}

	r := args[0].Call(frame).(reader)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	return scanner{s: s}
}

func write(f func(io.Writer, interface{}) error) wdte.Func {
	var gf wdte.GoFunc
	gf = func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		switch len(args) {
		case 0:
			return gf
		case 1:
			return wdte.GoFunc(func(frame wdte.Frame, more ...wdte.Func) wdte.Func {
				return gf(frame, append(more, args...)...)
			})
		}

		var w writer
		var d wdte.Func
		switch a0 := args[0].Call(frame).(type) {
		case writer:
			w = a0
			d = args[1].Call(frame)
		case wdte.Func:
			d = a0
			w = args[1].Call(frame).(writer)
		}

		err := f(w, d)
		if err != nil {
			return wdte.Error{Err: err, Frame: frame}
		}
		return w
	}
	return gf
}

// Write writes to a writer. It takes two arguments, one of which is
// the writer and one of which is the data. It is essentially
// equivalent to fmt.Fprint. It accepts the arguments in either order
// and, if given only one argument, returns a function that takes the
// other. In other words,
//
//     'Example' -> io.write io.stdout;
//
// and
//
//     io.stdout -> io.write 'Example';
//
// are equivalent.
//
// It returns the writer that it was passed.
func Write(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("write")
	return write(func(w io.Writer, v interface{}) error {
		_, err := fmt.Fprint(w, v)
		return err
	}).Call(frame, args...)
}

// Writeln is exactly like Write, but also writes a newline
// afterwards. It is essentially equivalent to fmt.Fprintln.
func Writeln(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("writeln")
	return write(func(w io.Writer, v interface{}) error {
		_, err := fmt.Fprintln(w, v)
		return err
	}).Call(frame, args...)
}

// Module returns a module for easy importing into an actual script.
// The imported functions have the same names as the functions in this
// package, except that the first letter is lowercase.
//
// In addition, it contains the following functions:
//
// * stdin, stdout, and stderr: Return readers or writers, as
//   appropriate, that wrap the standard I/O streams.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"stdin":  Reader{Reader: os.Stdin},
			"stdout": Writer{Writer: os.Stdout},
			"stderr": Writer{Writer: os.Stderr},

			"seek":  wdte.GoFunc(Seek),
			"close": wdte.GoFunc(Close),

			"combine": wdte.GoFunc(Combine),
			"copy":    wdte.GoFunc(Copy),

			"readString": wdte.GoFunc(ReadString),
			"string":     wdte.GoFunc(String),
			"lines":      wdte.GoFunc(Lines),
			"words":      wdte.GoFunc(Words),

			"write":   wdte.GoFunc(Write),
			"writeln": wdte.GoFunc(Writeln),
		},
	}
}