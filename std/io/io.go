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
	"github.com/DeedleFake/wdte/std"
	"github.com/DeedleFake/wdte/wdteutil"
)

// These variables are what are returned by the corresponding
// functions in this package. If a client wants to globally redirect
// input or output, they may simply change these variables.
var (
	Stdin  io.Reader = os.Stdin
	Stdout io.Writer = os.Stdout
	Stderr io.Writer = os.Stderr
)

type stdin struct{}

func (r stdin) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return r
}

func (stdin) Read(buf []byte) (int, error) {
	return Stdin.Read(buf)
}

func (stdin) String() string {
	return "<reader(stdin)>"
}

func (stdin) Reflect(name string) bool {
	return name == "Reader"
}

type stdout struct{}

func (w stdout) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return w
}

func (stdout) Write(data []byte) (int, error) {
	return Stdout.Write(data)
}

func (stdout) String() string {
	return "<writer(stdout)>"
}

func (stdout) Reflect(name string) bool {
	return name == "Writer"
}

type stderr struct{}

func (w stderr) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return w
}

func (stderr) Write(data []byte) (int, error) {
	return Stderr.Write(data)
}

func (stderr) String() string {
	return "<writer(stderr)>"
}

func (stderr) Reflect(name string) bool {
	return name == "Writer"
}

type reader interface {
	wdte.Func
	io.Reader
}

// Reader wraps an io.Reader, allowing it to be used as a WDTE
// function.
//
// Note that using this specific type is not necessary. Any wdte.Func
// that implements io.Reader is also accepted by the functions in this
// module.
type Reader struct {
	io.Reader
}

func (r Reader) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return r
}

func (r Reader) String() string {
	if inner, ok := r.Reader.(fmt.Stringer); ok {
		return inner.String()
	}

	return "<reader>"
}

func (r Reader) Reflect(name string) bool {
	return name == "Reader"
}

type writer interface {
	wdte.Func
	io.Writer
}

// Writer wraps an io.Writer, allowing it to be used as a WDTE
// function.
//
// Note that using this specific type is not necessary. Any wdte.Func
// that implements io.Writer is also accepted by the functions in this
// module.
type Writer struct {
	io.Writer
}

func (w Writer) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return w
}

func (w Writer) String() string {
	if inner, ok := w.Writer.(fmt.Stringer); ok {
		return inner.String()
	}

	return "<writer>"
}

func (w Writer) Reflect(name string) bool {
	return name == "Writer"
}

// Seek is a WDTE function with the following signatures:
//
//    seek s n w
//    (seek w) s n
//    (seek n w) s
//
// Returns s after seeking s to n, with a relative position denoted by
// w:
//
// If w is greater than 0, it seeks relative to the beginning of s.
//
// If w is equal to 0, it seeks relative to the current location in s.
//
// If w is less than 0, it seeks relative to the end of s.
func Seek(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("seek")

	if len(args) < 3 {
		return wdteutil.SaveArgsReverse(wdte.GoFunc(Seek), args...)
	}

	s := args[0].(io.Seeker)
	off := int64(args[1].(wdte.Number))
	rel := args[2].(wdte.Number)

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

// Close is a WDTE function with the following signatures:
//
//    close c
//
// Returns c after closing it.
func Close(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("close")

	if len(args) == 0 {
		return wdte.GoFunc(Close)
	}

	c := args[0].(io.Closer)
	err := c.Close()
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return c.(wdte.Func)
}

// Combine is a WDTE function with the following signatures:
//
//    combine a ...
//    (combine a) ...
//
// If the arguments passed are readers, it returns a reader that reads
// each until EOF before continuing to the next, and finally yielding
// EOF itself when the last reader does.
//
// If the arguments passed are writers, it returns a writer that
// writes each write to all of them in turn, only returning when they
// have all returned.
func Combine(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("combine")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Combine), args...)
	}

	switch a0 := args[0].(type) {
	case reader:
		r := make([]io.Reader, 1, len(args))
		r[0] = a0
		for _, a := range args[1:] {
			r = append(r, a.(reader))
		}
		return Reader{Reader: io.MultiReader(r...)}

	case writer:
		w := make([]io.Writer, 1, len(args))
		w[0] = a0
		for _, a := range args[1:] {
			w = append(w, a.(writer))
		}
		return Writer{Writer: io.MultiWriter(w...)}

	default:
		panic(fmt.Errorf("Unexpected argument type: %T", a0))
	}
}

// Copy is a WDTE function with the following signatures:
//
//    copy w r
//    (copy w) r
//    copy r w
//    (copy r) w
//
// Copies from the reader r into the writer w until r yields EOF.
// Returns whichever argument was given second.
//
// The reason for this return discrepency is to allow both variants of
// the function to be used more easily in chains. For example, both of
// the following work:
//
//    stdout -> copy stdin -> ... # Later elements will be given stdout.
//    stdin -> copy stdout -> ... # Later elements will be given stdin.
func Copy(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("copy")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Copy), args...)
	}

	var w writer
	var r reader
	var a1 wdte.Func
	switch a0 := args[0].(type) {
	case writer:
		w = a0
		r = args[1].(reader)
		a1 = r

	case reader:
		w = args[1].(writer)
		r = a0
		a1 = w
	}

	_, err := io.Copy(w, r)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return a1
}

// String is a WDTE function with the following signature:
//
//    string r
//
// Reads the entirety of the reader r and returns the result as a
// string.
func String(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("string")

	if len(args) == 0 {
		return wdte.GoFunc(String)
	}

	r := args[0].(reader)

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

func (s scanner) Reflect(name string) bool {
	return name == "Stream"
}

// Lines is a WDTE function with the following signature:
//
//    lines r
//
// Returns a stream.Stream that yields, as strings, successive lines
// read from the reader r.
func Lines(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("lines")

	if len(args) == 0 {
		return wdte.GoFunc(Lines)
	}

	r := args[0].(reader)
	return scanner{s: bufio.NewScanner(r)}
}

// Words is a WDTE function with the following signature:
//
//    words r
//
// Returns a stream.Stream that yields, as strings, successive words
// read from the reader r.
func Words(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("words")

	if len(args) == 0 {
		return wdte.GoFunc(Words)
	}

	r := args[0].(reader)
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	return scanner{s: s}
}

// Scan is a WDTE function with the following signatures:
//
//    scan r sep
//    (scan r) sep
//    scan sep r
//    (scan sep) r
//
// Returns a stream.Stream that yields sections of the reader r split
// around the separator string sep. For example,
//
//    readString 'this--is--an--example' -> scan '--'
//
// will return a stream.Stream that will yield 'this', 'is', 'an', and
// 'example'.
func Scan(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("scan")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Scan), args...)
	}

	var r reader
	var sep wdte.String
	switch a0 := args[0].(type) {
	case reader:
		r = a0
		sep = args[1].(wdte.String)
	case wdte.String:
		r = args[1].(reader)
		sep = a0
	}

	seb := []byte(sep)

	s := bufio.NewScanner(r)
	s.Split(func(data []byte, eof bool) (int, []byte, error) {
		start := bytes.Index(data, seb)
		if start < 0 {
			if eof {
				if len(data) == 0 {
					return 1, nil, nil
				}

				return len(data), data, nil
			}
			return 0, nil, nil
		}

		return start + len(seb), data[:start], nil
	})

	return scanner{s: s}
}

type runeStream struct {
	r io.RuneReader
}

func (r runeStream) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	return r
}

func (r runeStream) Next(frame wdte.Frame) (wdte.Func, bool) {
	c, _, err := r.r.ReadRune()
	if err != nil {
		if err == io.EOF {
			return nil, false
		}
		return wdte.Error{Frame: frame, Err: err}, true
	}
	return wdte.Number(c), true
}

func (r runeStream) Reflect(name string) bool {
	return name == "Stream"
}

// Runes is a WDTE function with the following signature:
//
//    runes r
//
// Returns a stream.Stream that yields individual Unicode characters
// from the reader r as numbers.
//
// TODO: Maybe it makes more sense for them to be yielded as strings
// with a length of one.
func Runes(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("runes")

	if len(args) == 0 {
		return wdte.GoFunc(Runes)
	}

	a := args[0]

	var r io.RuneReader
	switch a := a.(type) {
	case io.RuneReader:
		r = a
	case io.Reader:
		r = bufio.NewReader(a)
	default:
		panic(fmt.Errorf("Unexpected argument type: %T", a))
	}

	return runeStream{r: r}
}

func write(f func(io.Writer, interface{}) error) (gf wdte.Func) {
	return wdte.GoFunc(func(frame wdte.Frame, args ...wdte.Func) wdte.Func {
		if len(args) < 2 {
			return wdteutil.SaveArgsReverse(gf, args...)
		}

		var w writer
		var d wdte.Func
		switch a0 := args[0].(type) {
		case writer:
			w = a0
			d = args[1]
		case wdte.Func:
			d = a0
			w = args[1].(writer)
		}

		err := f(w, d)
		if err != nil {
			return wdte.Error{Err: err, Frame: frame}
		}
		return w
	})
}

// Write is a WDTE function with the following signatures:
//
//    write w d
//    (write w) d
//    write d w
//    (write d) w
//
// It writes the data d to the writer w in much the same way that Go's
// fmt.Fprint does. It returns w to allow for easier chaining.
//
// If both arguments are writers, it will consider either the first
// argument or the outer argument to be w.
func Write(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("write")
	return write(func(w io.Writer, v interface{}) error {
		_, err := fmt.Fprint(w, v)
		return err
	}).Call(frame, args...)
}

// Writeln is a WDTE function with the following signatures:
//
//    writeln w d
//    (writeln w) d
//    writeln d w
//    (writeln d) w
//
// It writes the data d to the writer w in much the same way that Go's
// fmt.Fprintln does. It returns w to allow for easier chaining.
//
// If both arguments are writers, it will consider either the first
// argument or the outer argument to be w.
func Writeln(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("writeln")
	return write(func(w io.Writer, v interface{}) error {
		_, err := fmt.Fprintln(w, v)
		return err
	}).Call(frame, args...)
}

// Panic is a WDTE function with the following signatures:
//
//    panic err
//    panic w err
//    panic desc err
//    panic w desc err
//
// Note that, somewhat unusually, Panic accepts its arguments in any order.
//
// It writes the given error to w, prepending the optional
// description in the form `desc: err` and appending a newline. It
// then returns the error. If an error occurs somewhere internally,
// such as while printing, that error is returned instead.
//
// If w is not given, it defaults to Stderr.
//
// Panic is primarily intended for use with the error chain operator.
// For example:
//
//    + a b -| panic 'Failed to add a and b';
func Panic(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("panic")

	w := writer(stderr{})
	var desc wdte.String
	var e error

	set := func(f wdte.Func) {
		switch f := f.(type) {
		case writer:
			w = f
		case wdte.String:
			desc = f + ": "
		case error:
			e = f

		default:
			panic(fmt.Errorf("Unexpected argument type: %T", f))
		}
	}

	n := 3
	if len(args) < 3 {
		n = len(args)
	}

	for i, arg := range args[:n] {
		args[i] = arg
		set(args[i])
	}
	if e == nil {
		return wdteutil.SaveArgs(wdte.GoFunc(Panic), args...)
	}

	_, err := fmt.Fprintf(w, "%v%v\n", desc, e)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return e.(wdte.Func)
}

// Scope is a scope containing the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"stdin":  stdin{},
	"stdout": stdout{},
	"stderr": stderr{},

	"seek":  wdte.GoFunc(Seek),
	"close": wdte.GoFunc(Close),

	"combine": wdte.GoFunc(Combine),
	"copy":    wdte.GoFunc(Copy),

	"string": wdte.GoFunc(String),
	"lines":  wdte.GoFunc(Lines),
	"words":  wdte.GoFunc(Words),
	"scan":   wdte.GoFunc(Scan),
	"runes":  wdte.GoFunc(Runes),

	"write":   wdte.GoFunc(Write),
	"writeln": wdte.GoFunc(Writeln),
	"panic":   wdte.GoFunc(Panic),
	//"panicln": wdte.GoFunc(Panicln),
})

func init() {
	std.Register("io", Scope)
}
