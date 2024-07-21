package adminservice

import (
	"context"
	"errors"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/internal/storage"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/google/uuid"
)

var (
	ErrDeleteUser           = errors.New("admin-service: can't delete user")
	ErrSetAccessLevel       = errors.New("admin-service: can't set access level")
	ErrGetAllUsersWithLevel = errors.New("admin-service: can't get all users with specified level")
)

type AdminService struct {
	storage storage.AdminStorage
}

func New(storage storage.AdminStorage) *AdminService {
	return &AdminService{
		storage: storage,
	}
}

func (as *AdminService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := as.storage.DeleteUser(ctx, userID)

	return errwrap.WrapIfErr(ErrDeleteUser, err)
}

// SetAccessLevel changes user's access level.
// WARNING: To update the access level, the user will need to log in to the account.
func (as *AdminService) SetAccessLevel(ctx context.Context, userID uuid.UUID, lvl model.AccessLevel) error {
	err := as.storage.SetAccessLevel(ctx, userID, lvl)

	return errwrap.WrapIfErr(ErrSetAccessLevel, err)
}

func (as *AdminService) GetAllUsersWithLevel(ctx context.Context, lvl model.AccessLevel) ([]*model.UserInfo, error) {
	users, err := as.storage.GetAllUsersWithLevel(ctx, lvl)
	if err != nil {
		return nil, errwrap.Wrap(ErrGetAllUsersWithLevel, err)
	}

	return users, nil
}
