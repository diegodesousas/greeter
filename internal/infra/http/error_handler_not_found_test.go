package http_test

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlerNotFound_Match(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matches sql.ErrNoRows",
			err:  sql.ErrNoRows,
			want: true,
		},
		{
			name: "matches wrapped sql.ErrNoRows",
			err:  pkgerrors.Wrap(sql.ErrNoRows, "user not found"),
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
			assert.Equal(t, tt.want, infrahttp.ErrorHandlerNotFound{}.Match(tt.err))
		})
	}
}

func TestErrorHandlerNotFound_Write(t *testing.T) {
	rec := httptest.NewRecorder()

	err := infrahttp.ErrorHandlerNotFound{}.Write(rec, sql.ErrNoRows)

	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"message":"resource not found"}`, rec.Body.String())
}

func TestErrorHandlerNotFound_MustLog(t *testing.T) {
	assert.False(t, infrahttp.ErrorHandlerNotFound{}.MustLog())
}