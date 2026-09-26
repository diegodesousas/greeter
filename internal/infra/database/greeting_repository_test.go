package database_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/diegodesousas/greeter/internal/infra/database"
	"github.com/diegodesousas/greeter/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var fixedTime = time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC)

// setTotal fills the COUNT(*) destination passed to Get.
func setTotal(total int) func(context.Context, interface{}, string, ...interface{}) {
	return func(_ context.Context, dest interface{}, _ string, _ ...interface{}) {
		*dest.(*int) = total
	}
}

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
			conn := mocks.NewMockConnection(t)
			conn.EXPECT().Exec(mock.Anything, mock.Anything, tt.greeting.Name, tt.greeting.GreetedAt).Return(nil, tt.execErr)
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
			name:      "returns greetings and total",
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
			getResult: 5,
			selectErr: errors.New("db error"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := mocks.NewMockConnection(t)
			conn.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything).Run(setTotal(tt.getResult)).Return(tt.getErr)
			if tt.getErr == nil {
				conn.EXPECT().Select(mock.Anything, mock.Anything, mock.Anything, tt.perPage, (tt.page-1)*tt.perPage).Return(tt.selectErr)
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
			conn := mocks.NewMockConnection(t)
			conn.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything, "%joao%").Run(setTotal(tt.getResult)).Return(tt.getErr)
			if tt.getErr == nil {
				conn.EXPECT().Select(mock.Anything, mock.Anything, mock.Anything, "%joao%", tt.perPage, (tt.page-1)*tt.perPage).Return(tt.selectErr)
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

func TestGreetingRepository_SearchEscapesLikeWildcards(t *testing.T) {
	tests := []struct {
		name        string
		search      string
		wantPattern string
	}{
		{name: "plain term is wrapped in wildcards", search: "joao", wantPattern: `%joao%`},
		{name: "percent is matched literally", search: "100%", wantPattern: `%100\%%`},
		{name: "underscore is matched literally", search: "a_b", wantPattern: `%a\_b%`},
		{name: "backslash is matched literally", search: `a\b`, wantPattern: `%a\\b%`},
		{name: "lone percent does not match everything", search: "%", wantPattern: `%\%%`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := mocks.NewMockConnection(t)
			conn.EXPECT().Get(mock.Anything, mock.Anything, mock.Anything, tt.wantPattern).Run(setTotal(0)).Return(nil)
			conn.EXPECT().Select(mock.Anything, mock.Anything, mock.Anything, tt.wantPattern, 10, 0).Return(nil)
			repo := database.NewGreetingRepository(conn)

			_, _, err := repo.Search(context.Background(), tt.search, 1, 10)

			require.NoError(t, err)
		})
	}
}
