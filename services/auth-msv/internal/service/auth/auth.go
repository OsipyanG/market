package authservice

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/OsipyanG/market/services/auth-msv/config"
	"github.com/OsipyanG/market/services/auth-msv/internal/cache"
	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/internal/storage"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrNewUser        = errors.New("auth-service: can't create new user")
	ErrLogin          = errors.New("auth-service: login error")
	ErrUpdateTokens   = errors.New("auth-service: can't update tokens")
	ErrUpdatePassword = errors.New("auth-service: can't update password")
	ErrLogout         = errors.New("auth-service: logout error")
	ErrGetJWTClaims   = errors.New("auth-service: can't get jwt claims")

	ErrCreateAccessToken  = errors.New("can't create access token")
	ErrInvalidMethod      = errors.New("invalid signature method")
	ErrInvalidAccessToken = errors.New("invalid access token")
)

type AuthService struct {
	config       *config.AuthServiceConfig
	storage      storage.UserStorage
	refreshCache cache.RefreshTokenCache
}

func New(conf *config.AuthServiceConfig, storage storage.UserStorage, refresh cache.RefreshTokenCache) *AuthService {
	return &AuthService{
		config:       conf,
		storage:      storage,
		refreshCache: refresh,
	}
}

func (ac *AuthService) NewUser(ctx context.Context, userCredentials *model.UserCredentials) (*model.Tokens, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, errwrap.Wrap(ErrNewUser, err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userCredentials.Password),
		ac.config.Security.HashCost)
	if err != nil {
		return nil, errwrap.Wrap(ErrNewUser, err)
	}

	hashedCred := &model.UserCredentials{
		Password: string(hashedPassword),
		Login:    userCredentials.Login,
	}

	err = ac.storage.CreateUser(ctx, userID, hashedCred)
	if err != nil {
		return nil, errwrap.Wrap(ErrNewUser, err)
	}

	tokens, err := ac.Login(ctx, userCredentials)
	if err != nil {
		return nil, errwrap.Wrap(ErrNewUser, err)
	}

	return tokens, nil
}

func (ac *AuthService) Login(ctx context.Context, userCredentials *model.UserCredentials) (*model.Tokens, error) {
	hashedPassword, err := ac.storage.GetPasswordByLogin(ctx, userCredentials.Login)
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userCredentials.Password))
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	userCredentials.Password = hashedPassword

	claims, err := ac.storage.GetJWTClaims(ctx, userCredentials)
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	accessToken, err := ac.createAccessToken(claims)
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	refreshToken, err := uuid.NewRandom()
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	err = ac.refreshCache.SetRefreshToken(refreshToken.String(), claims)
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	return &model.Tokens{
		Access:  accessToken,
		Refresh: refreshToken.String(),
	}, nil
}

func (ac *AuthService) UpdateTokens(_ context.Context, refreshToken string) (*model.Tokens, error) {
	claims, err := ac.refreshCache.GetJWTClaims(refreshToken)
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, err)
	}

	accessToken, err := ac.createAccessToken(claims)
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, err)
	}

	err = ac.refreshCache.DeleteRefreshToken(refreshToken)
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, err)
	}

	newRefresh, err := uuid.NewRandom()
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, err)
	}

	err = ac.refreshCache.SetRefreshToken(newRefresh.String(), claims)
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, err)
	}

	return &model.Tokens{
		Access:  accessToken,
		Refresh: newRefresh.String(),
	}, nil
}

func (ac *AuthService) UpdatePassword(ctx context.Context, userID uuid.UUID, userPasswords *model.UserPasswords) error {
	hashedPassword, err := ac.storage.GetPasswordByID(ctx, userID)
	if err != nil {
		return errwrap.Wrap(ErrUpdatePassword, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPasswords.Old))
	if err != nil {
		return errwrap.Wrap(ErrUpdatePassword, err)
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(userPasswords.New),
		ac.config.Security.HashCost)
	if err != nil {
		return errwrap.Wrap(ErrUpdatePassword, err)
	}

	err = ac.storage.UpdatePassword(ctx, userID, string(newHashedPassword))

	return errwrap.WrapIfErr(ErrUpdatePassword, err)
}

func (ac *AuthService) Logout(_ context.Context, refreshToken string) error {
	err := ac.refreshCache.DeleteRefreshToken(refreshToken)

	return errwrap.WrapIfErr(ErrLogout, err)
}

func (ac *AuthService) GetJWTClaims(_ context.Context, accessToken string) (*model.JWTClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidMethod
		}

		return ac.config.Security.SecretKey, nil
	})
	if err != nil {
		return nil, errwrap.Wrap(ErrGetJWTClaims, err)
	}

	if !token.Valid {
		return nil, errwrap.Wrap(ErrGetJWTClaims, ErrInvalidAccessToken)
	}

	claims, err := parseClaims(token.Claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (ac *AuthService) createAccessToken(claims *model.JWTClaims) (string, error) {
	jwtClaims := jwt.MapClaims{
		"user_id":      claims.UserID,
		"access_level": claims.AccessLvl,
		"exp": time.Now().Add(ac.config.Access.Timeout).Unix() +
			rand.Int63n(int64(ac.config.Access.Jitter.Seconds())),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	accessTokenStr, err := accessToken.SignedString(ac.config.Security.SecretKey)
	if err != nil {
		return "", errwrap.Wrap(ErrCreateAccessToken, err)
	}

	return accessTokenStr, nil
}
