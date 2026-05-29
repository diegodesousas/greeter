package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
	search_greetings "github.com/diegodesousas/greeter/internal/application/search_greetings"
	"github.com/diegodesousas/greeter/internal/infra/http/handlers"
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

type mockListGreetingsUseCase struct {
	capturedDTO list_greetings.DTO
	result      list_greetings.Output
	err         error
}

func (m *mockListGreetingsUseCase) Run(_ context.Context, dto list_greetings.DTO) (list_greetings.Output, error) {
	m.capturedDTO = dto
	return m.result, m.err
}

func TestHello(t *testing.T) {
	fixedTime := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		body          string
		useCaseResult greet.GreetingDTO
		useCaseErr    error
		wantErr       bool
		wantName      string
		wantMessage   string
	}{
		{
			name: "success decodes name from json body",
			body: `{"name":"Diego"}`,
			useCaseResult: greet.GreetingDTO{
				Message:   "Hello, Diego!",
				GreetedAt: fixedTime,
			},
			wantName:    "Diego",
			wantMessage: "Hello, Diego!",
		},
		{
			name:       "use case error is propagated",
			body:       `{"name":"Diego"}`,
			useCaseErr: errors.New("unexpected error"),
			wantErr:    true,
		},
		{
			name:    "empty body results in empty name passed to use case",
			body:    "",
			wantErr: false,
			useCaseResult: greet.GreetingDTO{
				Message:   "Hello, !",
				GreetedAt: fixedTime,
			},
			wantName:    "",
			wantMessage: "Hello, !",
		},
		{
			name:    "malformed json body returns error",
			body:    `{"name":`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockGreetUseCase{result: tt.useCaseResult, err: tt.useCaseErr}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader(tt.body))

			err := handlers.Hello(mock)(rec, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantName, mock.capturedDTO.Name)
			assert.Equal(t, http.StatusOK, rec.Code)

			var body greet.GreetingDTO
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.wantMessage, body.Message)
			assert.Equal(t, fixedTime.UTC(), body.GreetedAt.UTC())
		})
	}
}

func TestListGreetings(t *testing.T) {
	fixedTime := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		query          string
		useCaseResult  list_greetings.Output
		useCaseErr     error
		wantErr        bool
		wantPage       int
		wantPerPage    int
	}{
		{
			name:  "success parses page and per_page from query",
			query: "?page=2&per_page=5",
			useCaseResult: list_greetings.Output{
				Data: []list_greetings.GreetingDTO{
					{ID: "abc-123", Message: "Hello, Diego!", GreetedAt: fixedTime},
				},
				Pagination: list_greetings.PaginationDTO{Total: 1, Page: 2, PerPage: 5},
			},
			wantPage:    2,
			wantPerPage: 5,
		},
		{
			name:        "use case error is propagated",
			query:       "?page=1&per_page=10",
			useCaseErr:  errors.New("unexpected error"),
			wantErr:     true,
		},
		{
			name:        "missing query params result in zero values passed to use case",
			query:       "",
			useCaseErr:  errors.New("validation error"),
			wantErr:     true,
			wantPage:    0,
			wantPerPage: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockListGreetingsUseCase{result: tt.useCaseResult, err: tt.useCaseErr}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/greetings"+tt.query, nil)

			err := handlers.ListGreetings(mock)(rec, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.wantPage, mock.capturedDTO.Page)
			assert.Equal(t, tt.wantPerPage, mock.capturedDTO.PerPage)

			var body list_greetings.Output
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.useCaseResult.Pagination.Total, body.Pagination.Total)
		})
	}
}

type mockSearchGreetingsUseCase struct {
	capturedDTO search_greetings.DTO
	result      search_greetings.Output
	err         error
}

func (m *mockSearchGreetingsUseCase) Run(_ context.Context, dto search_greetings.DTO) (search_greetings.Output, error) {
	m.capturedDTO = dto
	return m.result, m.err
}

func TestSearchGreetings(t *testing.T) {
	fixedTime := time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		query         string
		useCaseResult search_greetings.Output
		useCaseErr    error
		wantErr       bool
		wantName      string
		wantPage      int
		wantPerPage   int
	}{
		{
			name:  "success parses name, page and per_page from query",
			query: "?name=jo%C3%A3o&page=1&per_page=10",
			useCaseResult: search_greetings.Output{
				Data: []search_greetings.GreetingDTO{
					{ID: "abc-123", Message: "Hello, João!", GreetedAt: fixedTime},
				},
				Pagination: search_greetings.PaginationDTO{Total: 1, Page: 1, PerPage: 10},
			},
			wantName:    "joão",
			wantPage:    1,
			wantPerPage: 10,
		},
		{
			name:       "use case error is propagated",
			query:      "?name=diego&page=1&per_page=10",
			useCaseErr: errors.New("unexpected error"),
			wantErr:    true,
		},
		{
			name:        "missing name results in empty string passed to use case",
			query:       "?page=1&per_page=10",
			useCaseErr:  errors.New("validation error"),
			wantErr:     true,
			wantName:    "",
			wantPage:    1,
			wantPerPage: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockSearchGreetingsUseCase{result: tt.useCaseResult, err: tt.useCaseErr}
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/greetings/search"+tt.query, nil)

			err := handlers.SearchGreetings(mock)(rec, req)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.wantName, mock.capturedDTO.Name)
			assert.Equal(t, tt.wantPage, mock.capturedDTO.Page)
			assert.Equal(t, tt.wantPerPage, mock.capturedDTO.PerPage)

			var body search_greetings.Output
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, tt.useCaseResult.Pagination.Total, body.Pagination.Total)
		})
	}
}