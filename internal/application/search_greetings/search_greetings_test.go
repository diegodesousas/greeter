package search_greetings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	search_greetings "github.com/diegodesousas/greeter/internal/application/search_greetings"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

type mockRepository struct {
	greetings []greeting.Greeting
	total     int
	err       error
}

func (m *mockRepository) Save(_ context.Context, _ greeting.Greeting) error { return m.err }

func (m *mockRepository) List(_ context.Context, _, _ int) ([]greeting.Greeting, int, error) {
	return m.greetings, m.total, m.err
}

func (m *mockRepository) Search(_ context.Context, _ string, _, _ int) ([]greeting.Greeting, int, error) {
	return m.greetings, m.total, m.err
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		dto         search_greetings.DTO
		repo        *mockRepository
		wantTotal   int
		wantLen     int
		wantErr     bool
		errContains string
	}{
		{
			name: "success with results",
			dto:  search_greetings.DTO{Name: "joão", Page: 1, PerPage: 10},
			repo: &mockRepository{
				greetings: []greeting.Greeting{
					{ID: "abc-123", Name: "João", Message: "Hello, João!", GreetedAt: fixedTime},
					{ID: "def-456", Name: "Joana", Message: "Hello, Joana!", GreetedAt: fixedTime},
				},
				total: 2,
			},
			wantTotal: 2,
			wantLen:   2,
		},
		{
			name:      "success with empty results",
			dto:       search_greetings.DTO{Name: "xpto", Page: 1, PerPage: 10},
			repo:      &mockRepository{greetings: []greeting.Greeting{}, total: 0},
			wantTotal: 0,
			wantLen:   0,
		},
		{
			name:        "empty name fails validation",
			dto:         search_greetings.DTO{Name: "", Page: 1, PerPage: 10},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "name",
		},
		{
			name:        "page less than one fails validation",
			dto:         search_greetings.DTO{Name: "diego", Page: 0, PerPage: 10},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "page",
		},
		{
			name:        "per_page zero fails validation",
			dto:         search_greetings.DTO{Name: "diego", Page: 1, PerPage: 0},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "per_page",
		},
		{
			name:        "per_page exceeds max fails validation",
			dto:         search_greetings.DTO{Name: "diego", Page: 1, PerPage: 101},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "per_page",
		},
		{
			name:    "repository error is propagated",
			dto:     search_greetings.DTO{Name: "diego", Page: 1, PerPage: 10},
			repo:    &mockRepository{err: errors.New("db connection failed")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := search_greetings.NewUseCase(tt.repo)

			result, err := useCase.Run(context.Background(), tt.dto)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, result.Pagination.Total)
			assert.Len(t, result.Data, tt.wantLen)
			assert.Equal(t, tt.dto.Page, result.Pagination.Page)
			assert.Equal(t, tt.dto.PerPage, result.Pagination.PerPage)
		})
	}
}
