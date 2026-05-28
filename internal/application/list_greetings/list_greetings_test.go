package list_greetings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
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

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		dto        list_greetings.DTO
		repo       *mockRepository
		wantTotal  int
		wantLen    int
		wantErr    bool
		errContains string
	}{
		{
			name: "success with results",
			dto:  list_greetings.DTO{Page: 1, PerPage: 10},
			repo: &mockRepository{
				greetings: []greeting.Greeting{
					{ID: "abc-123", Name: "Diego", Message: "Hello, Diego!", GreetedAt: fixedTime},
					{ID: "def-456", Name: "Ana", Message: "Hello, Ana!", GreetedAt: fixedTime},
				},
				total: 2,
			},
			wantTotal: 2,
			wantLen:   2,
		},
		{
			name:      "success with empty results",
			dto:       list_greetings.DTO{Page: 1, PerPage: 10},
			repo:      &mockRepository{greetings: []greeting.Greeting{}, total: 0},
			wantTotal: 0,
			wantLen:   0,
		},
		{
			name:        "page less than one",
			dto:         list_greetings.DTO{Page: 0, PerPage: 10},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "page",
		},
		{
			name:        "per_page zero",
			dto:         list_greetings.DTO{Page: 1, PerPage: 0},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "per_page",
		},
		{
			name:        "per_page exceeds max",
			dto:         list_greetings.DTO{Page: 1, PerPage: 101},
			repo:        &mockRepository{},
			wantErr:     true,
			errContains: "per_page",
		},
		{
			name:    "repository error",
			dto:     list_greetings.DTO{Page: 1, PerPage: 10},
			repo:    &mockRepository{err: errors.New("db connection failed")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := list_greetings.NewUseCase(tt.repo)

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
