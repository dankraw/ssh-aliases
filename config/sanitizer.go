package config

import "strings"

type Sanitizer struct{}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{}
}

func (s *Sanitizer) Sanitize(keyword string) string {
	withSpaces := strings.Replace(keyword, "_", " ", -1)
	titled := strings.Title(withSpaces)
	return strings.Replace(titled, " ", "", -1)
}
