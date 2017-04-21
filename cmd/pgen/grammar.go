package main

import (
	"io"

	"github.com/DeedleFake/wdte/scanner"
)

type Grammar struct {
	NTerms map[string]NTerm
}

func LoadGrammar(r io.Reader) (Grammar, error) {
	panic("Not implemented.")
}

type Term struct {
	scanner.Token
}

type NTerm struct {
	Name   string
	First  map[string]struct{}
	Follow map[string]struct{}
}

type Epsilon struct{}
