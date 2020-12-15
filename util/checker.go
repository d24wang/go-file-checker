package util

import (
	"unicode/utf8"
)

// ValidLine ...
func ValidLine(raw []byte) (bool, int) {
	pos := 1

	for len(raw) > 0 {
		r, size := utf8.DecodeRune(raw)
		if utf8.RuneError == r {
			return false, pos
		}
		raw = raw[size:]
		pos += size
	}
	return true, 0
}
