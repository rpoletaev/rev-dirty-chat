package helper

import (
	"strings"
)

func StringsContains(src []string, s string, ignoreCase bool) bool {
	if ignoreCase {
		return stringsContainsIgnoreKeys(src, s)
	}

	for i := 0; i < len(src); i++ {
		if strings.Compare(s, src[i]) == 0 {
			return true
		}
	}

	return false
}

func stringsContainsIgnoreKeys(src []string, s string) bool {
	s = strings.ToLower(s)

	for i := 0; i < len(src); i++ {
		if strings.Compare(s, strings.ToLower(src[i])) == 0 {
			return true
		}
	}

	return false
}

//StringsWithoutFirstEntry returns new slice which not contains first entry of s
func StringsWithoutFirstEntry(src []string, s string) []string {
	dest := make([]string, len(src))
	copy(dest, src)
	for i, str := range src {
		if str == s {
			dest = append(dest[:i], dest[i+1:]...)
			return dest
		}
	}

	return src
}
