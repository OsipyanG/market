package service

import (
	"context"
	"fmt"

	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/repository"
)

type CatalogService struct {
	catalogRepository repository.CatalogRepository
}

func NewCatalogService(catalogRepository repository.CatalogRepository) *CatalogService {
	return &CatalogService{
		catalogRepository: catalogRepository,
	}
}

func (s *CatalogService) GetCatalog(ctx context.Context, offset int, limit int) ([]model.Product, error) {
	result, err := s.catalogRepository.GetCatalog(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetCatalog, err)
	}

	return result, nil
}
