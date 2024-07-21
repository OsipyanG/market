package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/OsipyanG/market/services/warehouse-msv/config"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewClient(ctx context.Context, cfg *config.Postgres) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		cfg.User, cfg.Password, net.JoinHostPort(cfg.Host, cfg.Port), cfg.DB, cfg.SSLMode)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrRepositoryConnection, err)
	}

	if pool.Ping(ctx) != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrRepositoryConnection, err)
	}

	return pool, nil
}
