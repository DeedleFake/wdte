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
				panic("Not implemented.")
			}
			break
		}
		if !more {
			if err = s.Err(); err != nil {
				return nil, err
			}

			panic("Not implemented.")
		}

		switch gtok := gtok.(type) {
		case pgen.Term:
			if !tokensEqual(s.Tok(), gtok) {
				panic("Not implemented.")
			}

			cur.AddChild(&Term{
				Term: gtok,
				p:    cur,
			})
			more = s.Scan()

		case pgen.NTerm:
			rule := pgen.Table[pgen.Lookup{Term: toPGenTerm(s.Tok()), NTerm: gtok}]
			if rule == nil {
				panic("Not implemented.")
			}

			g.Push(rule)

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
	case Term:
		return (gtok.Type == stok.Type) && ((stok.Type != scanner.Keyword) || (gtok.Keyword == stok.Val))
	}

	panic(fmt.Errorf("Tried to compare non-terminal: %#v", gtok))
}

func toPGenTerm(tok scanner.Token) pgen.Term {
	return pgen.Term{
		Type:    tok.Type,
		Keyword: fmt.Sprintf("%v", tok.Val),
	}
}
