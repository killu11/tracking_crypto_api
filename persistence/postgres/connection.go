package postgres

import (
	"context"
	"crypto_api/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"time"
)

func NewPostgresConnection(cfg *config.PostgresConfig) (*sql.DB, error) {
	conn, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("create postgres connection: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err = conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return conn, nil
}
