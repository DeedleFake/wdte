package scanner

import "strings"

var (
	symbols = map[string]struct{}{
		".":  {},
		"->": {},
		"{":  {},
		"}":  {},
		"[":  {},
		"]":  {},
		"(":  {},
		")":  {},
		"=>": {},
		";":  {},
	}

	keywords = map[string]struct{}{
		"switch":  {},
		"default": {},
		"memo":    {},
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
