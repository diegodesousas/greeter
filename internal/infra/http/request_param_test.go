package http_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diegodesousas/go-devkit/pkg/validator"
	"github.com/go-chi/chi/v5"
	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQueryString(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  string
	}{
		{name: "returns value", query: "?name=diego", want: "diego"},
		{name: "trims surrounding whitespace", query: "?name=%20%20diego%09", want: "diego"},
		{name: "keeps inner whitespace", query: "?name=%20jo%C3%A3o%20silva%20", want: "joão silva"},
		{name: "whitespace-only value becomes empty", query: "?name=%20%20", want: ""},
		{name: "missing param returns empty", query: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.query, nil)

			assert.Equal(t, tt.want, infrahttp.QueryString(req, "name"))
		})
	}
}

func TestPathParam(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{name: "returns value", value: "diego", want: "diego"},
		{name: "trims surrounding whitespace", value: "  diego\t", want: "diego"},
		{name: "whitespace-only value becomes empty", value: "   ", want: ""},
		{name: "missing param returns empty", value: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			routeCtx := chi.NewRouteContext()
			if tt.value != "" {
				routeCtx.URLParams.Add("name", tt.value)
			}
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			assert.Equal(t, tt.want, infrahttp.PathParam(req, "name"))
		})
	}
}

func TestQueryInt(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		fallback    int
		want        int
		wantErrMsg  string
	}{
		{
			name:     "parses integer value",
			query:    "?page=3",
			fallback: 1,
			want:     3,
		},
		{
			name:     "missing param returns fallback",
			query:    "",
			fallback: 1,
			want:     1,
		},
		{
			name:     "empty param returns fallback",
			query:    "?page=",
			fallback: 1,
			want:     1,
		},
		{
			name:     "explicit zero is kept for validation downstream",
			query:    "?page=0",
			fallback: 1,
			want:     0,
		},
		{
			name:       "non-numeric value returns validation error",
			query:      "?page=abc",
			fallback:   1,
			wantErrMsg: "attribute page must be an integer",
		},
		{
			name:       "decimal value returns validation error",
			query:      "?page=1.5",
			fallback:   1,
			wantErrMsg: "attribute page must be an integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/"+tt.query, nil)

			got, err := infrahttp.QueryInt(req, "page", tt.fallback)

			if tt.wantErrMsg != "" {
				var validationErr validator.Error
				require.True(t, errors.As(err, &validationErr), "expected validator.Error, got %v", err)
				assert.Equal(t, tt.wantErrMsg, validationErr.Message)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
