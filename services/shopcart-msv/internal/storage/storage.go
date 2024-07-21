package storage

import (
	"context"
	"errors"

	"github.com/OsipyanG/market/services/shopcart-msv/internal/model"
	"github.com/google/uuid"
)

var (
	ErrAddProduct    = errors.New("storage: can't add product to the shopcart")
	ErrDeleteProduct = errors.New("storage: can't delete product from the shopcart")
	ErrGetProducts   = errors.New("storage: cant't get products from the shopcart")
	ErrClear         = errors.New("storage: can't clear shopcart")

	ErrStorageConnection = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected    = errors.New("no rows affected")
	ErrNotFound          = errors.New("not found")
)

type ShopCartStorage interface {
	AddProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error
	DeleteProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error
	GetProducts(ctx context.Context, userID uuid.UUID) ([]*model.Product, error)
	Clear(ctx context.Context, userID uuid.UUID) error
	Close()
}
