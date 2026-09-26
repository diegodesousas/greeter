package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewErrorHandler(t *testing.T) {
	tests := []struct {
		name        string
		setupWriter func(writer *mocks.MockErrorWriter)
		wantStatus  int
	}{
		{
			name: "calls matching writer",
			setupWriter: func(writer *mocks.MockErrorWriter) {
				writer.EXPECT().Match(mock.Anything).Return(true)
				writer.EXPECT().Write(mock.Anything, mock.Anything).RunAndReturn(func(w http.ResponseWriter, _ error) error {
					w.WriteHeader(http.StatusTeapot)
					return nil
				})
				writer.EXPECT().MustLog().Return(false)
			},
			wantStatus: http.StatusTeapot,
		},
		{
			name: "falls back to 500 when no writer matches",
			setupWriter: func(writer *mocks.MockErrorWriter) {
				writer.EXPECT().Match(mock.Anything).Return(false)
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "falls back to 500 when writer returns error",
			setupWriter: func(writer *mocks.MockErrorWriter) {
				writer.EXPECT().Match(mock.Anything).Return(true)
				writer.EXPECT().Write(mock.Anything, mock.Anything).Return(errors.New("write failed"))
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := mocks.NewMockErrorWriter(t)
			tt.setupWriter(writer)
			rec := httptest.NewRecorder()
			handler := infrahttp.NewErrorHandler(writer)

			handler(context.Background(), rec, errors.New("some error"))

			assert.Equal(t, tt.wantStatus, rec.Code)
		})
	}
}
