package utils

import "strings"

func RemoveWildcards(s string) string {
	s = strings.ReplaceAll(s, "*", "")
	return s
}
