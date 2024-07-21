package storage

import (
	"context"
	"errors"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/google/uuid"
)

var (
	ErrCreateUser     = errors.New("storage: can't create user")
	ErrGetJWTClaims   = errors.New("storage: can't get jwt claims")
	ErrUpdatePassword = errors.New("storage: can't update password")

	ErrDeleteUser        = errors.New("storage: can't delete user")
	ErrSetAccessLevel    = errors.New("storage: can't set access level")
	ErrAllUsersWithLevel = errors.New("storage: can't get all users with level")

	ErrStorageConnection = errors.New("storage: can't connect to the database")
	ErrNoRowsAffected    = errors.New("no rows affected")
	ErrNotFound          = errors.New("not found")
)

type UserStorage interface {
	CreateUser(ctx context.Context, userID uuid.UUID, userCredentials *model.UserCredentials) error
	GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error)
	GetPasswordByLogin(ctx context.Context, login string) (string, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	GetJWTClaims(ctx context.Context, userCredentials *model.UserCredentials) (*model.JWTClaims, error)
}

type AdminStorage interface {
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	SetAccessLevel(ctx context.Context, userID uuid.UUID, lvl model.AccessLevel) error
	GetAllUsersWithLevel(ctx context.Context, lvl model.AccessLevel) ([]*model.UserInfo, error)
}
