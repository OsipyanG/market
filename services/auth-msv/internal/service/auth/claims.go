package authservice

import (
	"errors"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidClaims = errors.New("invalid claims")

func parseClaims(claims jwt.Claims) (*model.JWTClaims, error) {
	mapClaims, ok := claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	userID, ok := mapClaims["user_id"].(string)
	if !ok {
		return nil, ErrInvalidClaims
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errwrap.Wrap(ErrInvalidClaims, err)
	}

	accessLvl, ok := mapClaims["access_level"].(float64)
	if !ok {
		return nil, ErrInvalidClaims
	}

	modelClaims := &model.JWTClaims{
		UserID:    userUUID,
		AccessLvl: model.AccessLevel(accessLvl),
	}

	return modelClaims, nil
}
