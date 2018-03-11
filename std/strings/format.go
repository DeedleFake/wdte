package strings

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/DeedleFake/wdte"
)

// Format formats a string. The first argument is a format
// specification, and all trailing arguments are values to be
// formatted into the string. It does *not* use Go's C-style
// formatting specification. Instead, format specifications are {}
// with optional flags placed between them. Flags may be any
// combination of the following:
//
//    #<num> The zero-based index of the argument to be inserted.
//           Subsequent formats will increment from here. In other
//           words, '{2} {}' will yield the third and fourth
//           arguments.
//    q      Place the value in quotes using strconv.Quote.
//    ?      Mark the value with it's underlying Go type, such as
//           wdte.Number(3).
//
// Any unknown flags will cause an error to be returned instead.
//
// If Format is only given one argument, it returns a function that formats its own arguments using that original argument as the format specification.
//
// TODO: Add more flags.
func Format(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("format")

	switch len(args) {
	case 0:
		return wdte.GoFunc(Format)
	case 1:
		return wdte.GoFunc(func(frame wdte.Frame, next ...wdte.Func) wdte.Func {
			return Format(frame, append(args, next...)...)
		})
	}

	var i int
	var out bytes.Buffer
	s := strings.NewReader(string(args[0].Call(frame).(wdte.String)))
	for {
		r, _, err := s.ReadRune()
		if err != nil {
			return wdte.String(out.String())
		}

		if r == '\\' {
			r, _, err := s.ReadRune()
			if err != nil {
				out.WriteRune('\\')
				return wdte.String(out.String())
			}
			out.WriteRune(r)
			continue
		}

		if r != '{' {
			out.WriteRune(r)
			continue
		}

		var flags formatFlags
		for {
			r, _, err := s.ReadRune()
			if err != nil {
				return wdte.Error{Frame: frame, Err: fmt.Errorf("Unterminated format specifier")}
			}
			if r == '}' {
				break
			}

			switch r {
			case '#':
				i = readIndex(s)
				if i < 0 {
					return wdte.Error{Frame: frame, Err: fmt.Errorf("Invalid index in format specifier")}
				}
			case 'q':
				flags |= formatQuote
			case '?':
				flags |= formatType
			default:
				return wdte.Error{Frame: frame, Err: fmt.Errorf("Unknown format flag: %q", r)}
			}
		}

		i++
		out.WriteString(flags.Format(args[i].Call(frame)))
	}
}

func readIndex(rr io.RuneScanner) (i int) {
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return -1
		}

		if !unicode.IsNumber(r) {
			rr.UnreadRune()
			return
		}

		i = (i * 10) + int(r-'0')
	}
}

type formatFlags uint

const (
	formatQuote formatFlags = 1 + iota
	formatType
)

func (ff formatFlags) Format(val interface{}) string {
	out := fmt.Sprint(val)

	if ff&formatQuote != 0 {
		out = strconv.Quote(out)
	}

	if ff&formatType != 0 {
		out = reflect.TypeOf(val).String() + "(" + out + ")"
	}

	return out
}
