package scanner

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"unicode"
)

type MacroMap map[string]func(string) ([]Token, error)

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
	macro string

	macroMap MacroMap
	macroBuf []Token
}

// New returns a new Scanner that reads from r.
func New(r io.Reader, macros MacroMap) *Scanner {
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

		macroMap: macros,
	}
}

// Scan reads the next token from the underlying io.Reader. If a token
// was successfully read, it returns true. It is designed to be used
// in a loop, similarly to bufio.Scanner's API.
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	if len(s.macroBuf) > 0 {
		tok := s.macroBuf[len(s.macroBuf)-1]
		s.macroBuf = s.macroBuf[:len(s.macroBuf)-1]
		s.setTok(tok.Type, tok.Val)
		return s.err == nil
	}

	s.tbuf.Reset()

	var eof bool
	state := s.whitespace
	for state != nil {
		r, err := s.read()
		switch err {
		case nil:
			eof = false
		case io.EOF:
			if eof {
				s.err = err
				s.setTok(EOF, nil)
				return true
			}

			r = '\n'
			eof = true

		default:
			s.err = err
			return false
		}

		state = state(r)
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

// Pos returns the line and column of the input that the scanner is
// currently on.
func (s *Scanner) Pos() (line, col int) {
	return s.line, s.col
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
	switch t {
	case Keyword:
		switch v {
		case ")", "]", "}":
			if (s.tok.Type != Keyword) || (s.tok.Val != ";") {
				s.unread(rune(v.(string)[0]))
				v = ";"
			}
		}

	case Macro:
		v := v.([2]string)

		macro := s.macroMap[v[0]]
		if macro == nil {
			break
		}

		toks, err := macro(v[1])
		if len(toks) > 0 {
			for i := len(toks) - 1; i >= 1; i-- {
				s.macroBuf = append(s.macroBuf, toks[i])
			}
			s.tok = toks[0]
		}
		s.err = err
		return
	}

	s.tok = Token{
		Line: s.tline,
		Col:  s.tcol,
		Type: t,
		Val:  v,
	}
}

type stateFunc func(rune) stateFunc

func (s *Scanner) whitespace(r rune) stateFunc {
	if r == '#' {
		return s.comment
	}

	if unicode.IsSpace(r) {
		return s.whitespace
	}

	if (r == '-') || (r == '.') {
		s.tline, s.tcol = s.line, s.col
		s.tbuf.WriteRune(r)
		return s.maybeNumber
	}

	if r == '@' {
		s.tline, s.tcol = s.line, s.col
		return s.macroName
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

func (s *Scanner) comment(r rune) stateFunc {
	if r == '\n' {
		return s.whitespace
	}

	return s.comment
}

func (s *Scanner) maybeNumber(r rune) stateFunc {
	if unicode.IsDigit(r) {
		s.tbuf.WriteRune(r)
		return s.number
	}

	s.unread(r)

	tbuf := []rune(s.tbuf.String())
	for i := len(tbuf) - 1; i >= 0; i-- {
		s.unread(tbuf[i])
	}
	s.tbuf.Reset()

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
	val := s.tbuf.String() + string(r)
	if k := symbolicPrefix(val); k != "" {
		if len(val) == len(k) {
			s.tbuf.WriteRune(r)
			return s.id
		}

		valr := []rune(val)
		for i := len(val) - 1; i >= len(k); i-- {
			s.unread(valr[i])
		}

		s.setTok(Keyword, k)
		return nil
	}

	if !unicode.IsSpace(r) {
		s.tbuf.WriteRune(r)

		if k := symbolicSuffix(val); k != "" {
			kr := []rune(k)
			for i := len(k) - 1; i >= 0; i-- {
				s.unread(kr[i])
			}

			t, val := ID, val[:len(val)-len(k)]
			if isKeyword(val) {
				t = Keyword
			}
			s.setTok(t, val)
			return nil
		}

		return s.id
	}

	t, val := ID, s.tbuf.String()
	if isKeyword(val) {
		t = Keyword
	}
	s.setTok(t, val)
	s.unread(r)
	return nil
}

func (s *Scanner) macroName(r rune) stateFunc {
	if unicode.IsDigit(r) || unicode.IsLetter(r) {
		s.tbuf.WriteRune(r)
		return s.macroName
	}

	s.quote = endQuote(r)
	s.macro = s.tbuf.String()
	s.tbuf.Reset()
	return s.macroInput
}

func (s *Scanner) macroInput(r rune) stateFunc {
	if r == s.quote {
		s.setTok(Macro, [2]string{s.macro, s.tbuf.String()})
		return nil
	}

	s.tbuf.WriteRune(r)
	return s.macroInput
}
