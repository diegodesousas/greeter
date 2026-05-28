package greeting_test

import (
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/stretchr/testify/assert"
)

var fixedTime = time.Date(2026, 5, 8, 12, 0, 0, 0, time.UTC)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		inputName   string
		wantMessage string
	}{
		{
			name:        "simple name",
			inputName:   "Diego",
			wantMessage: "Hello, Diego!",
		},
		{
			name:        "name with spaces",
			inputName:   "John Doe",
			wantMessage: "Hello, John Doe!",
		},
		{
			name:        "name with special characters",
			inputName:   "O'Brien",
			wantMessage: "Hello, O'Brien!",
		},
		{
			name:        "empty name",
			inputName:   "",
			wantMessage: "Hello, !",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := greeting.New(tt.inputName, fixedTime)

			assert.Empty(t, g.ID)
			assert.Equal(t, tt.inputName, g.Name)
			assert.Equal(t, tt.wantMessage, g.Message)
			assert.Equal(t, fixedTime, g.GreetedAt)
		})
	}
}