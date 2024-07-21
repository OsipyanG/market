package repository

import (
	"context"

	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/google/uuid"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string) error
	GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]model.Order, error)

	GetAllDeliveries(ctx context.Context) ([]model.Delivery, error)
	GetAllPendingOrders(ctx context.Context) ([]model.Order, error)
	CreateDelivery(ctx context.Context, delivery model.Delivery) error
	UpdateDeliveryStatus(ctx context.Context, orderID uuid.UUID, status string) error
}
