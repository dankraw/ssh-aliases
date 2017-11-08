package config

import "strings"

type KeywordSanitizer struct{}

func NewKeywordSanitizer() *KeywordSanitizer {
	return &KeywordSanitizer{}
}

func (s *KeywordSanitizer) Sanitize(keyword string) string {
	withSpaces := strings.Replace(keyword, "_", " ", -1)
	titled := strings.Title(withSpaces)
	return strings.Replace(titled, " ", "", -1)
}
