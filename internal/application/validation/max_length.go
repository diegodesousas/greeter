package validation

import "unicode/utf8"

// ExceedsMaxLength reports whether value has more than max characters. Length
// is counted in runes, not bytes, so accented and multi-byte characters count
// as one.
func ExceedsMaxLength(value string, max int) bool {
	return utf8.RuneCountInString(value) > max
}
