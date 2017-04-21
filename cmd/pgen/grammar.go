package main

import (
	"bufio"
	"io"
	"math"
	"strings"

	"github.com/DeedleFake/wdte/scanner"
)

type Grammar struct {
	First  map[NTerm]TermSet
	Follow map[NTerm]TermSet
}

func LoadGrammar(r io.Reader) (g Grammar, err error) {
	var cur NTerm
	rules := make(map[NTerm][][]interface{})

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if (len(line) == 0) || (line[0] == '#') {
			continue
		}

		parts := strings.SplitN(line, "->", 2)
		if len(parts) < 2 {
			parts = strings.SplitN(line, "|", 2)
			parts[0] = string(cur)
		}
		cur = NTerm(strings.TrimSpace(parts[0]))

		parts = strings.Fields(parts[1])

		rule := make([]interface{}, 0, len(parts))
		for _, p := range parts {
			rule = append(rule, part(p))
		}
		rules[cur] = append(rules[cur], rule)
	}
	if s.Err() != nil {
		return g, s.Err()
	}

	g.First = make(map[NTerm]TermSet)
	for nt := range rules {
		g.first(nt, rules)
	}

	g.Follow = make(map[NTerm]TermSet)
	for nt := range rules {
		g.follow(nt, rules)
	}

	return g, nil
}

func (g *Grammar) first(nt NTerm, rules map[NTerm][][]interface{}) TermSet {
	if s, ok := g.First[nt]; ok {
		return s
	}
	g.First[nt] = make(TermSet)

	for _, rule := range rules[nt] {
	loop:
		for _, p := range rule {
			switch p := p.(type) {
			case Term:
				g.First[nt].Add(p)

			case NTerm:
				g.First[nt].AddAll(g.first(p, rules))

				if !g.nullable(p, rules) {
					break loop
				}

			case Epsilon:
				g.First[nt].Add(Term{Type: EpsilonToken})
			}
		}
	}

	return g.First[nt]
}

func (g *Grammar) follow(nt NTerm, rules map[NTerm][][]interface{}) TermSet {
	g.Follow = make(map[NTerm]TermSet)

	panic("Not implemented.")
}

func (g *Grammar) nullable(nt NTerm, rules map[NTerm][][]interface{}) bool {
	return g.first(nt, rules).Contains(Term{Type: EpsilonToken})
}

func part(str string) interface{} {
	if (str[0] == '<') && (str[len(str)-1] == '>') {
		return NTerm(str)
	}

	if str == "ε" {
		return Epsilon{}
	}

	switch str {
	case "number":
		return Term{Type: scanner.Number}
	case "string":
		return Term{Type: scanner.String}
	case "id":
		return Term{Type: scanner.ID}
	}

	return Term{
		Type:    scanner.Keyword,
		Keyword: str,
	}
}

const EpsilonToken scanner.TokenType = scanner.TokenType(math.MaxUint32)

type Term struct {
	Type    scanner.TokenType
	Keyword string
}

type NTerm string

type Epsilon struct{}

type TermSet map[Term]struct{}

func (s TermSet) Add(t Term) {
	s[t] = struct{}{}
}

func (s TermSet) AddAll(o TermSet) {
	for t := range o {
		s.Add(t)
	}
}

func (s TermSet) Contains(t Term) bool {
	_, ok := s[t]
	return ok
}
