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

func TestErrorHandlerBadRequest_Match(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matches ErrBadRequest",
			err:  infrahttp.ErrBadRequest,
			want: true,
		},
		{
			name: "matches wrapped ErrBadRequest",
			err:  fmt.Errorf("%w: unexpected EOF", infrahttp.ErrBadRequest),
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
			assert.Equal(t, tt.want, infrahttp.ErrorHandlerBadRequest{}.Match(tt.err))
		})
	}
}

func TestErrorHandlerBadRequest_Write(t *testing.T) {
	rec := httptest.NewRecorder()

	err := infrahttp.ErrorHandlerBadRequest{}.Write(rec, infrahttp.ErrBadRequest)

	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"message":"invalid request body"}`, rec.Body.String())
}

func TestErrorHandlerBadRequest_MustLog(t *testing.T) {
	assert.False(t, infrahttp.ErrorHandlerBadRequest{}.MustLog())
}
