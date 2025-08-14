package scanner

import (
	"slices"
	"sort"
	"strings"
)

var (
	symbols = []string{
		".",
		"{",
		"}",
		"[",
		"]",
		"(",
		")",
		"(|",
		"|)",
		"=>",
		";",
		":",
		"->",
		"--",
		"-|",
		"(@",
	}

	keywords = []string{
		"let",
		"import",
	}
)

func init() {
	sort.Slice(symbols, func(i1, i2 int) bool { return len(symbols[i1]) > len(symbols[i2]) })
}

func isKeyword(str string) bool {
	return slices.Contains(keywords, str)
}

func symbolicPrefix(str string) (f string) {
	for _, k := range symbols {
		if strings.HasPrefix(str, k) {
			if len(k) > len(f) {
				f = k
			}
		}
	}

	return
}

func symbolicSuffix(str string) string {
	for _, k := range symbols {
		if strings.HasSuffix(str, k) {
			return k
		}
	}

	return ""
}

func isQuote(r rune) bool {
	return (r == '\'') || (r == '"')
}

func endQuote(r rune) rune {
	switch r {
	case '(':
		return ')'
	case ')':
		return '('
	case '[':
		return ']'
	case ']':
		return '['
	case '{':
		return '}'
	case '}':
		return '{'

	default:
		return r
	}
}
