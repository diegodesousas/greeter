package search_greetings_test

import (
	"context"
	"strings"
	"testing"

	search_greetings "github.com/diegodesousas/greeter/internal/application/search_greetings"
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

func TestValidator_Name(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errContains string
	}{
		{name: "simple term is valid", input: "joão", wantErr: false},
		{name: "empty term fails", input: "", wantErr: true, errContains: "name"},
		{name: "50 characters is valid", input: strings.Repeat("a", 50), wantErr: false},
		{name: "51 characters fails", input: strings.Repeat("a", 51), wantErr: true, errContains: "at most 50 characters"},
		{name: "50 accented characters is valid", input: strings.Repeat("ã", 50), wantErr: false},
		{name: "51 accented characters fails", input: strings.Repeat("ã", 51), wantErr: true, errContains: "at most 50 characters"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := search_greetings.NewUseCase(&noopRepository{})
			_, err := useCase.Run(context.Background(), search_greetings.DTO{Name: tt.input, Page: 1, PerPage: 10})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}
