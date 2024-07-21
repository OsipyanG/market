package postgres

import (
	"context"
	"fmt"

	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CatalogRepository struct {
	pool *pgxpool.Pool
}

func NewCatalogRepository(pool *pgxpool.Pool) *CatalogRepository {
	return &CatalogRepository{
		pool: pool,
	}
}

func (s *CatalogRepository) GetCatalog(ctx context.Context, offset int, limit int) ([]model.Product, error) {
	query := `
		SELECT product_id, name, description, available, quantity, price
		FROM products
		ORDER BY product_id
		LIMIT $1 OFFSET $2
	`

	rows, err := s.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrGetCatalog, err)
	}
	defer rows.Close()

	catalog := make([]model.Product, 0)

	for rows.Next() {
		var product model.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Available,
			&product.Quantity,
			&product.Price,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", repository.ErrGetCatalog, err)
		}

		catalog = append(catalog, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrGetCatalog, err)
	}

	return catalog, nil
}
