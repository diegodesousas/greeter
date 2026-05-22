package database_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	devkitsql "github.com/diegodesousas/go-devkit/pkg/database/sql"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/diegodesousas/greeter/internal/infra/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnection struct {
	execErr error
}

func (m *mockConnection) Exec(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, m.execErr
}

func (m *mockConnection) Get(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

func (m *mockConnection) Select(_ context.Context, _ interface{}, _ string, _ ...interface{}) error {
	return nil
}

func (m *mockConnection) Begin(_ context.Context) (devkitsql.Transaction, error) {
	return nil, nil
}

func (m *mockConnection) TransactionContext(_ context.Context, _ func(context.Context) error) error {
	return nil
}

func (m *mockConnection) Ping() error  { return nil }
func (m *mockConnection) Close() error { return nil }

var fixedTime = time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC)

func TestGreetingRepository_Save(t *testing.T) {
	tests := []struct {
		name     string
		greeting greeting.Greeting
		execErr  error
		wantErr  bool
	}{
		{
			name: "saves greeting successfully",
			greeting: greeting.Greeting{
				Name:      "Diego",
				Message:   "Hello, Diego!",
				GreetedAt: fixedTime,
			},
		},
		{
			name: "returns error when exec fails",
			greeting: greeting.Greeting{
				Name:      "Diego",
				Message:   "Hello, Diego!",
				GreetedAt: fixedTime,
			},
			execErr: errors.New("db error"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConnection{execErr: tt.execErr}
			repo := database.NewGreetingRepository(conn)

			err := repo.Save(context.Background(), tt.greeting)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
