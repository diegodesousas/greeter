package greet_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixedClock struct{ t time.Time }

func (f fixedClock) Now() time.Time { return f.t }

var fixedTime = time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		dto         greet.DTO
		wantMessage string
		wantErr     bool
		errContains string
	}{
		{
			name:        "success",
			dto:         greet.DTO{Name: "Diego"},
			wantMessage: "Hello, Diego!",
		},
		{
			name:        "name at max length",
			dto:         greet.DTO{Name: strings.Repeat("a", 50)},
			wantMessage: "Hello, " + strings.Repeat("a", 50) + "!",
		},
		{
			name:        "name empty",
			dto:         greet.DTO{Name: ""},
			wantErr:     true,
			errContains: "name",
		},
		{
			name:        "name exceeds max length",
			dto:         greet.DTO{Name: strings.Repeat("a", 51)},
			wantErr:     true,
			errContains: "name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := greet.NewUseCase(fixedClock{t: fixedTime})

			result, err := useCase.Run(context.Background(), tt.dto)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantMessage, result.Message)
			assert.Equal(t, fixedTime, result.GreetedAt)
		})
	}
}