// Package file provides functions for dealing with files.
package file

import (
	"os"

	"github.com/DeedleFake/wdte"
)

// File wraps an os.File, allowing it to be used as a WDTE function.
type File struct {
	*os.File
}

func (f File) Call(frame wdte.Frame, args ...wdte.Func) wdte.Func { // nolint
	return f
}

// Open opens a file and returns it as a reader.
func Open(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("open")

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

// Create creates a file, truncating it if it exists, and returns it
// as a writer.
func Create(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("create")

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

// Append opens a file for appending as a writer. If it doesn't exist
// already, it is created.
func Append(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.WithID("append")

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

// Module returns a module containing the functions in this package.
func Module() *wdte.Module {
	return &wdte.Module{
		Funcs: map[wdte.ID]wdte.Func{
			"open":   wdte.GoFunc(Open),
			"create": wdte.GoFunc(Create),
			"append": wdte.GoFunc(Append),
		},
	}
}
