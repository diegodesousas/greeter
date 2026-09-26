package validation

import (
	"context"
	"fmt"
	"unicode/utf8"

	"github.com/diegodesousas/go-devkit/pkg/validator"
)

// MaxLength returns a rule that fails when the string selected by value has
// more than max characters. Length is counted in runes, not bytes, so accented
// and multi-byte characters count as one.
func MaxLength[T any](field string, max int, value func(T) string) validator.Rule[T] {
	return func(_ context.Context, dto T) error {
		if utf8.RuneCountInString(value(dto)) > max {
			return validator.Error{
				Code:    "max_length",
				Message: fmt.Sprintf("attribute %s must have at most %d characters", field, max),
			}
		}

		return nil
	}
}
