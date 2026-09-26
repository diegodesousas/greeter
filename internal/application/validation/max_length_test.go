package validation_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/diegodesousas/greeter/internal/application/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type maxLengthDTO struct {
	Title string
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name        string
		max         int
		input       string
		wantErr     bool
		wantMessage string
	}{
		{name: "empty value is valid", max: 5, input: "", wantErr: false},
		{name: "value below max is valid", max: 5, input: "abc", wantErr: false},
		{name: "value at max is valid", max: 5, input: "abcde", wantErr: false},
		{name: "value above max fails", max: 5, input: "abcdef", wantErr: true, wantMessage: "attribute title must have at most 5 characters"},
		{name: "multi-byte value at max is valid", max: 5, input: strings.Repeat("ã", 5), wantErr: false},
		{name: "multi-byte value above max fails", max: 5, input: strings.Repeat("日", 6), wantErr: true, wantMessage: "attribute title must have at most 5 characters"},
		{name: "message uses the configured max", max: 2, input: "abc", wantErr: true, wantMessage: "attribute title must have at most 2 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := validation.MaxLength("title", tt.max, func(dto maxLengthDTO) string { return dto.Title })

			err := rule(context.Background(), maxLengthDTO{Title: tt.input})

			if !tt.wantErr {
				require.NoError(t, err)
				return
			}

			var validationErr validator.Error
			require.True(t, errors.As(err, &validationErr), "expected validator.Error, got %v", err)
			assert.Equal(t, validator.ErrorCode("max_length"), validationErr.Code)
			assert.Equal(t, tt.wantMessage, validationErr.Message)
		})
	}
}
