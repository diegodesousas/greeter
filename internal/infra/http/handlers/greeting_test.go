package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	"github.com/diegodesousas/greeter/internal/infra/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockGreetUseCase struct {
	capturedDTO greet.DTO
	result      greet.GreetingDTO
	err         error
}

func (m *mockGreetUseCase) Run(_ context.Context, dto greet.DTO) (greet.GreetingDTO, error) {
	m.capturedDTO = dto
	return m.result, m.err
}

func requestWithParam(method, target, key, value string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}

func TestHello(t *testing.T) {
	fixedTime := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		param         string
		useCaseResult greet.GreetingDTO
		useCaseErr    error
		wantErr       bool
		wantMessage   string
	}{
		{
			name:  "success",
			param: "Diego",
			useCaseResult: greet.GreetingDTO{
				Message:   "Hello, Diego!",
				GreetedAt: fixedTime,
			},
			wantMessage: "Hello, Diego!",
		},
		{
			name:       "use case error is propagated",
			param:      "Diego",
			useCaseErr: errors.New("unexpected error"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockGreetUseCase{result: tt.useCaseResult, err: tt.useCaseErr}
			rec := httptest.NewRecorder()
			req := requestWithParam(http.MethodGet, "/hello/Diego", "name", tt.param)

			err := handlers.Hello(mock)(rec, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.param, mock.capturedDTO.Name)
			assert.Equal(t, http.StatusOK, rec.Code)

			var body greet.GreetingDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantMessage, body.Message)
			assert.Equal(t, fixedTime.UTC(), body.GreetedAt.UTC())
		})
	}
}