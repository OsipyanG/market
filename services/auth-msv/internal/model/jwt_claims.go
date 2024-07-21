package model

import "github.com/google/uuid"

type JWTClaims struct {
	UserID    uuid.UUID
	AccessLvl AccessLevel
}
