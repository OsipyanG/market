package controller

import (
	"context"
	"errors"

	orderpb "github.com/OsipyanG/market/protos/order"
	"github.com/OsipyanG/market/services/order-msv/internal/converter"
	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/OsipyanG/market/services/order-msv/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderService interface {
	CreateOrder(ctx context.Context, customerID uuid.UUID, address string) (uuid.UUID, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, status string, jwtClaims model.JwtClaims) error
	GetAllOrders(ctx context.Context, customerID uuid.UUID) ([]model.Order, error)

	GetAllDeliveries(ctx context.Context, jwtClaims model.JwtClaims) ([]model.Delivery, error)
	GetAllPendingOrders(ctx context.Context, jwtClaims model.JwtClaims) ([]model.Order, error)
}

type OrderController struct {
	orderpb.UnimplementedOrderServiceServer
	s OrderService
}

func NewOrderController(s OrderService) *OrderController {
	return &OrderController{s: s}
}

// FIXME: спрятать ошибки внутренностей.
// FIXME: добавить логирование.
// FIXME: добавить валидацию всех входных данных.
// FIXME: добавить правильную обработку ошибок бизнес логики.
// FIXME: выкидывать ошибку Internal только в случае незапланированных ошибок.
// FIXME: удалить JwtClaims модель из бизнес логики, т.к. она используется только в одном месте, в данном случае легче просто добавить параметр в функцию.
func (or *OrderController) CreateOrder(ctx context.Context, req *orderpb.CreateOrderRequest) (*orderpb.CreateOrderResponse, error) {
	customerID, err := uuid.Parse(req.GetJwtClaims().GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid customer ID: %v", err)
	}

	orderID, err := or.s.CreateOrder(ctx, customerID, req.GetAddress())
	if err != nil {
		if errors.Is(err, service.ErrEmptyShopcart) {
			return nil, status.Errorf(codes.FailedPrecondition, "shopcart is empty")
		}

		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return &orderpb.CreateOrderResponse{OrderId: orderID.String()}, nil
}

func (or *OrderController) GetOrder(ctx context.Context, req *orderpb.GetOrderRequest) (*orderpb.GetOrderResponse, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	modelOrder, err := or.s.GetOrder(ctx, orderID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get order: %v", err)
	}

	protoOrder := converter.GetProtoOrder(modelOrder)

	return &orderpb.GetOrderResponse{Order: protoOrder}, nil
}

func (or *OrderController) UpdateOrderStatus(ctx context.Context, req *orderpb.UpdateOrderStatusRequest) (*emptypb.Empty, error) {
	orderID, err := uuid.Parse(req.GetOrderId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid order ID: %v", err)
	}

	orderStatus := req.GetStatus()

	userID, err := uuid.Parse(req.GetJwtClaims().GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	jwtClaims := model.JwtClaims{
		UserID:      userID,
		AccessLevel: int(req.GetJwtClaims().GetAccessLevel()),
	}

	err = or.s.UpdateOrderStatus(ctx, orderID, orderStatus, jwtClaims)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		} else if errors.Is(err, service.ErrInvalidStatus) {
			return nil, status.Errorf(codes.InvalidArgument, "invalid order status")
		} else if errors.Is(err, service.ErrPreCondition) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err) // FIXME спрятать ошибку внутренностей
	}

	return &emptypb.Empty{}, nil
}

func (or *OrderController) GetAllOrders(ctx context.Context, req *orderpb.GetAllOrdersRequest) (*orderpb.GetAllOrdersResponse, error) {
	customerID, err := uuid.Parse(req.GetJwtClaims().GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid customer ID: %v", err)
	}

	modelOrders, err := or.s.GetAllOrders(ctx, customerID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all orders")
	}

	protoOrders := make([]*orderpb.Order, 0, len(modelOrders))

	for _, modelOrder := range modelOrders {
		protoOrder := converter.GetProtoOrder(&modelOrder)
		protoOrders = append(protoOrders, protoOrder)
	}

	return &orderpb.GetAllOrdersResponse{Orders: protoOrders}, nil
}

func (or *OrderController) GetAllDeliveries(ctx context.Context, req *orderpb.GetAllDeliveriesRequest) (*orderpb.GetAllDeliveriesResponse, error) {
	userID, err := uuid.Parse(req.GetJwtClaims().GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	jwtClaims := model.JwtClaims{
		UserID:      userID,
		AccessLevel: int(req.GetJwtClaims().GetAccessLevel()),
	}

	modelDeliveries, err := or.s.GetAllDeliveries(ctx, jwtClaims)
	if err != nil {
		if errors.Is(err, service.ErrAccessDenied) {
			return nil, status.Errorf(codes.PermissionDenied, "access denied")
		}

		return nil, status.Errorf(codes.Internal, "failed to get all deliveries: %v", err)
	}

	protoDeliveries := make([]*orderpb.Delivery, 0, len(modelDeliveries))

	for _, modelDelivery := range modelDeliveries {
		protoDelivery := converter.GetProtoDelivery(modelDelivery)
		protoDeliveries = append(protoDeliveries, protoDelivery)
	}

	return &orderpb.GetAllDeliveriesResponse{Deliveries: protoDeliveries}, nil
}

func (or *OrderController) GetAllPendingOrders(ctx context.Context, req *orderpb.GetAllPendingOrdersRequest) (*orderpb.GetAllPendingOrdersResponse, error) {
	userID, err := uuid.Parse(req.GetJwtClaims().GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID: %v", err)
	}

	jwtClaims := model.JwtClaims{
		UserID:      userID,
		AccessLevel: int(req.GetJwtClaims().GetAccessLevel()),
	}

	modelOrders, err := or.s.GetAllPendingOrders(ctx, jwtClaims)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all pending orders: %v", err)
	}

	protoOrders := make([]*orderpb.Order, 0, len(modelOrders))

	for _, modelOrder := range modelOrders {
		protoOrder := converter.GetProtoOrder(&modelOrder)
		protoOrders = append(protoOrders, protoOrder)
	}

	return &orderpb.GetAllPendingOrdersResponse{Orders: protoOrders}, nil
}
