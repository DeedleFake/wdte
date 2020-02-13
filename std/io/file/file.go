// Package file provides functions for dealing with files.
package file

import (
	"fmt"
	"os"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

const (
	modeReader = 1 << iota
	modeWriter
)

// File wraps an os.File, allowing it to be used as a WDTE function.
// While it contains unexproted fields, it is safe for a client to
// simply wrap an *os.File in it manually.
//
// A file is considered a "File" by reflection, as well as a "Reader"
// if it is opened for reading and a "Writer" if it is opened for
// writing. If the file was created manually, it will be considered
// both a reader and a writer.
type File struct {
	*os.File
	mode int
}

func (f File) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return f
}

func (f File) String() string { // nolint
	return fmt.Sprintf("<file %q>", f.Name())
}

func (f File) Reflect(name string) bool { // nolint
	return name == "File" ||
		(((f.mode == 0) || (f.mode&modeReader != 0)) && (name == "Reader")) ||
		(((f.mode == 0) || (f.mode&modeWriter != 0)) && (name == "Writer"))
}

// Open is a WDTE function with the following signature:
//
//    open path
//
// Opens the file at path and returns it.
func Open(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("open")

	if len(args) == 0 {
		return wdte.GoFunc(Open)
	}

	path := args[0].(wdte.String)
	file, err := os.Open(string(path))
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return File{File: file, mode: modeReader}
}

// Create is a WDTE function with the following signature:
//
//    create path
//
// Creates the file at path, truncating it if it already exists, and
// returns it.
func Create(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("create")

	if len(args) == 0 {
		return wdte.GoFunc(Create)
	}

	path := args[0].(wdte.String)
	file, err := os.Create(string(path))
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return File{File: file, mode: modeWriter}
}

// Append is a WDTE function with the following signature:
//
//    append path
//
// Opens the file at path for appending, creating it if it doesn't
// already exist, and returns it.
func Append(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("append")

	if len(args) == 0 {
		return wdte.GoFunc(Append)
	}

	path := args[0].(wdte.String)
	file, err := os.OpenFile(string(path), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return File{File: file, mode: modeWriter}
}

// Scope is a scope containing the functions in this package.
var Scope = wdte.S().Map(map[wdte.ID]wdte.Func{
	"open":   wdte.GoFunc(Open),
	"create": wdte.GoFunc(Create),
	"append": wdte.GoFunc(Append),
})

func init() {
	std.Register("io/file", Scope)
}
