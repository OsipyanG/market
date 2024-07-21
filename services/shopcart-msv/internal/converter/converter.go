package converter

import (
	"errors"

	proto "github.com/OsipyanG/market/protos/shopcart"
	"github.com/OsipyanG/market/services/shopcart-msv/internal/model"
	"github.com/OsipyanG/market/services/shopcart-msv/pkg/errwrap"
	"github.com/google/uuid"
)

var ErrCantConvertToModel = errors.New("can't convert to model product")

func ConvertToModelProduct(protoProduct *proto.Product) (*model.Product, error) {
	if protoProduct == nil {
		return nil, ErrCantConvertToModel
	}

	id, err := uuid.Parse(protoProduct.GetId())
	if err != nil {
		return nil, errwrap.Wrap(ErrCantConvertToModel, err)
	}

	return &model.Product{
		ID:       id,
		Quantity: int(protoProduct.GetQuantity()),
	}, nil
}

func ConvertFromModelProduct(product *model.Product) *proto.Product {
	return &proto.Product{
		Id:       product.ID.String(),
		Quantity: uint64(product.Quantity),
	}
}
