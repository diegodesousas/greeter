package validation_test

import (
	"strings"
	"testing"

	"github.com/diegodesousas/greeter/internal/application/validation"
	"github.com/stretchr/testify/assert"
)

func TestExceedsMaxLength(t *testing.T) {
	tests := []struct {
		name  string
		value string
		max   int
		want  bool
	}{
		{name: "empty value does not exceed", value: "", max: 5, want: false},
		{name: "value below max does not exceed", value: "abc", max: 5, want: false},
		{name: "value at max does not exceed", value: "abcde", max: 5, want: false},
		{name: "value above max exceeds", value: "abcdef", max: 5, want: true},
		{name: "accented characters count as one", value: strings.Repeat("ã", 5), max: 5, want: false},
		{name: "multi-byte characters above max exceed", value: strings.Repeat("日", 6), max: 5, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, validation.ExceedsMaxLength(tt.value, tt.max))
		})
	}
}
