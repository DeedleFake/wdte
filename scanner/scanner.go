package scanner

import (
	"bufio"
	"bytes"
	"io"
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

func (s *Scanner) number(r rune) stateFunc {
	panic("Not implemented.")
}

func (s *Scanner) string(r rune) stateFunc {
	panic("Not implemented.")
}

func (s *Scanner) id(r rune) stateFunc {
	panic("Not implemented.")
}
