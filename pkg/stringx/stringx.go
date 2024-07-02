package stringx

import (
	"strings"
	"unicode"
)

func IsWhiteSpace(s string) bool {
	fields := strings.FieldsFunc(s, unicode.IsSpace)
	return len(fields) == 0
}

func ContainsAny(s string, subStr ...string) bool {
	for _, str := range subStr {
		if strings.Contains(s, str) {
			return true
		}
	}
	return false
}

func TrimWhiteSpace(s string) string {
	r := strings.NewReplacer(" ", "", "\t", "", "\n", "", "\f", "", "\r", "")
	return r.Replace(s)
}
