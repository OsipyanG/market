package service

import "errors"

var (
	ErrCreateOrder  = errors.New("orderService: failed to create order")
	ErrGetOrder     = errors.New("orderService: failed to get order")
	ErrUpdateOrder  = errors.New("orderService: failed to update order")
	ErrGetAllOrders = errors.New("orderService: failed to get all orders")

	ErrEmptyShopcart = errors.New("orderService: shopcart is empty")

	ErrGetAllDeliveries    = errors.New("orderService: failed to get all deliveries")
	ErrGetAllPendingOrders = errors.New("orderService: failed to get all pending orders")

	ErrAccessDenied  = errors.New("access denied")
	ErrPreCondition  = errors.New("the order status cannot be changed to this state")
	ErrInvalidStatus = errors.New("invalid status")
)
