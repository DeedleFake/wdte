package ast

import (
	"github.com/DeedleFake/wdte/ast/internal/pgen"
	"github.com/DeedleFake/wdte/scanner"
)

// A Node represents a node of the AST.
type Node interface {
	// Parent returns the Node's parent, or nil if it is the root node.
	Parent() Node

	// Children returns the node's children in left-to-right order.
	Children() []Node
}

// A Term is a Node that represents a terminal, such as a string or a
// keyword. Terms are always leaf nodes.
type Term struct {
	tok scanner.Token

	t pgen.Term
	p Node
}

func (t Term) Parent() Node {
	return t.p
}

// Tok returns the scanner token that the node was generated from.
func (t Term) Tok() scanner.Token {
	return t.tok
}

func (t Term) Children() []Node {
	return nil
}

// An NTerm is a Node that represents a non-terminal. NTerms are
// always parent nodes.
type NTerm struct {
	nt pgen.NTerm
	p  Node
	c  []Node
}

// Name returns the name of the non-terminal, minus the surrounding `<` and `>`.
func (nt NTerm) Name() string {
	return string(nt.nt)
}

func (nt NTerm) Parent() Node {
	return nt.p
}

// AddChild adds a child to the right-hand side of the NTerm's list of
// children.
func (nt *NTerm) AddChild(n Node) {
	if nt == nil {
		return
	}

	nt.c = append(nt.c, n)
}

func (nt NTerm) Children() []Node {
	return nt.c
}

// An Epsilon is a special terminal which represnts a non-action.
type Epsilon struct {
	p Node
}

func (e Epsilon) Parent() Node {
	return e.p
}

func (e Epsilon) Children() []Node {
	return nil
}
