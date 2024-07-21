package grpc

import (
	"context"

	warehousepb "github.com/OsipyanG/market/protos/warehouse"
	"github.com/OsipyanG/market/services/order-msv/internal/model"
	"github.com/OsipyanG/market/services/order-msv/internal/transport/warehouse"
	"github.com/OsipyanG/market/services/order-msv/pkg/errwrap"
	"github.com/google/uuid"
)

type Client struct {
	client warehousepb.WarehouseClient
}

func NewClient(client warehousepb.WarehouseClient) *Client {
	return &Client{client: client}
}

func (c *Client) ReserveProducts(ctx context.Context, products []model.OrderItem) error {
	productsProto := modelToProto(products)

	_, err := c.client.ReserveProducts(ctx, &warehousepb.ReserveProductsRequest{
		Products: productsProto,
	})
	if err != nil {
		return errwrap.Wrap(warehouse.ErrReserveProducts, err)
	}

	return nil
}

func (c *Client) FreeReservedProducts(ctx context.Context, products []model.OrderItem) error {
	productsProto := modelToProto(products)

	_, err := c.client.FreeReservedProducts(ctx, &warehousepb.FreeReservedProductsRequest{
		Products: productsProto,
	})
	if err != nil {
		return errwrap.Wrap(warehouse.ErrFreeProducts, err)
	}

	return nil
}

func (c *Client) DeleteReservedProducts(ctx context.Context, products []model.OrderItem) error {
	productsProto := modelToProto(products)

	_, err := c.client.DeleteReservedProducts(ctx, &warehousepb.DeleteReservedProductsRequest{
		Products: productsProto,
	})
	if err != nil {
		return errwrap.Wrap(warehouse.ErrDeleteProducts, err)
	}

	return nil
}

func (c *Client) GetProductsPrices(ctx context.Context, ids []uuid.UUID) ([]model.OrderItem, error) {

	productIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		productIDs = append(productIDs, id.String())
	}

	products, err := c.client.GetProductPrices(ctx, &warehousepb.GetProductsPricesRequest{
		ProductIds: productIDs,
	})
	if err != nil {
		return nil, errwrap.Wrap(warehouse.ErrGetProductsPrices, err)
	}

	orderItems := make([]model.OrderItem, 0, len(products.Products))
	for _, product := range products.Products {
		orderItems = append(orderItems, model.OrderItem{
			ProductID: uuid.MustParse(product.ProductId),
			Price:     int(product.Price),
		})

	}

	return orderItems, nil
}

func modelToProto(products []model.OrderItem) []*warehousepb.ProductQuantity {
	productsProto := make([]*warehousepb.ProductQuantity, 0, len(products))

	for _, item := range products {
		productsProto = append(productsProto, &warehousepb.ProductQuantity{
			ProductId: item.ProductID.String(),
			Quantity:  int64(item.Quantity),
		})
	}

	return productsProto
}
