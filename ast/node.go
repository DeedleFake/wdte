package ast

import (
	"github.com/DeedleFake/wdte/ast/internal/pgen"
	"github.com/DeedleFake/wdte/scanner"
)

type Node interface {
	Parent() Node
	Children() []Node
}

type Term struct {
	tok scanner.Token

	t pgen.Term
	p Node
}

func (t Term) Parent() Node {
	return t.p
}

func (t Term) Tok() scanner.Token {
	return t.tok
}

func (t Term) Children() []Node {
	return nil
}

type NTerm struct {
	nt pgen.NTerm
	p  Node
	c  []Node
}

func (nt NTerm) Parent() Node {
	return nt.p
}

func (nt *NTerm) AddChild(n Node) {
	nt.c = append(nt.c, n)
}

func (nt NTerm) Children() []Node {
	return nt.c
}

type Epsilon struct {
	p Node
}

func (e Epsilon) Parent() Node {
	return e.p
}

func (e Epsilon) Children() []Node {
	return nil
}
