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
			if _, ok := s.Tok().Val.(scanner.EOF); !ok {
				return nil, parseError(s, fmt.Errorf("EOF expected, but found %T", s.Tok().Val))
			}
			return cur, nil
		}
	}
}

func tokensEqual(stok scanner.Token, gtok pgen.Token) bool {
	switch gtok := gtok.(type) {
	case pgen.Term:
		st, ok := stok.Val.(scanner.Keyword)
		return (ok && (st == scanner.Keyword(gtok.Keyword))) || (toPGenTerm(stok) == gtok)
	}

	panic(fmt.Errorf("Tried to compare non-terminal: %#v", gtok))
}

func toPGenTerm(tok scanner.Token) pgen.Token {
	switch tok := tok.Val.(type) {
	case scanner.ID:
		return pgen.Term{Type: pgen.ID}
	case scanner.String:
		return pgen.Term{Type: pgen.String}
	case scanner.Number:
		return pgen.Term{Type: pgen.Number}
	case scanner.Keyword:
		return pgen.Term{Type: pgen.Keyword, Keyword: string(tok)}
	case scanner.EOF:
		return pgen.EOF{}
	}

	panic(fmt.Errorf("Unexpected token type: %T", tok.Val))
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
