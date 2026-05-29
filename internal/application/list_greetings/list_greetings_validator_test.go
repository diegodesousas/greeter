package list_greetings_test

import (
	"context"
	"testing"

	list_greetings "github.com/diegodesousas/greeter/internal/application/list_greetings"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopRepository struct{}

func (n *noopRepository) Save(_ context.Context, _ greeting.Greeting) error { return nil }

func (n *noopRepository) List(_ context.Context, _, _ int) ([]greeting.Greeting, int, error) {
	return []greeting.Greeting{}, 0, nil
}

func (n *noopRepository) Search(_ context.Context, _ string, _, _ int) ([]greeting.Greeting, int, error) {
	return []greeting.Greeting{}, 0, nil
}

func TestValidator_Page(t *testing.T) {
	tests := []struct {
		name        string
		page        int
		wantErr     bool
		errContains string
	}{
		{name: "page one is valid", page: 1, wantErr: false},
		{name: "page greater than one is valid", page: 5, wantErr: false},
		{name: "page zero fails", page: 0, wantErr: true, errContains: "page"},
		{name: "negative page fails", page: -1, wantErr: true, errContains: "page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := list_greetings.NewUseCase(&noopRepository{})
			_, err := useCase.Run(context.Background(), list_greetings.DTO{Page: tt.page, PerPage: 10})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestValidator_PerPage(t *testing.T) {
	tests := []struct {
		name        string
		perPage     int
		wantErr     bool
		errContains string
	}{
		{name: "per_page one is valid", perPage: 1, wantErr: false},
		{name: "per_page 100 is valid", perPage: 100, wantErr: false},
		{name: "per_page 50 is valid", perPage: 50, wantErr: false},
		{name: "per_page zero fails", perPage: 0, wantErr: true, errContains: "per_page"},
		{name: "per_page 101 fails", perPage: 101, wantErr: true, errContains: "per_page"},
		{name: "negative per_page fails", perPage: -1, wantErr: true, errContains: "per_page"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := list_greetings.NewUseCase(&noopRepository{})
			_, err := useCase.Run(context.Background(), list_greetings.DTO{Page: 1, PerPage: tt.perPage})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}
