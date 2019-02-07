package scanner

import "strings"

var (
	symbols = map[string]struct{}{
		".":  {},
		"{":  {},
		"}":  {},
		"[":  {},
		"]":  {},
		"(":  {},
		")":  {},
		"=>": {},
		";":  {},
		":":  {},
		"->": {},
		"--": {},
		"-|": {},
		"(@": {},
	}

	keywords = map[string]struct{}{
		"memo":   {},
		"let":    {},
		"import": {},
	}
)

func isKeyword(str string) bool {
	_, ok := keywords[str]
	return ok
}

func symbolicPrefix(str string) (f string) {
	for k := range symbols {
		if strings.HasPrefix(str, k) {
			if len(k) > len(f) {
				f = k
			}
		}
	}

	return
}

func symbolicSuffix(str string) string {
	for k := range symbols {
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
