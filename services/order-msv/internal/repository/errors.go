package repository

import "errors"

var (
	ErrCreateOrder       = errors.New("orderRepository: can't create order")
	ErrGetOrder          = errors.New("orderRepository: can't get info about order")
	ErrGetAllOrders      = errors.New("orderRepository: can't get info about all orders")
	ErrUpdateOrderStatus = errors.New("orderRepository: can't update order status")

	ErrNoRowsAffected       = errors.New("orderRepository: no rows affected")
	ErrRepositoryConnection = errors.New("orderRepository: error connecting to the repository")

	ErrGetAllDeliveries     = errors.New("orderRepository: can't get info about all deliveries")
	ErrGetAllPendingOrders  = errors.New("orderRepository: can't get info about all pending orders")
	ErrCreateDelivery       = errors.New("orderRepository: can't create delivery")
	ErrUpdateDeliveryStatus = errors.New("orderRepository: can't update delivery status")
)
