package model

import "github.com/google/uuid"

type Delivery struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
	Status    string
}
