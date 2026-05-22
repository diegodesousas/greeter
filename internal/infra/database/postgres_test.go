package database_test

import (
	"testing"

	devkitsql "github.com/diegodesousas/go-devkit/pkg/database/sql"
	"github.com/diegodesousas/greeter/internal/infra/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPostgresConnection(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		wantErr error
	}{
		{
			name: "returns error when host is unreachable",
			setup: func() {
				viper.Set("DB_HOST", "localhost")
				viper.Set("DB_PORT", "9999")
				viper.Set("DB_USER", "postgres")
				viper.Set("DB_PASSWORD", "postgres")
				viper.Set("DB_NAME", "greeter")
				viper.Set("DB_SSL_MODE", "disable")
				viper.Set("DB_MAX_OPEN_CONN", 10)
				viper.Set("DB_MAX_IDLE_CONN", 5)
				viper.Set("DB_CONN_MAX_IDLE_TIME", 30)
				viper.Set("DB_CONN_MAX_LIFETIME", 60)
			},
			wantErr: devkitsql.ErrConn,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			conn, err := database.NewPostgresConnection()

			require.Error(t, err)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Nil(t, conn)
		})
	}
}
