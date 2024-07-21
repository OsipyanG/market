package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/OsipyanG/market/services/order-msv/config"
	"github.com/OsipyanG/market/services/order-msv/internal/repository"
	"github.com/OsipyanG/market/services/order-msv/pkg/errwrap"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, cfg config.Postgres) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", cfg.User, cfg.Password, net.JoinHostPort(cfg.Host, cfg.Port), cfg.Name, cfg.SSLMode)

	dbPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, errwrap.Wrap(repository.ErrRepositoryConnection, err)
	}

	if err := dbPool.Ping(ctx); err != nil {
		return nil, errwrap.Wrap(repository.ErrRepositoryConnection, err)
	}

	return dbPool, nil
}
