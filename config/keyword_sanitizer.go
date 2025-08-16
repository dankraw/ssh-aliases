package config

import (
	"strings"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func sanitize(keyword string) string {
	withSpaces := strings.ReplaceAll(keyword, "_", " ")
	titled := cases.Title(language.English, cases.NoLower).String(withSpaces)
	return strings.ReplaceAll(titled, " ", "")
}
