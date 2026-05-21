package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	pkgerrors "github.com/pkg/errors"
	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorHandlerValidation_Match(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "matches validator.Error",
			err:  validator.NewRequiredError("name"),
			want: true,
		},
		{
			name: "matches wrapped validator.Error",
			err:  pkgerrors.Wrap(validator.NewRequiredError("name"), "validation failed"),
			want: true,
		},
		{
			name: "does not match regular error",
			err:  errors.New("other error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, infrahttp.ErrorHandlerValidation{}.Match(tt.err))
		})
	}
}

func TestErrorHandlerValidation_Write(t *testing.T) {
	rec := httptest.NewRecorder()
	validationErr := validator.NewRequiredError("name")

	err := infrahttp.ErrorHandlerValidation{}.Write(rec, validationErr)

	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.JSONEq(t, `{"message":"attribute name is required"}`, rec.Body.String())
}

func TestErrorHandlerValidation_MustLog(t *testing.T) {
	assert.False(t, infrahttp.ErrorHandlerValidation{}.MustLog())
}