package grpc

import (
	"context"

	jwtpb "github.com/OsipyanG/market/protos/jwt"
	shopcartpb "github.com/OsipyanG/market/protos/shopcart"
	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/OsipyanG/market/services/order-msv/internal/transport/shopcart"
	"github.com/OsipyanG/market/services/order-msv/pkg/errwrap"
	"github.com/google/uuid"
)

type Client struct {
	client shopcartpb.UserShopcartClient
}

func NewClient(client shopcartpb.UserShopcartClient) *Client {
	return &Client{client: client}
}

func (c *Client) GetProducts(ctx context.Context, userID uuid.UUID) ([]model.OrderItem, error) {
	resp, err := c.client.GetProducts(ctx, &shopcartpb.RequestByID{
		JwtClaims: &jwtpb.JWTClaims{UserId: userID.String()},
	})
	if err != nil {
		return nil, errwrap.Wrap(shopcart.ErrGetProducts, err)
	}

	orderItems := make([]model.OrderItem, 0, len(resp.GetProducts()))

	for _, item := range resp.GetProducts() {
		productID, err := uuid.Parse(item.GetId())
		if err != nil {
			return nil, errwrap.Wrap(shopcart.ErrGetProducts, err)
		}

		orderItems = append(orderItems, model.OrderItem{
			ProductID: productID,
			Quantity:  int(item.GetQuantity()),
		})
	}

	return orderItems, nil
}

func (c *Client) Clear(ctx context.Context, userID uuid.UUID) error {
	_, err := c.client.Clear(ctx, &shopcartpb.RequestByID{JwtClaims: &jwtpb.JWTClaims{UserId: userID.String()}})
	if err != nil {
		return errwrap.Wrap(shopcart.ErrClear, err)
	}

	return nil
}
