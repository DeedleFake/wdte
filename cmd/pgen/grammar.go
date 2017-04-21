package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/DeedleFake/wdte/scanner"
)

type Grammar struct {
	NTerms map[string]NTerm
}

func LoadGrammar(r io.Reader) (g Grammar, err error) {
	g = Grammar{
		NTerms: make(map[string]NTerm),
	}

	var cur NTerm
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()

		parts := strings.SplitN(line, "->", 2)
		if len(parts) < 2 {
			parts = strings.SplitN(line, "|", 2)
			parts[0] = cur.Name
		}
		name := strings.TrimSpace(parts[0])
		if nt, ok := g.NTerms[name]; ok {
			cur = nt
		}
		cur.Name = name
	}
	if s.Err() != nil {
		return g, s.Err()
	}

	panic("Not implemented.")
}

type Term scanner.TokenType

type NTerm struct {
	Name   string
	First  map[scanner.TokenType]struct{}
	Follow map[scanner.TokenType]struct{}
}

type Epsilon struct{}
