package greet_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fixedClock struct{ t time.Time }

func (f fixedClock) Now() time.Time { return f.t }

type mockRepository struct{ err error }

func (m *mockRepository) Save(_ context.Context, _ greeting.Greeting) error { return m.err }

func (m *mockRepository) List(_ context.Context, _, _ int) ([]greeting.Greeting, int, error) {
	return nil, 0, m.err
}

func (m *mockRepository) Search(_ context.Context, _ string, _, _ int) ([]greeting.Greeting, int, error) {
	return nil, 0, m.err
}

var fixedTime = time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		dto         greet.DTO
		repo        *mockRepository
		wantMessage string
		wantErr     bool
		errContains string
	}{
		{
			name:        "success",
			dto:         greet.DTO{Name: "Diego"},
			repo:        &mockRepository{},
			wantMessage: "Hello, Diego!",
		},
		{
			name:        "name at max length",
			dto:         greet.DTO{Name: strings.Repeat("a", 50)},
			repo:        &mockRepository{},
			wantMessage: "Hello, " + strings.Repeat("a", 50) + "!",
		},
		{
			name:        "name empty",
			dto:         greet.DTO{Name: ""},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "name",
		},
		{
			name:        "name exceeds max length",
			dto:         greet.DTO{Name: strings.Repeat("a", 51)},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "name",
		},
		{
			name:    "repository error",
			dto:     greet.DTO{Name: "Diego"},
			repo:    &mockRepository{err: errors.New("db connection failed")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := greet.NewUseCase(fixedClock{t: fixedTime}, tt.repo)

			result, err := useCase.Run(context.Background(), tt.dto)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantMessage, result.Message)
			assert.Equal(t, fixedTime, result.GreetedAt)
		})
	}
}