package shutdown_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/infra/shutdown"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGraceful(t *testing.T) {
	errServer := errors.New("server shutdown failed")
	errDatabase := errors.New("database close failed")

	tests := []struct {
		name       string
		stepErrs   []error
		wantErrIs  []error
		wantCalled []int
	}{
		{
			name:       "runs every step in order",
			stepErrs:   []error{nil, nil, nil},
			wantCalled: []int{0, 1, 2},
		},
		{
			name:       "keeps running remaining steps after a failure",
			stepErrs:   []error{errServer, nil},
			wantErrIs:  []error{errServer},
			wantCalled: []int{0, 1},
		},
		{
			name:       "joins errors from every failing step",
			stepErrs:   []error{errServer, errDatabase},
			wantErrIs:  []error{errServer, errDatabase},
			wantCalled: []int{0, 1},
		},
		{
			name:       "no steps is a no-op",
			stepErrs:   nil,
			wantCalled: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var called []int
			steps := make([]shutdown.Step, len(tt.stepErrs))
			for i, stepErr := range tt.stepErrs {
				steps[i] = func(context.Context) error {
					called = append(called, i)
					return stepErr
				}
			}

			err := shutdown.Graceful(context.Background(), time.Second, steps...)

			assert.Equal(t, tt.wantCalled, called)
			if tt.wantErrIs == nil {
				require.NoError(t, err)
				return
			}
			for _, target := range tt.wantErrIs {
				assert.ErrorIs(t, err, target)
			}
		})
	}
}

func TestGraceful_StepsShareTheTimeout(t *testing.T) {
	closed := false

	start := time.Now()
	err := shutdown.Graceful(context.Background(), 50*time.Millisecond,
		func(ctx context.Context) error {
			_, hasDeadline := ctx.Deadline()
			require.True(t, hasDeadline, "step context must have a deadline")

			<-ctx.Done() // simulates a server that never finishes draining
			return ctx.Err()
		},
		func(context.Context) error {
			closed = true
			return nil
		},
	)

	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.True(t, closed, "later steps must still run after a timeout")
	assert.Less(t, time.Since(start), time.Second, "shutdown must not hang past the timeout")
}
