package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// A Grammar is a map from non-terminals to lists of rules.
type Grammar map[NTerm][]Rule

// LoadGrammar loads a grammar from an io.Reader.
func LoadGrammar(r io.Reader, detectAmbiguity bool) (g Grammar, err error) {
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

	if detectAmbiguity {
		g.detectAmbiguities()
	}

	return g, s.Err()
}

func (g Grammar) detectAmbiguities() {
	type lookup struct {
		NTerm Token
		Term  Token
	}
	known := make(map[lookup]struct{}, len(g))

	for nterm := range g {
		for term, rule := range g.First(nterm) {
			if isEpsilon(term) {
				continue
			}

			_, ambiguous := known[lookup{nterm, term}]
			if ambiguous {
				fmt.Fprintf(os.Stderr, "Ambiguity for (%v, %v) from %v\n", nterm, term, rule)
			}

			known[lookup{nterm, term}] = struct{}{}
		}

		if g.Nullable(nterm) {
			for term, rule := range g.Follow(nterm) {
				if isEpsilon(term) {
					continue
				}

				_, ambiguous := known[lookup{nterm, term}]
				if ambiguous {
					fmt.Fprintf(os.Stderr, "Ambiguity for (%v, %v) from %v\n", nterm, term, rule)
				}

				known[lookup{nterm, term}] = struct{}{}
			}
		}
	}
}

// Nullable returns true if the given token is nullable according to
// the grammar.
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

	case EOF:
		return false
	}

	panic(fmt.Errorf("Unexpected token type: %T", tok))
}

// First returns the first set of the given token. For terminals, the
// first set only contains the terminal in question. For
// non-terminals, the first set is the set of all terminals which can
// appear first in rules that the non-terminal maps to, including ε.
func (g Grammar) First(tok Token) TokenSet {
	ts := make(TokenSet)

	switch tok := tok.(type) {
	case Term, Epsilon, EOF:
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

// Follow returns the follow set of the given non-terminal. The follow
// set consists of every terminal which can appear immediately after
// the non-terminal in the grammar, excluding ε.
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

// A Rule is an ordered list of tokens.
type Rule []Token

// Epsilon returns true if the rule contains only a single ε.
func (r Rule) Epsilon() bool {
	return (len(r) == 1) && isEpsilon(r[0])
}

// A Token represnts one of four things:
// * A terminal, which is anything that the scanner sends to the parser.
// * A non-terminal, which is structure in the grammar's definition.
// * Epsilon, written as ε, which represents a non-action.
// * EOF, written as Ω.
type Token interface{}

// NewToken determines the type of a token from the given string and
// returns the appropriate token.
func NewToken(str string) Token {
	if (str[0] == '<') && (str[len(str)-1] == '>') {
		return NTerm(str)
	}

	switch str {
	case "ε":
		return Epsilon{}
	case "Ω":
		return EOF{}
	}

	return Term(str)
}

// Term is a terminal token.
type Term string

// NTerm is a non-terminal token.
type NTerm string

// Epsilon is an ε token.
type Epsilon struct{}

func isEpsilon(t Token) bool {
	_, ok := t.(Epsilon)
	return ok
}

// EOF is an Ω token.
type EOF struct{}

// A TokenSet is an unordered set of tokens mapped to the rules that
// put them there.
type TokenSet map[Token]Rule

// Add adds a token to the token set, mapping it to the given rule. If
// the token is already in the token set, this is a no-op.
func (s TokenSet) Add(t Token, r Rule) {
	if _, ok := s[t]; ok {
		return
	}

	s[t] = r
}

// AddAll adds all of the tokens from another token set to this one,
// mapping all of them to the given rule.
func (s TokenSet) AddAll(o TokenSet, r Rule) {
	for t := range o {
		s.Add(t, r)
	}
}

// Contains returns true if t is in s.
func (s TokenSet) Contains(t Token) bool {
	_, ok := s[t]
	return ok
}

// Remove removes t from s.
func (s TokenSet) Remove(t Token) {
	delete(s, t)
}
