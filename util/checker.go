package util

import (
	"unicode/utf8"
)

// ValidLine ...
func ValidLine(raw []byte) (bool, int) {
	pos := 1

	for len(raw) > 0 {
		if r, size := utf8.DecodeRune(raw); utf8.RuneError != r {
			raw = raw[size:]
			pos += size
		} else {
			return false, pos
		}
	}
	return true, 0
}
