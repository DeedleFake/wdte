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
	"github.com/DeedleFake/wdte/wdteutil"
)

// Format is a WDTE function with the following signatures:
//
//    format tmpl ...
//
// Format has some special rules for returning a partial function. For
// more information, see below.
//
// This is the general-purpose string formatting function of the
// standard library, similar to Go's fmt.Sprintf(). Unlike
// fmt.Sprintf(), however, format uses a custom formatting
// specification. A format in the string tmpl is of the form {} with
// optional flags placed between them. Flags may be any combination of
// the following:
//
//    #<num> The zero-based index of the argument to be inserted.
//           Subsequent formats will increment from here. In other
//           words, '{2} {}' will yield the third and fourth
//           arguments.
//    q      Place the value in quotes using strconv.Quote.
//    ?      Mark the value with it's underlying Go type, such as
//           wdte.Number(3).
//
// Format's rules for returning a partial function are dependant on
// the value of the first argument. Specifically, if the first
// argument attempts to substitute in more arguments than were given,
// a partial function will be returned. For example,
//
//    format '' # Returns ''
//    format '{}' 3 # Returns '3'
//    format '{}' # Returns a partial function.
//    (format '{} {}' 3) 'example' # Returns '3 example'
//
// Note that the total number of arguments required is the smallest
// number necessary to perform every substitution specified by the
// first argument. For example,
//
//    format '{3} {}'
//
// will return a partial function that requires 5 arguments before it
// will return the formatted string.
//
// TODO: Add more flags.
func Format(frame wdte.Frame, args ...wdte.Func) wdte.Func {
	frame = frame.Sub("format")

	if len(args) < 2 {
		return wdteutil.SaveArgs(wdte.GoFunc(Format), args...)
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
		if i >= len(args) {
			// TODO: Cache the current buffer.
			return wdteutil.SaveArgs(wdte.GoFunc(Format), args...)
		}

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
