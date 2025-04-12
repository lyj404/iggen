package utils

import (
	"strings"
	"unicode"
)

// FuzzySearch
func FuzzySearch(templates []string, term string) []string {
	term = strings.ToLower(term)
	matches := make([]string, 0)

	for _, t := range templates {
		lower := strings.ToLower(t)
		if strings.Contains(lower, term) {
			matches = append(matches, t)
			continue
		}

		if compareSimilar(lower, term) {
			matches = append(matches, t)
		}
	}

	return matches
}

// compareSimilar比较相似的字符
func compareSimilar(a, b string) bool {
	a = removeSpecialChars(a)
	b = removeSpecialChars(b)
	return strings.HasPrefix(a, b) || strings.HasSuffix(a, b)
}

// removeSpecialChars移除特诉字符
func removeSpecialChars(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, s)
}
