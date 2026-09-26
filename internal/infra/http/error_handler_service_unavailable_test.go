package http_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlerServiceUnavailable_Match(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matches ErrServiceUnavailable",
			err:  infrahttp.ErrServiceUnavailable,
			want: true,
		},
		{
			name: "matches wrapped ErrServiceUnavailable",
			err:  fmt.Errorf("%w: database: connection refused", infrahttp.ErrServiceUnavailable),
			want: true,
		},
		{
			name: "does not match other errors",
			err:  errors.New("other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, infrahttp.ErrorHandlerServiceUnavailable{}.Match(tt.err))
		})
	}
}

func TestErrorHandlerServiceUnavailable_Write(t *testing.T) {
	rec := httptest.NewRecorder()

	err := infrahttp.ErrorHandlerServiceUnavailable{}.Write(rec, infrahttp.ErrServiceUnavailable)

	require.NoError(t, err)
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.JSONEq(t, `{"message":"service unavailable"}`, rec.Body.String())
}

func TestErrorHandlerServiceUnavailable_MustLog(t *testing.T) {
	assert.True(t, infrahttp.ErrorHandlerServiceUnavailable{}.MustLog())
}
