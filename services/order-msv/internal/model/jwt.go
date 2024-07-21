package model

import "github.com/google/uuid"

type JwtClaims struct {
	UserID      uuid.UUID
	AccessLevel int
}
