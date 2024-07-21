package converter

import (
	"github.com/OsipyanG/market/protos/auth"
	"github.com/OsipyanG/market/protos/jwt"
	"github.com/OsipyanG/market/services/auth-msv/internal/model"
)

func ConvertToModelCredentials(cred *auth.UserCredentials) *model.UserCredentials {
	if cred == nil {
		return nil
	}

	return &model.UserCredentials{
		Login:    cred.GetLogin(),
		Password: cred.GetPassword(),
	}
}

func ConvertFromModelTokens(tokens *model.Tokens) *auth.Tokens {
	if tokens == nil {
		return nil
	}

	return &auth.Tokens{
		AccessToken:  tokens.Access,
		RefreshToken: tokens.Refresh,
	}
}

func ConvertFromModelClaims(claims *model.JWTClaims) *jwt.JWTClaims {
	if claims == nil {
		return nil
	}

	return &jwt.JWTClaims{
		UserId:      claims.UserID.String(),
		AccessLevel: int32(claims.AccessLvl),
	}
}

func ConvertFromModelUser(user *model.UserInfo) *auth.User {
	return &auth.User{
		Id:          user.ID.String(),
		Login:       user.Login,
		AccessLevel: uint32(user.AccessLvl),
	}
}
