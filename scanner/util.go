package scanner

import "strings"

var (
	symbols = map[string]struct{}{
		".":  struct{}{},
		"->": struct{}{},
		"{":  struct{}{},
		"}":  struct{}{},
		"[":  struct{}{},
		"]":  struct{}{},
		":=": struct{}{},
		";":  struct{}{},
	}

	keywords = map[string]struct{}{
		"switch":  struct{}{},
		"default": struct{}{},
	}
)

func isKeyword(str string) bool {
	_, ok := keywords[str]
	return ok
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
