package controller

import (
	"context"

	warehousepb "github.com/OsipyanG/market/protos/warehouse"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CatalogService interface {
	GetCatalog(ctx context.Context, offset int, limit int) ([]model.Product, error)
}

type CatalogController struct {
	warehousepb.UnimplementedCatalogServer
	service CatalogService
}

func NewCatalogController(service CatalogService) *CatalogController {
	return &CatalogController{
		service: service,
	}
}

func (c *CatalogController) GetCatalog(ctx context.Context, req *warehousepb.GetCatalogRequest) (
	*warehousepb.GetCatalogResponse, error,
) {
	result, err := c.service.GetCatalog(ctx, int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	protoProduct := make([]*warehousepb.Product, len(result))
	for i, product := range result {
		protoProduct[i] = &warehousepb.Product{
			ProductId:   product.ID.String(),
			Name:        product.Name,
			Description: product.Description,
			Available:   product.Available,
			Quantity:    product.Quantity,
			Price:       product.Price,
		}
	}

	return &warehousepb.GetCatalogResponse{Products: protoProduct}, nil
}
