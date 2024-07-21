package warehouse

import "errors"

var (
	ErrReserveProducts   = errors.New("warehouseClient: failed to reserve products")
	ErrFreeProducts      = errors.New("warehouseClient: failed to free products")
	ErrDeleteProducts    = errors.New("warehouseClient: failed to delete products")
	ErrGetProductsPrices = errors.New("warehouseClient: failed to get products prices")
)
