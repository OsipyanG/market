package repository

import "errors"

var (
	ErrReserveProducts        = errors.New("repository: can't reserve products")
	ErrFreeReservedProducts   = errors.New("repository: can't free reserved products")
	ErrDeleteReservedProducts = errors.New("repository: can't delete reserved products")
	ErrGetProductPrices       = errors.New("repository: can't get product prices")
	ErrNotSoManyReserved      = errors.New("there is not so many reserved products")
	ErrNotEnoughStock         = errors.New("not enough stock for products")

	ErrGetCatalog = errors.New("repository: can't get catalog")

	ErrRepositoryConnection = errors.New("repository: can't connect to the database")
	ErrNoRowsAffected       = errors.New("no rows affected")
)
