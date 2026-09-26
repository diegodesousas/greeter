package list_greetings_test

import (
	"context"
	"errors"
	"testing"
	"time"

	list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

func TestRun(t *testing.T) {
	tests := []struct {
		name       string
		dto        list_greetings.DTO
		setupRepo  func(repo *mocks.MockGreetingRepository)
		wantTotal  int
		wantLen    int
		wantErr    bool
		errContains string
	}{
		{
			name: "success with results",
			dto:  list_greetings.DTO{Page: 1, PerPage: 10},
			setupRepo: func(repo *mocks.MockGreetingRepository) {
				repo.EXPECT().List(mock.Anything, 1, 10).Return([]greeting.Greeting{
					{ID: "abc-123", Name: "Diego", Message: "Hello, Diego!", GreetedAt: fixedTime},
					{ID: "def-456", Name: "Ana", Message: "Hello, Ana!", GreetedAt: fixedTime},
				}, 2, nil)
			},
			wantTotal: 2,
			wantLen:   2,
		},
		{
			name:      "success with empty results",
			dto:       list_greetings.DTO{Page: 1, PerPage: 10},
			setupRepo: func(repo *mocks.MockGreetingRepository) {
				repo.EXPECT().List(mock.Anything, 1, 10).Return([]greeting.Greeting{}, 0, nil)
			},
			wantTotal: 0,
			wantLen:   0,
		},
		{
			name:        "page less than one",
			dto:         list_greetings.DTO{Page: 0, PerPage: 10},
			wantErr:     true,
			errContains: "page",
		},
		{
			name:        "per_page zero",
			dto:         list_greetings.DTO{Page: 1, PerPage: 0},
			wantErr:     true,
			errContains: "per_page",
		},
		{
			name:        "per_page exceeds max",
			dto:         list_greetings.DTO{Page: 1, PerPage: 101},
			wantErr:     true,
			errContains: "per_page",
		},
		{
			name:    "repository error",
			dto:     list_greetings.DTO{Page: 1, PerPage: 10},
			setupRepo: func(repo *mocks.MockGreetingRepository) {
				repo.EXPECT().List(mock.Anything, 1, 10).Return(nil, 0, errors.New("db connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := mocks.NewMockGreetingRepository(t)
			if tt.setupRepo != nil {
				tt.setupRepo(repo)
			}

			useCase := list_greetings.NewUseCase(repo)

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
