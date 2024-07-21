package repository

import (
	"context"

	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"github.com/google/uuid"
)

type WarehouseRepository interface {
	ReserveProducts(ctx context.Context, products []model.ProductQuantity) error
	FreeReservedProducts(ctx context.Context, products []model.ProductQuantity) error
	DeleteReservedProducts(ctx context.Context, products []model.ProductQuantity) error
	GetProductPrices(ctx context.Context, productIDs []uuid.UUID) ([]model.ProductPrice, error)
}

type CatalogRepository interface {
	GetCatalog(ctx context.Context, offset int, limit int) ([]model.Product, error)
}
