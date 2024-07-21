package service

import (
	"context"
	"fmt"

	"github.com/OsipyanG/market/services/warehouse-msv/internal/model"
	"github.com/OsipyanG/market/services/warehouse-msv/internal/repository"
	"github.com/google/uuid"
)

type WarehouseService struct {
	repository repository.WarehouseRepository
}

func NewWarehouseService(repository repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{
		repository: repository,
	}
}

func (s *WarehouseService) ReserveProducts(ctx context.Context, products []model.ProductQuantity) error {
	err := s.repository.ReserveProducts(ctx, products)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReserveProduct, err)
	}

	return nil
}

func (s *WarehouseService) FreeReservedProducts(ctx context.Context, products []model.ProductQuantity) error {
	err := s.repository.FreeReservedProducts(ctx, products)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFreeReservedProducts, err)
	}

	return nil
}

func (s *WarehouseService) DeleteReservedProducts(ctx context.Context, products []model.ProductQuantity) error {
	err := s.repository.DeleteReservedProducts(ctx, products)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrDeleteReservedProduct, err)
	}

	return nil
}

func (s *WarehouseService) GetProductPrices(ctx context.Context, productIDs []uuid.UUID) ([]model.ProductPrice, error) {
	prices, err := s.repository.GetProductPrices(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrGetProductPrices, err)
	}

	return prices, nil
}
