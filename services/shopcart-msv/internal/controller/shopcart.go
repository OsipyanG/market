package controller

import (
	"context"
	"errors"

	"github.com/OsipyanG/market/protos/jwt"
	"github.com/OsipyanG/market/protos/shopcart"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/converter"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/model"
	"github.com/OsipyanG/market/services/shopcart-msv/pkg/errwrap"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
)

var (
	ErrAddProduct    = errors.New("controller: can't add product")
	ErrDeleteProduct = errors.New("controller: can't delete product")
	ErrGetProducts   = errors.New("controller: can't get products")
	ErrClear         = errors.New("controller: can't clear shopcart")

	ErrNoJWTClaims   = errors.New("no jwt claims")
	ErrInvalidUserID = errors.New("invalid user_id")
)

type ShopCartService interface {
	AddProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error
	DeleteProduct(ctx context.Context, userID uuid.UUID, product *model.Product) error
	GetProducts(ctx context.Context, userID uuid.UUID) ([]*model.Product, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type ShopCartController struct {
	shopcart.UnimplementedUserShopcartServer
	service ShopCartService
}

func New(service ShopCartService) *ShopCartController {
	return &ShopCartController{
		service: service,
	}
}

func (s *ShopCartController) AddProduct(ctx context.Context, req *shopcart.RequestByIDWithProduct) (*empty.Empty, error) {
	userID, err := parseUserID(req.GetJwtClaims())
	if err != nil {
		return nil, errwrap.Wrap(ErrAddProduct, err)
	}

	product, err := converter.ConvertToModelProduct(req.GetProduct())
	if err != nil {
		return nil, errwrap.Wrap(ErrAddProduct, err)
	}

	err = s.service.AddProduct(ctx, userID, product)
	if err != nil {
		return nil, errwrap.Wrap(ErrAddProduct, err)
	}

	return &empty.Empty{}, nil
}

func (s *ShopCartController) DeleteProduct(ctx context.Context, req *shopcart.RequestByIDWithProduct) (*empty.Empty, error) {
	userID, err := parseUserID(req.GetJwtClaims())
	if err != nil {
		return nil, errwrap.Wrap(ErrDeleteProduct, err)
	}

	product, err := converter.ConvertToModelProduct(req.GetProduct())
	if err != nil {
		return nil, errwrap.Wrap(ErrDeleteProduct, err)
	}

	err = s.service.DeleteProduct(ctx, userID, product)
	if err != nil {
		return nil, errwrap.Wrap(ErrDeleteProduct, err)
	}

	return &empty.Empty{}, nil
}

func (s *ShopCartController) GetProducts(ctx context.Context, req *shopcart.RequestByID) (*shopcart.GetProductsResponse, error) {
	userID, err := parseUserID(req.GetJwtClaims())
	if err != nil {
		return nil, errwrap.Wrap(ErrGetProducts, err)
	}

	products, err := s.service.GetProducts(ctx, userID)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetProducts, err)
	}

	grpcProducts := make([]*shopcart.Product, len(products))
	for i := range products {
		grpcProducts[i] = converter.ConvertFromModelProduct(products[i])
	}

	return &shopcart.GetProductsResponse{
		Products: grpcProducts,
	}, nil
}

func (s *ShopCartController) Clear(ctx context.Context, req *shopcart.RequestByID) (*empty.Empty, error) {
	userID, err := parseUserID(req.GetJwtClaims())
	if err != nil {
		return nil, errwrap.Wrap(ErrClear, err)
	}

	err = s.service.Clear(ctx, userID)
	if err != nil {
		return nil, errwrap.Wrap(ErrClear, err)
	}

	return &empty.Empty{}, nil
}

func parseUserID(claims *jwt.JWTClaims) (uuid.UUID, error) {
	if claims == nil {
		return uuid.Nil, ErrNoJWTClaims
	}

	userID, err := uuid.Parse(claims.GetUserId())
	if err != nil {
		return uuid.Nil, errwrap.Wrap(ErrInvalidUserID, err)
	}

	return userID, nil
}
