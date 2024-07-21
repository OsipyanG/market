package model

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID
	CustomerID uuid.UUID
	Status     string
	Address    string
	Items      []OrderItem
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type OrderItem struct {
	ProductID uuid.UUID
	Quantity  int
	Price     int
}
