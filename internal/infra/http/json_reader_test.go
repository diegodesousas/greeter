package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type readJsonPayload struct {
	Name string `json:"name"`
}

func TestReadJson(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		wantErrIs error
		want      readJsonPayload
	}{
		{
			name: "decodes valid json body",
			body: `{"name":"Diego"}`,
			want: readJsonPayload{Name: "Diego"},
		},
		{
			name: "empty body leaves destination untouched",
			body: "",
			want: readJsonPayload{},
		},
		{
			name:      "malformed json returns bad request error",
			body:      `{"name":`,
			wantErrIs: infrahttp.ErrBadRequest,
		},
		{
			name:      "invalid field type returns bad request error",
			body:      `{"name":123}`,
			wantErrIs: infrahttp.ErrBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))

			var got readJsonPayload
			err := infrahttp.ReadJson(req, &got)

			if tt.wantErrIs != nil {
				require.ErrorIs(t, err, tt.wantErrIs)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
