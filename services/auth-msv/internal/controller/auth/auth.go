package auth

import (
	"context"
	"errors"

	"github.com/OsipyanG/market/protos/auth"
	"github.com/OsipyanG/market/protos/jwt"
	"github.com/OsipyanG/market/services/auth-msv/internal/controller"
	"github.com/OsipyanG/market/services/auth-msv/internal/converter"
	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
)

var (
	ErrNewUser        = errors.New("auth controller: can't create new user")
	ErrLogin          = errors.New("auth controller: login error")
	ErrUpdateTokens   = errors.New("auth controller: can't update tokens")
	ErrLogout         = errors.New("auth controller: logout error")
	ErrUpdatePassword = errors.New("auth controller: can't update password")
	ErrGetJWTClaims   = errors.New("auth controller: can't get jwt claims")
)

type Service interface {
	NewUser(ctx context.Context, userCredentials *model.UserCredentials) (*model.Tokens, error)
	Login(ctx context.Context, userCredentials *model.UserCredentials) (*model.Tokens, error)
	UpdateTokens(ctx context.Context, refreshToken string) (*model.Tokens, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, userPasswords *model.UserPasswords) error
	Logout(ctx context.Context, refreshToken string) error
	GetJWTClaims(ctx context.Context, accessToken string) (*model.JWTClaims, error)
}

type Controller struct {
	auth.UnimplementedAuthServer
	service Service
}

func New(s Service) *Controller {
	return &Controller{
		service: s,
	}
}

func (ac *Controller) NewUser(ctx context.Context, credentials *auth.UserCredentials) (*auth.Tokens, error) {
	if credentials == nil {
		return nil, errwrap.Wrap(ErrNewUser, controller.ErrNoPayload)
	}

	modelCredentials := converter.ConvertToModelCredentials(credentials)

	modelTokens, err := ac.service.NewUser(ctx, modelCredentials)
	if err != nil {
		return nil, errwrap.Wrap(ErrNewUser, err)
	}

	return converter.ConvertFromModelTokens(modelTokens), nil
}

func (ac *Controller) Login(ctx context.Context, credentials *auth.UserCredentials) (*auth.Tokens, error) {
	if credentials == nil {
		return nil, errwrap.Wrap(ErrLogin, controller.ErrNoPayload)
	}

	modelCredentials := converter.ConvertToModelCredentials(credentials)

	modelTokens, err := ac.service.Login(ctx, modelCredentials)
	if err != nil {
		return nil, errwrap.Wrap(ErrLogin, err)
	}

	return converter.ConvertFromModelTokens(modelTokens), nil
}

func (ac *Controller) UpdateTokens(ctx context.Context, refresh *auth.RefreshToken) (*auth.Tokens, error) {
	if refresh == nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, controller.ErrNoPayload)
	}

	modelTokens, err := ac.service.UpdateTokens(ctx, refresh.GetValue())
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdateTokens, err)
	}

	return converter.ConvertFromModelTokens(modelTokens), nil
}

func (ac *Controller) UpdatePassword(ctx context.Context, req *auth.RequestUpdatePassword) (*empty.Empty, error) {
	if req == nil {
		return nil, errwrap.Wrap(ErrNewUser, controller.ErrNoPayload)
	}

	passwords := &model.UserPasswords{
		New: req.GetNewPassword(),
		Old: req.GetOldPassword(),
	}

	claims := req.GetJwtClaims()
	if claims == nil {
		return nil, errwrap.Wrap(ErrUpdatePassword, controller.ErrNoJWTClaims)
	}

	userID, err := uuid.Parse(claims.GetUserId())
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdatePassword, err)
	}

	err = ac.service.UpdatePassword(ctx, userID, passwords)
	if err != nil {
		return nil, errwrap.Wrap(ErrUpdatePassword, err)
	}

	return &empty.Empty{}, nil
}

func (ac *Controller) Logout(ctx context.Context, refresh *auth.RefreshToken) (*empty.Empty, error) {
	err := ac.service.Logout(ctx, refresh.GetValue())
	if err != nil {
		return nil, errwrap.Wrap(ErrLogout, err)
	}

	return &empty.Empty{}, nil
}

func (ac *Controller) GetJWTClaims(ctx context.Context, access *auth.AccessToken) (*jwt.JWTClaims, error) {
	claims, err := ac.service.GetJWTClaims(ctx, access.GetValue())
	if err != nil {
		return nil, errwrap.Wrap(ErrGetJWTClaims, err)
	}

	return converter.ConvertFromModelClaims(claims), nil
}
