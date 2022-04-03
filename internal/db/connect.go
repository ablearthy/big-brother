package db

import (
	"big-brother/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

var conn *pgx.Conn

func Connect(cfg config.DbConfig) (err error) {
	conn, err = pgx.Connect(context.Background(), getDbUrl(cfg))
	return
}

func GetConn() *pgx.Conn {
	return conn
}

func getDbUrl(cfg config.DbConfig) string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname, cfg.Sslmode)
}
