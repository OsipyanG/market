package service

import "errors"

var (
	ErrGetCatalog = errors.New("catalogService: can't get catalog")

	ErrReserveProduct        = errors.New("warehouseService: can't reserve product")
	ErrFreeReservedProducts  = errors.New("warehouseService: can't free reserved products")
	ErrDeleteReservedProduct = errors.New("warehouseService: can't delete reserved product")
	ErrGetProductPrices      = errors.New("warehouseService: can't get product prices")
)
