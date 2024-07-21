package cache

import (
	"errors"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
)

var (
	ErrSetRefreshToken    = errors.New("cache: cant set refresh token")
	ErrGetJWTClaims       = errors.New("cache: cant get jwt claims")
	ErrUpdateToken        = errors.New("cache: cant update token")
	ErrDeleteRefreshToken = errors.New("cache: cant delete refresh token")

	ErrConnection = errors.New("cache: connection error")
)

type RefreshTokenCache interface {
	SetRefreshToken(refreshToken string, claims *model.JWTClaims) error
	GetJWTClaims(refreshToken string) (*model.JWTClaims, error)
	DeleteRefreshToken(refreshToken string) error
}
