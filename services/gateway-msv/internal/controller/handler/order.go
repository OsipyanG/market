package handler

import (
	"context"
	"net/http"

	jwt "github.com/OsipyanG/market/protos/jwt"
	order "github.com/OsipyanG/market/protos/order"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	client order.OrderServiceClient
}

func NewOrderHandler(o order.OrderServiceClient) *OrderHandler {
	return &OrderHandler{client: o}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	createOrderRequest := &order.CreateOrderRequest{}

	if err := c.ShouldBindJSON(createOrderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	claims := c.MustGet("claims")
	createOrderRequest.JwtClaims = claims.(*jwt.JWTClaims)

	order, err := h.client.CreateOrder(context.TODO(), createOrderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	getOrderRequest := &order.GetOrderRequest{}

	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order_id is required"})

		return
	}

	claims := c.MustGet("claims")
	getOrderRequest.JwtClaims = claims.(*jwt.JWTClaims)
	getOrderRequest.OrderId = orderID

	order, err := h.client.GetOrder(context.TODO(), getOrderRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	getAllOrdersRequest := &order.GetAllOrdersRequest{}

	claims := c.MustGet("claims")
	getAllOrdersRequest.JwtClaims = claims.(*jwt.JWTClaims)

	orders, err := h.client.GetAllOrders(context.TODO(), getAllOrdersRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	updateOrderStatusRequest := &order.UpdateOrderStatusRequest{}

	if err := c.ShouldBindJSON(updateOrderStatusRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	claims := c.MustGet("claims")
	updateOrderStatusRequest.JwtClaims = claims.(*jwt.JWTClaims)

	_, err := h.client.UpdateOrderStatus(context.TODO(), updateOrderStatusRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, "order status updated")
}

func (h *OrderHandler) GetAllPendingOrders(c *gin.Context) {
	getAllPendingOrdersReq := &order.GetAllPendingOrdersRequest{}

	if err := c.ShouldBindJSON(getAllPendingOrdersReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	claims := c.MustGet("claims")
	getAllPendingOrdersReq.JwtClaims = claims.(*jwt.JWTClaims)

	orders, err := h.client.GetAllPendingOrders(context.TODO(), getAllPendingOrdersReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetAllDeliveries(c *gin.Context) {
	getAllDeliveriesReq := &order.GetAllDeliveriesRequest{}

	if err := c.ShouldBindJSON(getAllDeliveriesReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	claims := c.MustGet("claims")
	getAllDeliveriesReq.JwtClaims = claims.(*jwt.JWTClaims)

	orders, err := h.client.GetAllDeliveries(context.TODO(), getAllDeliveriesReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, orders)
}
