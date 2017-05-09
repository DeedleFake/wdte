package ast

import (
	"fmt"
	"io"

	"github.com/DeedleFake/wdte/ast/internal/pgen"
	"github.com/DeedleFake/wdte/scanner"
)

func Parse(r io.Reader) (ast Node, err error) {
	ast = &NTerm{NTerm: "script"}
	s := scanner.New(r)
	g := tokenStack{"EOF", pgen.NTerm("script")}

	more := s.Scan()
	cur := ast.(*NTerm)
	for len(g) > 0 {
		gtok := g.Pop()
		if gtok == nil {
			cur = cur.Parent().(*NTerm)
			continue
		}
		if gtok == "EOF" {
			if more {
				return nil, parseError(s, fmt.Errorf("EOF expected, but found %v", s.Tok().Type))
			}
			break
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
				return nil, parseError(s, fmt.Errorf("Expected %v (<%v>), but found %v", gtok, cur.NTerm, s.Tok().Val))
			}

			val := s.Tok().Val
			cur.AddChild(&Term{
				Term: gtok,
				p:    cur,
				v:    val,
			})
			more = s.Scan()

		case pgen.NTerm:
			rule := pgen.Table[pgen.Lookup{Term: toPGenTerm(s.Tok()), NTerm: gtok}]
			if rule == nil {
				return nil, parseError(s, fmt.Errorf("No rule for (%v, <%v>)", toPGenTerm(s.Tok()), gtok))
			}

			g.PushRule(rule)

			child := &NTerm{
				NTerm: gtok,
				p:     cur,
			}
			cur.AddChild(child)
			cur = child

		case pgen.Epsilon:
			cur.AddChild(&Epsilon{
				p: cur,
			})
		}
	}

	return ast, nil
}

func tokensEqual(stok scanner.Token, gtok pgen.Token) bool {
	switch gtok := gtok.(type) {
	case pgen.Term:
		return (gtok.Type == stok.Type) && ((stok.Type != scanner.Keyword) || (gtok.Keyword == stok.Val))
	}

	panic(fmt.Errorf("Tried to compare non-terminal: %#v", gtok))
}

func toPGenTerm(tok scanner.Token) pgen.Term {
	var keyword string
	if tok.Type == scanner.Keyword {
		keyword = fmt.Sprintf("%v", tok.Val)
	}

	return pgen.Term{
		Type:    tok.Type,
		Keyword: keyword,
	}
}

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
