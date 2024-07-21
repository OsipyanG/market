package model

import "github.com/google/uuid"

type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Available   int64
	Quantity    int64
	Price       int64
}

type ProductQuantity struct {
	ID       uuid.UUID
	Quantity int64
}

type ProductPrice struct {
	ID    uuid.UUID
	Price int64
}
