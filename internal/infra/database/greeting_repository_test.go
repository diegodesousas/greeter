package database_test

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	devkitsql "github.com/diegodesousas/go-devkit/pkg/database/sql"
	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/diegodesousas/greeter/internal/infra/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockConnection struct {
	execErr   error
	getErr    error
	selectErr error
	getResult interface{}
	selectResult interface{}
}

func (m *mockConnection) Exec(_ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
	return nil, m.execErr
}

func (m *mockConnection) Get(_ context.Context, dest interface{}, _ string, _ ...interface{}) error {
	if m.getErr != nil {
		return m.getErr
	}
	if m.getResult != nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(m.getResult))
	}
	return nil
}

func (m *mockConnection) Select(_ context.Context, dest interface{}, _ string, _ ...interface{}) error {
	if m.selectErr != nil {
		return m.selectErr
	}
	if m.selectResult != nil {
		reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(m.selectResult))
	}
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

func TestGreetingRepository_List(t *testing.T) {
	tests := []struct {
		name          string
		page          int
		perPage       int
		getResult     int
		getErr        error
		selectErr     error
		wantTotal     int
		wantLen       int
		wantFirstName string
		wantErr       bool
	}{
		{
			name:          "returns greetings and total",
			page:          1,
			perPage:       10,
			getResult:     2,
			wantTotal:     2,
			wantLen:       0,
		},
		{
			name:      "count query error is propagated",
			page:      1,
			perPage:   10,
			getErr:    errors.New("db error"),
			wantErr:   true,
		},
		{
			name:      "select query error is propagated",
			page:      1,
			perPage:   10,
			getResult: 5,
			selectErr: errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConnection{
				getResult: tt.getResult,
				getErr:    tt.getErr,
				selectErr: tt.selectErr,
			}
			repo := database.NewGreetingRepository(conn)

			greetings, total, err := repo.List(context.Background(), tt.page, tt.perPage)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.Len(t, greetings, tt.wantLen)
		})
	}
}

func TestGreetingRepository_Search(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		perPage   int
		getResult int
		getErr    error
		selectErr error
		wantTotal int
		wantLen   int
		wantErr   bool
	}{
		{
			name:      "returns matching greetings and total",
			page:      1,
			perPage:   10,
			getResult: 2,
			wantTotal: 2,
			wantLen:   0,
		},
		{
			name:    "count query error is propagated",
			page:    1,
			perPage: 10,
			getErr:  errors.New("db error"),
			wantErr: true,
		},
		{
			name:      "select query error is propagated",
			page:      1,
			perPage:   10,
			getResult: 1,
			selectErr: errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := &mockConnection{
				getResult: tt.getResult,
				getErr:    tt.getErr,
				selectErr: tt.selectErr,
			}
			repo := database.NewGreetingRepository(conn)

			greetings, total, err := repo.Search(context.Background(), "joao", tt.page, tt.perPage)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, total)
			assert.Len(t, greetings, tt.wantLen)
		})
	}
}
