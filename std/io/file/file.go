// Package file provides functions for dealing with files.
package file

import (
	"os"

	"github.com/DeedleFake/wdte"
	"github.com/DeedleFake/wdte/std"
)

// File wraps an os.File, allowing it to be used as a WDTE function.
type File struct {
	*os.File
}

func (f File) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return f
}

// Open is a WDTE function with the following signature:
//
//    open path
//
// Opens the file at path and returns it.
func Open(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Open)
	}

	path := args[0].Call(frame).(wdte.String)
	file, err := os.Open(string(path))
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return File{File: file}
}

// Create is a WDTE function with the following signature:
//
//    create path
//
// Creates the file at path, truncating it if it already exists, and
// returns it.
func Create(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Create)
	}

	path := args[0].Call(frame).(wdte.String)
	file, err := os.Create(string(path))
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return File{File: file}
}

// Append is a WDTE function with the following signature:
//
//    append path
//
// Opens the file at path for appending, creating it if it doesn't
// already exist, and returns it.
func Append(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	if len(args) == 0 {
		return wdte.GoFunc(Append)
	}

	path := args[0].Call(frame).(wdte.String)
	file, err := os.OpenFile(string(path), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return wdte.Error{Err: err, Frame: frame}
	}
	return File{File: file}
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
