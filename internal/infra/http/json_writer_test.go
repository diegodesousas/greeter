package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		writer  func(t *testing.T) http.ResponseWriter
		wantErr bool
		wantBody string
	}{
		{
			name:     "serializes struct to json",
			input:    struct {
				Message string `json:"message"`
			}{Message: "ok"},
			writer:   func(_ *testing.T) http.ResponseWriter { return httptest.NewRecorder() },
			wantBody: `{"message":"ok"}`,
		},
		{
			name:    "returns error when writer fails",
			input:   struct{}{},
			writer: func(t *testing.T) http.ResponseWriter {
				w := mocks.NewMockResponseWriter(t)
				w.EXPECT().Write(mock.Anything).Return(0, errors.New("write failed"))
				return w
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.writer(t)

			err := infrahttp.WriteJson(w, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			rec := w.(*httptest.ResponseRecorder)
			assert.JSONEq(t, tt.wantBody, rec.Body.String())
		})
	}
}