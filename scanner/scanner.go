package scanner

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode"
)

// A Scanner tokenizes runes from an io.Reader.
type Scanner struct {
	r               io.RuneReader
	rbuf            []rune
	line, col, pcol int

	tok         Token
	tline, tcol int
	err         error

	tbuf  bytes.Buffer
	quote rune
}

// New returns a new Scanner that reads from r.
func New(r io.Reader) *Scanner {
	var rr io.RuneReader
	switch r := r.(type) {
	case io.RuneReader:
		rr = r
	default:
		rr = bufio.NewReader(r)
	}

	return &Scanner{
		r:    rr,
		line: 1,
	}
}

// Scan reads the next token from the underlying io.Reader. If a token
// was successfully read, it returns true. It is designed to be used
// in a loop, similarly to bufio.Scanner's API.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	s.tbuf.Reset()

	state := s.whitespace
	for (state != nil) && (s.err == nil) {
		r, err := s.read()
		if err != nil {
			s.err = err

			if err == io.EOF {
				r = '\n'
			}
		}

		state = state(r)

		if (s.err == io.EOF) && (state != nil) {
			return false
		}
	}

	return true
}

// Tok returns the latest token scanned. If there was an error or a
// token hasn't been scanned yet, its return value is undefined.
func (s *Scanner) Tok() Token {
	return s.tok
}

// Err returns the error that stopped the scanner, if any.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}

	return s.err
}

func (s *Scanner) read() (r rune, err error) {
	defer func() {
		s.col++

		if r == '\n' {
			s.line++
			s.pcol = s.col
			s.col = 0
		}
	}()

	if len(s.rbuf) > 0 {
		r = s.rbuf[len(s.rbuf)-1]
		s.rbuf = s.rbuf[:len(s.rbuf)-1]
		return
	}

	r, _, err = s.r.ReadRune()
	return
}

func (s *Scanner) unread(r rune) {
	s.rbuf = append(s.rbuf, r)

	s.col--
	if r == '\n' {
		s.line--
		s.col = s.pcol
	}
}

func (s *Scanner) setTok(t TokenType, v interface{}) {
	s.tok = Token{
		Line: s.tline,
		Col:  s.tcol,
		Type: t,
		Val:  v,
	}
}

type stateFunc func(rune) stateFunc

func (s *Scanner) whitespace(r rune) stateFunc {
	if unicode.IsSpace(r) {
		return s.whitespace
	}

	if r == '-' {
		s.tline, s.tcol = s.line, s.col
		s.unread(r)
		return s.negative
	}

	if unicode.IsDigit(r) {
		s.tline, s.tcol = s.line, s.col
		s.unread(r)
		return s.number
	}

	if isQuote(r) {
		s.tline, s.tcol = s.line, s.col
		s.quote = r
		return s.string
	}

	s.unread(r)
	s.tline, s.tcol = s.line, s.col
	return s.id
}

func (s *Scanner) negative(r rune) stateFunc {
	if r == '-' {
		s.tbuf.WriteRune(r)
		return s.negative
	}

	if unicode.IsDigit(r) {
		s.tbuf.WriteRune(r)
		return s.number
	}

	s.unread(r)
	return s.id
}

func (s *Scanner) number(r rune) stateFunc {
	if unicode.IsDigit(r) || (r == '.') {
		s.tbuf.WriteRune(r)
		return s.number
	}

	val, _ := strconv.ParseFloat(s.tbuf.String(), 64)
	s.setTok(Number, val)

	s.unread(r)
	return nil
}

func (s *Scanner) string(r rune) stateFunc {
	if r == '\\' {
		return s.escape
	}

	if r != s.quote {
		s.tbuf.WriteRune(r)
		return s.string
	}

	s.setTok(String, s.tbuf.String())

	return nil
}

func (s *Scanner) escape(r rune) stateFunc {
	switch r {
	case 'n':
		s.tbuf.WriteRune('\n')
	case 't':
		s.tbuf.WriteRune('\t')
	case '\n':
	default:
		s.tbuf.WriteRune(r)
	}

	return s.string
}

func (s *Scanner) id(r rune) stateFunc {
	if !unicode.IsSpace(r) {
		s.tbuf.WriteRune(r)

		// TODO: Find a way to do this without allocating and copying.
		val := s.tbuf.String()
		if k := getKeywordSuffix(val); k != "" {
			// BUG: This only works so long as the set of keywords doesn't
			// contain any which are contain other keywords as prefixes.
			if len(val) == len(k) {
				s.setTok(Keyword, val)
				return nil
			}

			for i := len(k) - 1; i >= 0; i-- {
				s.unread(rune(k[i]))
			}

			s.setTok(ID, val[:len(val)-len(k)])
			return nil
		}

		return s.id
	}

	s.setTok(ID, s.tbuf.String())
	s.unread(r)
	return nil
}
