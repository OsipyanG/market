package service

import (
	"context"
	"errors"

	"github.com/OsipyanG/market/services/shopcart-msv/internal/model"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/storage"
	"github.com/OsipyanG/market/services/shopcart-msv/pkg/errwrap"
	"github.com/google/uuid"
)

var (
	ErrAddProduct    = errors.New("service: can't add product")
	ErrDeleteProduct = errors.New("service: can't delete product")
	ErrGetProducts   = errors.New("service: cant't get products")
	ErrClear         = errors.New("service: can't clear shopcart")
)

type ShopCartService struct {
	storage storage.ShopCartStorage
}

func New(storage storage.ShopCartStorage) *ShopCartService {
	return &ShopCartService{
		storage: storage,
	}
}

func (s *ShopCartService) AddProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error {
	err := s.storage.AddProduct(ctx, userID, product)
	if err != nil {
		return errwrap.Wrap(ErrAddProduct, err)
	}

	return nil
}

func (s *ShopCartService) DeleteProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error {
	err := s.storage.DeleteProduct(ctx, userID, product)
	if err != nil {
		return errwrap.Wrap(ErrAddProduct, err)
	}

	return nil
}

func (s *ShopCartService) GetProducts(ctx context.Context, userID uuid.UUID) ([]*model.Product, error) {
	products, err := s.storage.GetProducts(ctx, userID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetProducts, err)
	}

	return products, nil
}

func (s *ShopCartService) Clear(ctx context.Context, userID uuid.UUID) error {
	err := s.storage.Clear(ctx, userID)
	if err != nil {
		return errwrap.Wrap(ErrGetProducts, err)
	}

	return nil
}
