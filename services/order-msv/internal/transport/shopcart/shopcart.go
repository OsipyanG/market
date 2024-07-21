package shopcart

import (
	"context"

	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/google/uuid"
)

type Client interface {
	GetProducts(ctx context.Context, userID uuid.UUID) ([]model.OrderItem, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}
