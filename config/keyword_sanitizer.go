package config

import "strings"

func sanitize(keyword string) string {
	withSpaces := strings.Replace(keyword, "_", " ", -1)
	titled := strings.Title(withSpaces)
	return strings.Replace(titled, " ", "", -1)
}
