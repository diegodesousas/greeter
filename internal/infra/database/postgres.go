package database

import (
	"time"

	devkitsql "github.com/diegodesousas/go-devkit/pkg/database/sql"
	"github.com/spf13/viper"
)

func NewPostgresConnection() (devkitsql.Connection, error) {
	cfg := devkitsql.Config{
		Host:            viper.GetString("DB_HOST"),
		Port:            viper.GetString("DB_PORT"),
		User:            viper.GetString("DB_USER"),
		Password:        viper.GetString("DB_PASSWORD"),
		Database:        viper.GetString("DB_NAME"),
		SSLMode:         viper.GetString("DB_SSL_MODE"),
		MaxOpenConn:     viper.GetInt("DB_MAX_OPEN_CONN"),
		MaxIdleConn:     viper.GetInt("DB_MAX_IDLE_CONN"),
		ConnMaxIdleTime: viper.GetDuration("DB_CONN_MAX_IDLE_TIME") * time.Second,
		ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME") * time.Second,
	}

	return devkitsql.New(cfg)
}
