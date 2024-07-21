package postgres

import (
	"context"

	"github.com/OsipyanG/market/services/shopcart-msv/internal/model"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/storage"
	"github.com/OsipyanG/market/services/shopcart-msv/pkg/errwrap"
	"github.com/google/uuid"
)

func (s *ShopCartStorage) AddProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error {
	_, err := s.pool.Exec(ctx, "CALL add_product($1, $2, $3)", userID, product.ID, product.Quantity)
	if err != nil {
		return errwrap.Wrap(storage.ErrAddProduct, err)
	}

	return nil
}

func (s *ShopCartStorage) DeleteProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error {
	var rowsAffected int

	err := s.pool.QueryRow(ctx, "SELECT delete_product($1, $2, $3)", userID,
		product.ID, product.Quantity).Scan(&rowsAffected)
	if err != nil {
		return errwrap.Wrap(storage.ErrAddProduct, err)
	}

	if rowsAffected == 0 {
		return errwrap.Wrap(storage.ErrAddProduct, storage.ErrNoRowsAffected)
	}

	return nil
}

func (s *ShopCartStorage) GetProducts(ctx context.Context, userID uuid.UUID) ([]*model.Product, error) {
	rows, err := s.pool.Query(ctx, "SELECT product_id, quantity FROM Products WHERE user_id=$1", userID)
	if err != nil {
		return nil, errwrap.Wrap(storage.ErrGetProducts, err)
	}
	defer rows.Close()

	products := []*model.Product{}

	for rows.Next() {
		product := &model.Product{}

		err := rows.Scan(&product.ID, &product.Quantity)
		if err != nil {
			return nil, errwrap.Wrap(storage.ErrGetProducts, err)
		}

		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, errwrap.Wrap(storage.ErrGetProducts, err)
	}

	return products, nil
}

func (s *ShopCartStorage) Clear(ctx context.Context, userID uuid.UUID) error {
	cmdTag, err := s.pool.Exec(ctx, "DELETE FROM Products WHERE user_id=$1", userID)
	if err != nil {
		return errwrap.Wrap(storage.ErrClear, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(storage.ErrClear, storage.ErrNoRowsAffected)
	}

	return nil
}
