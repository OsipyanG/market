package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/OsipyanG/market/services/auth-msv/config"
	"github.com/OsipyanG/market/services/auth-msv/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, config *config.PostgresConfig) (*Storage, error) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.User, config.Password, net.JoinHostPort(config.Host, config.Port), config.DB, config.SSLMode)

	dbpool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	if dbpool.Ping(ctx) != nil {
		return nil, fmt.Errorf("%w: %w", storage.ErrStorageConnection, err)
	}

	return &Storage{dbpool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}
