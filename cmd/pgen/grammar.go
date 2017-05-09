package main

import (
	"bufio"
	"io"
	"strings"

	"github.com/DeedleFake/wdte/scanner"
)

type Grammar struct {
	First  map[NTerm]TermSet
	Follow map[NTerm]TermSet
}

func LoadGrammar(r io.Reader) (g Grammar, err error) {
	var cur NTerm
	rules := make(map[NTerm][]Rule)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if (len(line) == 0) || (line[0] == '#') {
			continue
		}

		parts := strings.SplitN(line, "->", 2)
		if line[0] == '|' {
			parts = strings.SplitN(line, "|", 2)
			parts[0] = string(cur)
		}
		cur = NTerm(strings.TrimSpace(parts[0]))

		parts = strings.Fields(parts[1])

		rule := make(Rule, 0, 1+len(parts))
		rule = append(rule, cur)
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

func (g *Grammar) first(nt NTerm, rules map[NTerm][]Rule) TermSet {
	if s, ok := g.First[nt]; ok {
		return s
	}
	g.First[nt] = make(TermSet)

	for _, rule := range rules[nt] {
	loop:
		for _, p := range rule[1:] {
			switch p := p.(type) {
			case Term:
				g.First[nt].Add(p, rule)
				break loop

			case NTerm:
				g.First[nt].AddAll(g.first(p, rules), rule)

				if !g.nullable(p, rules) {
					break loop
				}
			}
		}
	}

	return g.First[nt]
}

func (g *Grammar) follow(nt NTerm, rules map[NTerm][]Rule) TermSet {
	g.Follow = make(map[NTerm]TermSet)

	panic("Not implemented.")
}

func (g *Grammar) nullable(nt NTerm, rules map[NTerm][]Rule) bool {
outer:
	for _, rule := range rules[nt] {
		for _, p := range rule[1:] {
			switch p := p.(type) {
			case Term:
				continue outer

			case NTerm:
				if p == nt {
					continue
				}

			case Epsilon:
				continue
			}

			return true
		}
	}

	return false
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

type Rule []interface{}

type Term struct {
	Type    scanner.TokenType
	Keyword string
}

type NTerm string

type Epsilon struct{}

type TermSet map[Term]Rule

func (s TermSet) Add(t Term, r Rule) {
	if _, ok := s[t]; ok {
		return
	}

	s[t] = r
}

func (s TermSet) AddAll(o TermSet, r Rule) {
	for t := range o {
		s.Add(t, r)
	}
}

func (s TermSet) Contains(t Term) bool {
	_, ok := s[t]
	return ok
}
