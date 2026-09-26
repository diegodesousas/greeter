package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	infrahttp "github.com/diegodesousas/greeter/internal/infra/http"
	"github.com/diegodesousas/greeter/internal/infra/http/handlers"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthReadiness(t *testing.T) {
	pingErr := errors.New("dial tcp: connection refused")

	tests := []struct {
		name      string
		pingErr   error
		wantErrIs []error
	}{
		{
			name: "database reachable returns ok",
		},
		{
			name:      "database unreachable returns service unavailable error",
			pingErr:   pingErr,
			wantErrIs: []error{infrahttp.ErrServiceUnavailable},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewMockPinger(t)
			db.EXPECT().Ping().Return(tt.pingErr).Once()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/readiness", nil)

			err := handlers.HealthReadiness(db)(rec, req)

			if tt.wantErrIs != nil {
				require.Error(t, err)
				for _, target := range tt.wantErrIs {
					assert.ErrorIs(t, err, target)
				}
				assert.Contains(t, err.Error(), tt.pingErr.Error())
				assert.Empty(t, rec.Body.String(), "handler must not write the response on error")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)

			var body map[string]string
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, "ok", body["message"])
		})
	}
}
