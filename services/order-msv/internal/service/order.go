package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/OsipyanG/market/services/order-msv/internal/repository"
	"github.com/OsipyanG/market/services/order-msv/internal/transport/shopcart"
	"github.com/OsipyanG/market/services/order-msv/internal/transport/warehouse"
	"github.com/OsipyanG/market/services/order-msv/pkg/errwrap"
	"github.com/google/uuid"
)

type OrderService struct {
	repository repository.OrderRepository
	shopcart   shopcart.Client
	warehouse  warehouse.Client
}

func NewOrderService(repo repository.OrderRepository, shopcartClient shopcart.Client, warehouseClient warehouse.Client) *OrderService {
	return &OrderService{repository: repo, shopcart: shopcartClient, warehouse: warehouseClient}
}

func (os *OrderService) CreateOrder(ctx context.Context, customerID uuid.UUID, address string) (uuid.UUID, error) {
	orderItems, err := os.shopcart.GetProducts(ctx, customerID)
	if err != nil {
		return uuid.Nil, errwrap.Wrap(ErrCreateOrder, err)
	}

	if len(orderItems) == 0 {
		return uuid.Nil, errwrap.Wrap(ErrCreateOrder, ErrEmptyShopcart)
	}

	err = os.warehouse.ReserveProducts(ctx, orderItems)
	if err != nil {
		return uuid.Nil, errwrap.Wrap(ErrCreateOrder, err)
	}

	ids := make([]uuid.UUID, 0, len(orderItems))
	for _, item := range orderItems {
		ids = append(ids, item.ProductID)
	}

	prices, err := os.warehouse.GetProductsPrices(ctx, ids)

	for i := range orderItems {
		for j := range prices {
			if orderItems[i].ProductID == prices[j].ProductID {
				orderItems[i].Price = prices[j].Price
			}
		}
	}

	order := model.Order{
		CustomerID: customerID,
		Items:      orderItems,
		Address:    address,
	}

	orderID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, errwrap.Wrap(ErrCreateOrder, err)
	}

	order.ID = orderID

	err = os.repository.CreateOrder(ctx, order)
	if err != nil {
		err := os.warehouse.FreeReservedProducts(ctx, orderItems)
		if err != nil {
			slog.Error("failed to free reserved products: ", "err", err)
		}

		return uuid.Nil, errwrap.Wrap(ErrCreateOrder, err)
	}

	defer func() {
		err = os.shopcart.Clear(ctx, customerID)
		if err != nil {
			slog.Error("failed to clear shopcart: ", "err", err)
		}
	}()

	return order.ID, nil
}

func (os *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	order, err := os.repository.GetOrder(ctx, orderID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetOrder, err)
	}

	return order, nil
}

func (os *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string, jwtClaims model.JwtClaims) error {
	order, err := os.repository.GetOrder(ctx, orderID)
	if err != nil {
		return errwrap.Wrap(ErrUpdateOrder, err)
	}

	switch status {
	case PendingStatus:
		if order.Status != AssemblyStatus {
			return fmt.Errorf("%w:  %w %s", ErrUpdateOrder, ErrPreCondition, "order status is not Assembly")
		}

		if jwtClaims.AccessLevel != WarehouseWorkerLevel {
			return errwrap.Wrap(ErrUpdateOrder, ErrAccessDenied)
		}

	case DeliveryStatus:
		if order.Status != PendingStatus {
			return fmt.Errorf("%w:  %w %s", ErrUpdateOrder, ErrPreCondition, "order status is not Pending")
		}

		if jwtClaims.AccessLevel != CourierLevel {
			return errwrap.Wrap(ErrUpdateOrder, ErrAccessDenied)
		}

		err := os.repository.CreateDelivery(ctx, model.Delivery{
			OrderID:   orderID,
			CourierID: jwtClaims.UserID,
			Status:    DeliveryStatus,
		})
		if err != nil {
			return errwrap.Wrap(ErrUpdateOrder, err)
		}

	case CompletedStatus:
		if order.Status != DeliveryStatus {
			return fmt.Errorf("%w:  %w %s", ErrUpdateOrder, ErrPreCondition, "order status is not Delivery")
		}

		if jwtClaims.AccessLevel != CourierLevel {
			return errwrap.Wrap(ErrUpdateOrder, ErrAccessDenied)
		}

		err := os.warehouse.DeleteReservedProducts(ctx, order.Items)
		if err != nil {
			return errwrap.Wrap(ErrUpdateOrder, err)
		}

		err = os.repository.UpdateDeliveryStatus(ctx, orderID, status)
		if err != nil {
			return errwrap.Wrap(ErrUpdateOrder, err)
		}
	case CancelledStatus:
		// TODO: добавить логику, что бы пользователь мог отменить свой заказ
		// TODO: добавить проверку, что пользователь отменяет свой заказ
		// TODO: Добавить здесь source что бы понять откуда пришел запрос, от пользователя как клиента или курьера
		if !(jwtClaims.AccessLevel == CourierLevel && order.Status == DeliveryStatus) {
			return errwrap.Wrap(ErrUpdateOrder, ErrAccessDenied)
		}

		err = os.warehouse.FreeReservedProducts(ctx, order.Items)
		if err != nil {
			return errwrap.Wrap(ErrUpdateOrder, err)
		}

		err = os.repository.UpdateDeliveryStatus(ctx, orderID, status)
		if err != nil {
			return errwrap.Wrap(ErrUpdateOrder, err)
		}

	default:
		return errwrap.Wrap(ErrUpdateOrder, ErrInvalidStatus)
	}

	err = os.repository.UpdateOrderStatus(ctx, orderID, status)
	if err != nil {
		return errwrap.Wrap(ErrUpdateOrder, err)
	}

	return nil
}

func (os *OrderService) GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]model.Order, error) {
	orders, err := os.repository.GetAllOrders(ctx, customerID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllOrders, err)
	}

	return orders, nil
}

func (os *OrderService) GetAllDeliveries(ctx context.Context, jwtClaims model.JwtClaims) ([]model.Delivery, error) {
	if jwtClaims.AccessLevel != CourierLevel && jwtClaims.AccessLevel != AdminLevel {
		return nil, errwrap.Wrap(ErrGetAllDeliveries, ErrAccessDenied)
	}

	deliveries, err := os.repository.GetAllDeliveries(ctx)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllDeliveries, err)
	}

	return deliveries, nil
}

func (os *OrderService) GetAllPendingOrders(ctx context.Context, jwtClaims model.JwtClaims) ([]model.Order, error) {
	if jwtClaims.AccessLevel != CourierLevel && jwtClaims.AccessLevel != AdminLevel {
		return nil, errwrap.Wrap(ErrGetAllPendingOrders, ErrAccessDenied)
	}

	deliveries, err := os.repository.GetAllPendingOrders(ctx)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllPendingOrders, err)
	}

	return deliveries, nil
}
