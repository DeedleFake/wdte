package ast

import (
	"fmt"
	"io"

	"github.com/DeedleFake/wdte/ast/internal/pgen"
	"github.com/DeedleFake/wdte/scanner"
)

// Parse parses a full script, returning the root node of the AST.
func Parse(r io.Reader, macros scanner.MacroMap) (Node, error) {
	return parse(r, tokenStack{pgen.NTerm("script")}, pgen.Table, macros)
}

func parse(r io.Reader, g tokenStack, table map[pgen.Lookup]pgen.Rule, macros scanner.MacroMap) (ast Node, err error) {
	s := scanner.New(r, macros)

	more := s.Scan()
	var cur *NTerm
	for {
		gtok := g.Pop()
		if gtok == nil {
			cur = cur.Parent().(*NTerm)
			continue
		}
		if !more {
			if err = s.Err(); err != nil {
				return nil, parseError(s, err)
			}

			return nil, parseError(s, fmt.Errorf("Expected %v, but found EOF", gtok))
		}

		switch gtok := gtok.(type) {
		case pgen.Term:
			if !tokensEqual(s.Tok(), gtok) {
				return nil, parseError(s, fmt.Errorf("Expected %v (<%v>), but found %v", gtok, cur.nt, s.Tok().Val))
			}

			cur.AddChild(&Term{
				tok: s.Tok(),

				t: gtok,
				p: cur,
			})
			more = s.Scan()

		case pgen.NTerm:
			rule := table[pgen.Lookup{Term: toPGenTerm(s.Tok()), NTerm: gtok}]
			if rule == nil {
				return nil, parseError(s, fmt.Errorf("No rule for (%v, <%v>)", toPGenTerm(s.Tok()), gtok))
			}

			g.PushRule(rule)

			child := &NTerm{
				nt: gtok,
				p:  cur,
			}
			cur.AddChild(child)
			cur = child

		case pgen.Epsilon:
			cur.AddChild(&Epsilon{
				p: cur,
			})

		case pgen.EOF:
			if s.Tok().Type != scanner.EOF {
				return nil, parseError(s, fmt.Errorf("EOF expected, but found %v", s.Tok().Type))
			}
			return cur, nil
		}
	}
}

func tokensEqual(stok scanner.Token, gtok pgen.Token) bool {
	switch gtok := gtok.(type) {
	case pgen.Term:
		return (gtok.Type == stok.Type) && ((stok.Type != scanner.Keyword) || (gtok.Keyword == stok.Val))
	}

	panic(fmt.Errorf("Tried to compare non-terminal: %#v", gtok))
}

func toPGenTerm(tok scanner.Token) pgen.Token {
	var keyword string
	switch tok.Type {
	case scanner.Keyword:
		keyword = fmt.Sprintf("%v", tok.Val)
	case scanner.EOF:
		return pgen.EOF{}
	}

	return pgen.Term{
		Type:    tok.Type,
		Keyword: keyword,
	}
}

// A ParseError is returned if an error happens during parsing.
type ParseError struct {
	Line, Col int
	Err       error
}

func parseError(s *scanner.Scanner, err error) ParseError {
	line, col := s.Pos()
	return ParseError{
		Line: line, Col: col,
		Err: err,
	}
}

func (err ParseError) Error() string {
	return fmt.Sprintf("%v:%v: %v", err.Line, err.Col, err.Err)
}
