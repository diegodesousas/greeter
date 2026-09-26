package greet_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/application/greet"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

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
			clock := mocks.NewMockClock(t)
			clock.EXPECT().Now().Return(time.Time{}).Maybe()
			repo := mocks.NewMockGreetingRepository(t)
			repo.EXPECT().Save(mock.Anything, mock.Anything).Return(nil).Maybe()

			useCase := greet.NewUseCase(clock, repo)
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
