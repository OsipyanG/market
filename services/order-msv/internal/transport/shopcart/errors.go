package shopcart

import "errors"

var (
	ErrGetProducts = errors.New("shopcartClient: can't get products from shopcart")
	ErrClear       = errors.New("shopcartClient: can't clear shopcart")
)
