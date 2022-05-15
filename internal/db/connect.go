package db

import (
	"big-brother/internal/config"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

var conn *pgxpool.Pool

func Connect(cfg config.DbConfig) (err error) {
	conn, err = pgxpool.Connect(context.Background(), getDbUrl(cfg))
	return
}

func GetConn() *pgxpool.Pool {
	return conn
}

func getDbUrl(cfg config.DbConfig) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname, cfg.Sslmode)
}
