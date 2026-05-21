package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type errWriter struct {
	header http.Header
}

func (e errWriter) Header() http.Header        { return e.header }
func (e errWriter) Write([]byte) (int, error)  { return 0, errors.New("write failed") }
func (e errWriter) WriteHeader(int)            {}

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		writer  func() http.ResponseWriter
		wantErr bool
		wantBody string
	}{
		{
			name:     "serializes struct to json",
			input:    struct {
				Message string `json:"message"`
			}{Message: "ok"},
			writer:   func() http.ResponseWriter { return httptest.NewRecorder() },
			wantBody: `{"message":"ok"}`,
		},
		{
			name:    "returns error when writer fails",
			input:   struct{}{},
			writer:  func() http.ResponseWriter { return errWriter{header: http.Header{}} },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := tt.writer()

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