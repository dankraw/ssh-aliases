package config

import "strings"

type keywordSanitizer struct{}

func newKeywordSanitizer() *keywordSanitizer {
	return &keywordSanitizer{}
}

func (s *keywordSanitizer) sanitize(keyword string) string {
	withSpaces := strings.Replace(keyword, "_", " ", -1)
	titled := strings.Title(withSpaces)
	return strings.Replace(titled, " ", "", -1)
}
