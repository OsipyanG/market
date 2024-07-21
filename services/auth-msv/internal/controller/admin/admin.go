package admin

import (
	"context"
	"errors"

	auth "github.com/OsipyanG/market/protos/auth"
	jwt "github.com/OsipyanG/market/protos/jwt"
	"github.com/OsipyanG/market/services/auth-msv/internal/controller"
	"github.com/OsipyanG/market/services/auth-msv/internal/converter"
	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
)

var (
	ErrDeleteUser           = errors.New("admin controller: can't delete user")
	ErrSetAccessLevel       = errors.New("admin controller: can't set access level")
	ErrGetAllUsersWithLevel = errors.New("admin controller: can't get all users with specified level")

	ErrInvalidAccessLevel = errors.New("invalid access level is specified")
)

type Service interface {
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	SetAccessLevel(ctx context.Context, userID uuid.UUID, lvl model.AccessLevel) error
	GetAllUsersWithLevel(ctx context.Context, lvl model.AccessLevel) ([]*model.UserInfo, error)
}

type Controller struct {
	auth.UnimplementedAuthAdminServer
	service Service
}

func New(s Service) *Controller {
	return &Controller{
		service: s,
	}
}

func (ac *Controller) DeleteUser(ctx context.Context, req *auth.RequestByUserID) (*empty.Empty, error) {
	if req == nil {
		return nil, errwrap.Wrap(ErrDeleteUser, controller.ErrNoPayload)
	}

	if err := isAdmin(req.GetJwtClaims()); err != nil {
		return nil, errwrap.Wrap(ErrGetAllUsersWithLevel, err)
	}

	strUserID := req.GetUserId()

	userID, err := uuid.Parse(strUserID)
	if err != nil {
		return nil, errwrap.Wrap(ErrDeleteUser, err)
	}

	err = ac.service.DeleteUser(ctx, userID)
	if err != nil {
		return nil, errwrap.Wrap(ErrDeleteUser, err)
	}

	return &empty.Empty{}, nil
}

func (ac *Controller) SetAccessLevel(ctx context.Context, req *auth.SetAccessLevelRequest) (*empty.Empty, error) {
	if req == nil {
		return nil, errwrap.Wrap(ErrSetAccessLevel, controller.ErrNoPayload)
	}

	if err := isAdmin(req.GetJwtClaims()); err != nil {
		return nil, errwrap.Wrap(ErrSetAccessLevel, err)
	}

	strUserID := req.GetUserId()

	userID, err := uuid.Parse(strUserID)
	if err != nil {
		return nil, errwrap.Wrap(ErrSetAccessLevel, err)
	}

	lvl := model.AccessLevel(req.GetLvl())
	if !lvl.IsValid() {
		return nil, errwrap.Wrap(ErrSetAccessLevel, ErrInvalidAccessLevel)
	}

	err = ac.service.SetAccessLevel(ctx, userID, lvl)
	if err != nil {
		return nil, errwrap.Wrap(ErrSetAccessLevel, err)
	}

	return &empty.Empty{}, nil
}

func (ac *Controller) GetAllUsersWithLevel(ctx context.Context, req *auth.RequestByLevel) (*auth.UsersInfoResponse, error) {
	if req == nil {
		return nil, errwrap.Wrap(ErrGetAllUsersWithLevel, controller.ErrNoPayload)
	}

	if err := isAdmin(req.GetJwtClaims()); err != nil {
		return nil, errwrap.Wrap(ErrGetAllUsersWithLevel, err)
	}

	lvl := model.AccessLevel(req.GetLvl())
	if !lvl.IsValid() {
		return nil, errwrap.Wrap(ErrGetAllUsersWithLevel, ErrInvalidAccessLevel)
	}

	users, err := ac.service.GetAllUsersWithLevel(ctx, lvl)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllUsersWithLevel, err)
	}

	response := &auth.UsersInfoResponse{
		Users: make([]*auth.User, len(users)),
	}

	for i := range users {
		response.Users[i] = converter.ConvertFromModelUser(users[i])
	}

	return response, nil
}

func isAdmin(claims *jwt.JWTClaims) error {
	if claims == nil {
		return controller.ErrNoJWTClaims
	}

	if claims.GetAccessLevel() != int32(model.Admin) {
		return controller.ErrAccessDenied
	}

	return nil
}
