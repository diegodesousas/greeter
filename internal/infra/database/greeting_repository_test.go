package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/diegodesousas/greeter/internal/domain/greeting"
	"github.com/diegodesousas/greeter/internal/infra/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestConnection(t *testing.T) database.Connection {
	t.Helper()

	viper.Set("DB_HOST", "localhost")
	viper.Set("DB_PORT", "5432")
	viper.Set("DB_USER", "postgres")
	viper.Set("DB_PASSWORD", "postgres")
	viper.Set("DB_NAME", "greeter")
	viper.Set("DB_SSL_MODE", "disable")
	viper.Set("DB_MAX_OPEN_CONN", 10)
	viper.Set("DB_MAX_IDLE_CONN", 5)
	viper.Set("DB_CONN_MAX_IDLE_TIME", 30)
	viper.Set("DB_CONN_MAX_LIFETIME", 60)

	conn, err := database.NewPostgresConnection()
	if err != nil {
		t.Skip("database not available:", err)
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		t.Skip("database not reachable:", err)
	}

	t.Cleanup(func() {
		conn.Close()
	})

	return conn
}

func TestGreetingRepository_Save(t *testing.T) {
	tests := []struct {
		name     string
		greeting greeting.Greeting
		wantErr  bool
	}{
		{
			name: "saves greeting successfully",
			greeting: greeting.Greeting{
				Name:      "Diego",
				Message:   "Hello, Diego!",
				GreetedAt: time.Date(2026, 5, 22, 12, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "saves greeting with long name",
			greeting: greeting.Greeting{
				Name:      "Aaaaabbbbbcccccdddddeeeee",
				Message:   "Hello, Aaaaabbbbbcccccdddddeeeee!",
				GreetedAt: time.Date(2026, 5, 22, 15, 30, 0, 0, time.UTC),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := newTestConnection(t)
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
