package scanner

func isKeyword(str string) bool {
	switch str {
	case ".", "->", "(", ")", "{", "}", "[", "]", "=", ";":
		return true
	}

	return false
}

func isQuote(r rune) bool {
	return (r == '\'') || (r == '"')
}
