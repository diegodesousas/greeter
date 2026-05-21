package clock_test

import (
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/infra/clock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, clock.New())
}

func TestNow(t *testing.T) {
	before := time.Now().UTC()
	got := clock.New().Now()
	after := time.Now().UTC()

	assert.Equal(t, time.UTC, got.Location())
	assert.False(t, got.Before(before), "Now() should not be before the call")
	assert.False(t, got.After(after), "Now() should not be after the call")
}