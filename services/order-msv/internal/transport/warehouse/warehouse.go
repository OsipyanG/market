package warehouse

import (
	"context"

	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/google/uuid"
)

type Client interface {
	ReserveProducts(ctx context.Context, products []model.OrderItem) error
	FreeReservedProducts(ctx context.Context, products []model.OrderItem) error
	DeleteReservedProducts(ctx context.Context, products []model.OrderItem) error
	GetProductsPrices(ctx context.Context, ids []uuid.UUID) ([]model.OrderItem, error)
}
