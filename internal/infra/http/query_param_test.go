package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegodesousas/go-devkit/pkg/validator"
	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryInt(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		fallback    int
		want        int
		wantErrMsg  string
	}{
		{
			name:     "parses integer value",
			query:    "?page=3",
			fallback: 1,
			want:     3,
		},
		{
			name:     "missing param returns fallback",
			query:    "",
			fallback: 1,
			want:     1,
		},
		{
			name:     "empty param returns fallback",
			query:    "?page=",
			fallback: 1,
			want:     1,
		},
		{
			name:     "explicit zero is kept for validation downstream",
			query:    "?page=0",
			fallback: 1,
			want:     0,
		},
		{
			name:       "non-numeric value returns validation error",
			query:      "?page=abc",
			fallback:   1,
			wantErrMsg: "attribute page must be an integer",
		},
		{
			name:       "decimal value returns validation error",
			query:      "?page=1.5",
			fallback:   1,
			wantErrMsg: "attribute page must be an integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.query, nil)

			got, err := infrahttp.QueryInt(req, "page", tt.fallback)

			if tt.wantErrMsg != "" {
				var validationErr validator.Error
				require.True(t, errors.As(err, &validationErr), "expected validator.Error, got %v", err)
				assert.Equal(t, tt.wantErrMsg, validationErr.Message)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
