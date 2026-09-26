package greet_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type noopClock struct{}

func (noopClock) Now() time.Time { return time.Time{} }

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
		{name: "simple name is valid", input: "Diego", wantErr: false},
		{name: "empty name fails", input: "", wantErr: true, errContains: "name"},
		{name: "50 ascii characters is valid", input: strings.Repeat("a", 50), wantErr: false},
		{name: "51 ascii characters fails", input: strings.Repeat("a", 51), wantErr: true, errContains: "at most 50 characters"},
		{name: "50 accented characters is valid", input: strings.Repeat("ã", 50), wantErr: false},
		{name: "51 accented characters fails", input: strings.Repeat("ã", 51), wantErr: true, errContains: "at most 50 characters"},
		{name: "50 multi-byte characters is valid", input: strings.Repeat("日", 50), wantErr: false},
		{name: "accented name under limit is valid", input: "João Conceição Araújo de Magalhães Gonçalves", wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := greet.NewUseCase(noopClock{}, &noopRepository{})
			_, err := useCase.Run(context.Background(), greet.DTO{Name: tt.input})

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			require.NoError(t, err)
		})
	}
}
