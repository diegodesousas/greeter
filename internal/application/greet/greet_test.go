package greet_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		dto         greet.DTO
		setupRepo   func(repo *mocks.MockGreetingRepository)
		wantMessage string
		wantErr     bool
		errContains string
	}{
		{
			name: "success",
			dto:  greet.DTO{Name: "Diego"},
			setupRepo: func(repo *mocks.MockGreetingRepository) {
				repo.EXPECT().Save(mock.Anything, greeting.New("Diego", fixedTime)).Return(nil)
			},
			wantMessage: "Hello, Diego!",
		},
		{
			name: "name at max length",
			dto:  greet.DTO{Name: strings.Repeat("a", 50)},
			setupRepo: func(repo *mocks.MockGreetingRepository) {
				repo.EXPECT().Save(mock.Anything, greeting.New(strings.Repeat("a", 50), fixedTime)).Return(nil)
			},
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
		{
			name: "repository error",
			dto:  greet.DTO{Name: "Diego"},
			setupRepo: func(repo *mocks.MockGreetingRepository) {
				repo.EXPECT().Save(mock.Anything, mock.Anything).Return(errors.New("db connection failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clock := mocks.NewMockClock(t)
			clock.EXPECT().Now().Return(fixedTime).Maybe()

			repo := mocks.NewMockGreetingRepository(t)
			if tt.setupRepo != nil {
				tt.setupRepo(repo)
			}

			useCase := greet.NewUseCase(clock, repo)

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
