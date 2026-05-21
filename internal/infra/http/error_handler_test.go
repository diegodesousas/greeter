package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/stretchr/testify/assert"
)

type mockErrorWriter struct {
	match    bool
	writeErr error
	mustLog  bool
	written  bool
}

func (m *mockErrorWriter) Match(_ error) bool { return m.match }
func (m *mockErrorWriter) MustLog() bool      { return m.mustLog }
func (m *mockErrorWriter) Write(w http.ResponseWriter, _ error) error {
	m.written = true
	if m.writeErr != nil {
		return m.writeErr
	}
	w.WriteHeader(http.StatusTeapot)
	return nil
}

func TestNewErrorHandler(t *testing.T) {
	tests := []struct {
		name       string
		writer     *mockErrorWriter
		wantStatus int
		wantWritten bool
	}{
		{
			name:        "calls matching writer",
			writer:      &mockErrorWriter{match: true},
			wantStatus:  http.StatusTeapot,
			wantWritten: true,
		},
		{
			name:        "falls back to 500 when no writer matches",
			writer:      &mockErrorWriter{match: false},
			wantStatus:  http.StatusInternalServerError,
			wantWritten: false,
		},
		{
			name:        "falls back to 500 when writer returns error",
			writer:      &mockErrorWriter{match: true, writeErr: errors.New("write failed")},
			wantStatus:  http.StatusInternalServerError,
			wantWritten: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			handler := infrahttp.NewErrorHandler(tt.writer)

			handler(context.Background(), rec, errors.New("some error"))

			assert.Equal(t, tt.wantStatus, rec.Code)
			assert.Equal(t, tt.wantWritten, tt.writer.written)
		})
	}
}