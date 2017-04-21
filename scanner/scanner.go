package scanner

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode"
)

type Scanner struct {
	r    io.RuneReader
	rbuf []rune

	tok Token
	err error

	tbuf  bytes.Buffer
	quote rune
}

func New(r io.Reader) *Scanner {
	var rr io.RuneReader
	switch r := r.(type) {
	case io.RuneReader:
		rr = r
	default:
		rr = bufio.NewReader(r)
	}

	return &Scanner{
		r: rr,
	}
}

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
	}

	return true
}

func (s *Scanner) Tok() Token {
	return s.tok
}

func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}

	return s.err
}

func (s *Scanner) read() (rune, error) {
	if len(s.rbuf) > 0 {
		r := s.rbuf[len(s.rbuf)-1]
		s.rbuf = s.rbuf[:len(s.rbuf)-1]
		return r, nil
	}

	r, _, err := s.r.ReadRune()
	return r, err
}

func (s *Scanner) unread(r rune) {
	s.rbuf = append(s.rbuf, r)
}

type stateFunc func(rune) stateFunc

func (s *Scanner) whitespace(r rune) stateFunc {
	if unicode.IsSpace(r) {
		return s.whitespace
	}

	if r == '-' {
		s.unread(r)
		return s.negative
	}

	if unicode.IsDigit(r) {
		s.unread(r)
		return s.number
	}

	if isQuote(r) {
		s.quote = r
		return s.string
	}

	s.unread(r)
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
	s.tok = Number{
		Val: val,
	}

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

	s.tok = String{
		Val: s.tbuf.String(),
	}

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

		val := s.tbuf.String()
		if k := getKeywordSuffix(val); k != "" {
			// BUG: This only works so long as the set of keywords doesn't
			// contain any which are contain other keywords as prefixes.
			if len(val) == len(k) {
				s.tok = Keyword{
					Val: val,
				}
				return nil
			}

			for i := len(k) - 1; i >= 0; i-- {
				s.unread(rune(k[i]))
			}

			s.tok = ID{
				Val: val[:len(val)-len(k)],
			}
			return nil
		}

		return s.id
	}

	switch val := s.tbuf.String(); isKeyword(val) {
	case true:
		s.tok = Keyword{
			Val: val,
		}
		if !isKeyword(string(r)) {
			s.unread(r)
		}

	case false:
		s.tok = ID{
			Val: val,
		}
		s.unread(r)
	}

	return nil
}
