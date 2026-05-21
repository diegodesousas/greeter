package http_test

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "sql.ErrNoRows maps to 404",
			err:        sql.ErrNoRows,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "validator.Error maps to 422",
			err:        validator.NewRequiredError("name"),
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "unknown error maps to 500",
			err:        errors.New("unexpected error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			infrahttp.ErrorHandler(context.Background(), rec, tt.err)

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}