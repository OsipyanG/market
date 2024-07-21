package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WarehouseRepository struct {
	pool *pgxpool.Pool
}

func NewWarehouseRepository(pool *pgxpool.Pool) *WarehouseRepository {
	return &WarehouseRepository{
		pool: pool,
	}
}

func (s *WarehouseRepository) ReserveProducts(ctx context.Context, products []model.ProductQuantity) error {
	checkQuery := `SELECT available FROM products WHERE product_id = $1`
	updateQuery := `UPDATE products SET available = available - $1 WHERE product_id = $2`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repository.ErrReserveProducts, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("Error rolling back transaction", "error:", err)
		}
	}()

	for _, product := range products {
		var available int64

		err := tx.QueryRow(ctx, checkQuery, product.ID).Scan(&available)
		if err != nil {
			return fmt.Errorf("%w: %w", repository.ErrReserveProducts, err)
		}

		if available < product.Quantity {
			return fmt.Errorf("%w: %w %s", repository.ErrReserveProducts,
				repository.ErrNotEnoughStock, product.ID.String())
		}

		cmdTag, err := tx.Exec(ctx, updateQuery, product.Quantity, product.ID)
		if err != nil {
			return fmt.Errorf("%w: %w", repository.ErrReserveProducts, err)
		}

		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("%w: %w", repository.ErrReserveProducts, repository.ErrNoRowsAffected)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repository.ErrReserveProducts, err)
	}

	return nil
}

func (s *WarehouseRepository) FreeReservedProducts(ctx context.Context, products []model.ProductQuantity) error {
	checkQuery := `SELECT available, quantity FROM products WHERE product_id = $1`
	updateQuery := `UPDATE Products SET available = available + $1 WHERE product_id = $2`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repository.ErrFreeReservedProducts, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("Error rolling back transaction", "error:", err)
		}
	}()

	for _, product := range products {
		var available, quantity int64

		err := tx.QueryRow(ctx, checkQuery, product.ID).Scan(&available, &quantity)
		if err != nil {
			return fmt.Errorf("%w: %w", repository.ErrFreeReservedProducts, err)
		}

		if product.Quantity > quantity-available {
			return fmt.Errorf("%w: %w %s", repository.ErrFreeReservedProducts,
				repository.ErrNotSoManyReserved, product.ID.String())
		}

		cmdTag, err := tx.Exec(ctx, updateQuery, product.Quantity, product.ID)
		if err != nil {
			return fmt.Errorf("%w: %w", repository.ErrFreeReservedProducts, err)
		}

		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("%w: %w", repository.ErrFreeReservedProducts, repository.ErrNoRowsAffected)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repository.ErrFreeReservedProducts, err)
	}

	return nil
}

func (s *WarehouseRepository) DeleteReservedProducts(ctx context.Context, products []model.ProductQuantity) error {
	checkQuery := `SELECT available, quantity FROM products WHERE product_id = $1`
	updateQuery := `UPDATE Products SET quantity = quantity - $1 WHERE product_id = $2`

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repository.ErrDeleteReservedProducts, err)
	}

	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("Error rolling back transaction", "error:", err)
		}
	}()

	for _, product := range products {
		var available, quantity int64

		err := tx.QueryRow(ctx, checkQuery, product.ID).Scan(&available, &quantity)
		if err != nil {
			return fmt.Errorf("%w: %w", repository.ErrDeleteReservedProducts, err)
		}

		if product.Quantity > quantity-available {
			return fmt.Errorf("%w: %w %s", repository.ErrDeleteReservedProducts,
				repository.ErrNotSoManyReserved, product.ID.String())
		}

		cmdTag, err := tx.Exec(ctx, updateQuery, product.Quantity, product.ID)
		if err != nil {
			return fmt.Errorf("%w: %w", repository.ErrDeleteReservedProducts, err)
		}

		if cmdTag.RowsAffected() == 0 {
			return fmt.Errorf("%w: %w", repository.ErrFreeReservedProducts, repository.ErrNoRowsAffected)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", repository.ErrDeleteReservedProducts, err)
	}

	return nil
}

func (s *WarehouseRepository) GetProductPrices(ctx context.Context, productIDs []uuid.UUID) ([]model.ProductPrice, error) {
	query := `
		SELECT product_id, price
		FROM products WHERE product_id = ANY($1)
		`

	rows, err := s.pool.Query(ctx, query, productIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrGetProductPrices, err)
	}
	defer rows.Close()

	productsPrices := make([]model.ProductPrice, 0)

	for rows.Next() {
		var ProductPrice model.ProductPrice

		err := rows.Scan(&ProductPrice.ID, &ProductPrice.Price)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", repository.ErrGetProductPrices, err)
		}

		productsPrices = append(productsPrices, ProductPrice)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("%w: %w", repository.ErrGetProductPrices, err)
	}

	return productsPrices, nil
}
