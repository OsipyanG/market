package memcached

import (
	"strconv"
	"strings"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/google/uuid"
)

func pack(claims *model.JWTClaims) string {
	var value strings.Builder

	value.WriteString(claims.UserID.String())
	value.WriteString("|")
	value.WriteString(strconv.Itoa(int(claims.AccessLvl)))

	return value.String()
}

func unpack(value string) (*model.JWTClaims, error) {
	claims := strings.Split(value, "|")

	userID, err := uuid.Parse(claims[0])
	if err != nil {
		return nil, errwrap.Wrap(ErrUnpackValue, err)
	}

	lvl, err := strconv.Atoi(claims[1])
	if err != nil {
		return nil, errwrap.Wrap(ErrUnpackValue, err)
	}

	return &model.JWTClaims{
		UserID:    userID,
		AccessLvl: model.AccessLevel(lvl),
	}, nil
}
