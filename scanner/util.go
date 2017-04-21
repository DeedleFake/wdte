package scanner

import "strings"

var (
	keywords = map[string]struct{}{
		".":  struct{}{},
		"->": struct{}{},
		"(":  struct{}{},
		")":  struct{}{},
		"{":  struct{}{},
		"}":  struct{}{},
		"[":  struct{}{},
		"]":  struct{}{},
		"=":  struct{}{},
		";":  struct{}{},
	}
)

func isKeyword(str string) bool {
	_, ok := keywords[str]
	return ok
}

func getKeywordSuffix(str string) string {
	for k := range keywords {
		if strings.HasSuffix(str, k) {
			return k
		}
	}

	return ""
}

func isQuote(r rune) bool {
	return (r == '\'') || (r == '"')
}
