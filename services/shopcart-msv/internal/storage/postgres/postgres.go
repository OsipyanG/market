package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/OsipyanG/market/services/shopcart-msv/config"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopCartStorage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, config *config.PostgresConfig) (*ShopCartStorage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.User, config.Password, net.JoinHostPort(config.Host, config.Port), config.DB, config.SSLMode)

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	if dbpool.Ping(ctx) != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	return &ShopCartStorage{dbpool}, nil
}

func (s *ShopCartStorage) Close() {
	s.pool.Close()
}
