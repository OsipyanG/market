package controller

import (
	"context"

	warehousepb "github.com/OsipyanG/market/protos/warehouse"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type WarehouseService interface {
	ReserveProducts(ctx context.Context, products []model.ProductQuantity) error
	FreeReservedProducts(ctx context.Context, products []model.ProductQuantity) error
	DeleteReservedProducts(ctx context.Context, products []model.ProductQuantity) error
	GetProductPrices(ctx context.Context, productIDs []uuid.UUID) ([]model.ProductPrice, error)
}

type WarehouseController struct {
	warehousepb.UnimplementedWarehouseServer
	service WarehouseService
}

func NewWarehouseController(service WarehouseService) *WarehouseController {
	return &WarehouseController{service: service}
}

func (c *WarehouseController) ReserveProducts(ctx context.Context, req *warehousepb.ReserveProductsRequest) (*empty.Empty, error) {
	productsRequest := req.GetProducts()

	ok := ValidateProductQuantity(productsRequest)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product quantity")
	}

	products := make([]model.ProductQuantity, 0, len(productsRequest))

	for _, product := range productsRequest {
		productID, err := uuid.Parse(product.GetProductId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		products = append(products, model.ProductQuantity{
			ID:       productID,
			Quantity: product.GetQuantity(),
		})
	}

	err := c.service.ReserveProducts(ctx, products)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (c *WarehouseController) FreeReservedProducts(ctx context.Context, req *warehousepb.FreeReservedProductsRequest) (*empty.Empty, error) {
	productsRequest := req.GetProducts()

	ok := ValidateProductQuantity(productsRequest)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product quantity")
	}

	products := make([]model.ProductQuantity, 0, len(productsRequest))

	for _, product := range productsRequest {
		productID, err := uuid.Parse(product.GetProductId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		products = append(products, model.ProductQuantity{
			ID:       productID,
			Quantity: product.GetQuantity(),
		})
	}

	err := c.service.FreeReservedProducts(ctx, products)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (c *WarehouseController) DeleteReservedProducts(ctx context.Context, req *warehousepb.DeleteReservedProductsRequest) (*empty.Empty, error) {
	productsRequest := req.GetProducts()

	ok := ValidateProductQuantity(productsRequest)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product quantity")
	}

	products := make([]model.ProductQuantity, 0, len(productsRequest))

	for _, product := range productsRequest {
		productID, err := uuid.Parse(product.GetProductId())
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		products = append(products, model.ProductQuantity{
			ID:       productID,
			Quantity: product.GetQuantity(),
		})
	}

	err := c.service.DeleteReservedProducts(ctx, products)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func ValidateProductQuantity(products []*warehousepb.ProductQuantity) bool {
	for _, product := range products {
		if product.GetQuantity() <= 0 {
			return false
		}
	}

	return true
}

func (c *WarehouseController) GetProductPrices(ctx context.Context, req *warehousepb.GetProductsPricesRequest) (*warehousepb.GetProductsPricesResponse, error) {
	productIDs := req.GetProductIds()

	ids := make([]uuid.UUID, 0, len(productIDs))

	for _, id := range productIDs {
		ProductUUID, err := uuid.Parse(id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}

		ids = append(ids, ProductUUID)
	}

	prices, err := c.service.GetProductPrices(ctx, ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &warehousepb.GetProductsPricesResponse{
		Products: make([]*warehousepb.ProductPrice, 0, len(prices)),
	}

	for _, price := range prices {
		response.Products = append(response.Products, &warehousepb.ProductPrice{
			ProductId: price.ID.String(),
			Price:     price.Price,
		})
	}

	return response, nil
}
