package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Grammar map[NTerm][]Rule

func LoadGrammar(r io.Reader) (g Grammar, err error) {
	var cur NTerm
	g = make(map[NTerm][]Rule)

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
		for _, p := range parts {
			rule = append(rule, NewToken(p))
		}
		g[cur] = append(g[cur], rule)
	}
	return g, s.Err()
}

func (g Grammar) Terms() TokenSet {
	ts := make(TokenSet)
	for _, rules := range g {
		for _, rule := range rules {
			for _, p := range rule {
				if p, ok := p.(Term); ok {
					ts.Add(p, nil)
				}
			}
		}
	}

	return ts
}

func (g Grammar) Nullable(tok Token) bool {
	switch tok := tok.(type) {
	case Term:
		return false

	case NTerm:
	rules:
		for _, rule := range g[tok] {
			if rule.Epsilon() {
				return true
			}

			for _, p := range rule {
				if !g.Nullable(p) {
					continue rules
				}
			}

			return true
		}

		return false

	case Epsilon:
		return true
	}

	panic(fmt.Errorf("Unexpected token type: %T", tok))
}

func (g Grammar) First(tok Token) TokenSet {
	ts := make(TokenSet)

	switch tok := tok.(type) {
	case Term, Epsilon:
		ts.Add(tok, nil)

	case NTerm:
		for _, rule := range g[tok] {
			for _, p := range rule {
				ts.AddAll(g.First(p), rule)
				if !g.Nullable(p) {
					break
				}
			}
		}
	}

	return ts
}

func (g Grammar) Follow(nt NTerm) TokenSet {
	return g.followWithout(nt, nil)
}

func (g Grammar) followWithout(nt NTerm, ignore []NTerm) TokenSet {
	isIgnored := func(nt NTerm) bool {
		for _, i := range ignore {
			if i == nt {
				return true
			}
		}

		return false
	}

	ts := make(TokenSet)

	for name, rules := range g {
		for _, rule := range rules {
			for i, tok := range rule {
				if tok == nt {
					if i == len(rule)-1 {
						if !isIgnored(name) {
							ts.AddAll(g.followWithout(name, append(ignore, nt)), rule)
						}
						continue
					}

					for i := i + 1; i < len(rule); i++ {
						ts.AddAll(g.First(rule[i]), rule)
						if !g.Nullable(rule[i]) {
							break
						}

						if i == len(rule)-1 {
							if !isIgnored(name) {
								ts.AddAll(g.followWithout(name, append(ignore, nt)), rule)
							}
							continue
						}
					}
				}
			}
		}
	}

	ts.Remove(Epsilon{})
	return ts
}

type Rule []Token

func (r Rule) Epsilon() bool {
	return (len(r) == 1) && (isEpsilon(r[0]))
}

type Token interface{}

func NewToken(str string) Token {
	if (str[0] == '<') && (str[len(str)-1] == '>') {
		return NTerm(str)
	}

	if str == "ε" {
		return Epsilon{}
	}

	return Term(str)
}

type Term string

type NTerm string

type Epsilon struct{}

func isEpsilon(t interface{}) bool {
	_, ok := t.(Epsilon)
	return ok
}

type TokenSet map[Token]Rule

func (s TokenSet) Add(t Token, r Rule) {
	if _, ok := s[t]; ok {
		return
	}

	s[t] = r
}

func (s TokenSet) AddAll(o TokenSet, r Rule) {
	for t := range o {
		s.Add(t, r)
	}
}

func (s TokenSet) Contains(t Token) bool {
	_, ok := s[t]
	return ok
}

func (s TokenSet) Remove(t Token) {
	delete(s, t)
}
