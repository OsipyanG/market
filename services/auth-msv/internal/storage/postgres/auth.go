package postgres

import (
	"context"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/internal/storage"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/google/uuid"
)

func (s *Storage) CreateUser(ctx context.Context, userID uuid.UUID, userCredentials *model.UserCredentials) error {
	cmdTag, err := s.pool.Exec(ctx, "INSERT INTO Users (id, login, password, access_level) VALUES ($1, $2, $3, $4)",
		userID, userCredentials.Login, userCredentials.Password, model.Buyer)
	if err != nil {
		return errwrap.Wrap(storage.ErrCreateUser, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(storage.ErrCreateUser, storage.ErrNoRowsAffected)
	}

	return nil
}

func (s *Storage) GetPasswordByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var password string

	err := s.pool.QueryRow(ctx, "SELECT password FROM Users WHERE id=$1",
		userID).Scan(&password)
	if err != nil {
		return "", errwrap.Wrap(storage.ErrUpdatePassword, err)
	}

	return password, nil
}

func (s *Storage) GetPasswordByLogin(ctx context.Context, login string) (string, error) {
	var password string

	err := s.pool.QueryRow(ctx, "SELECT password FROM Users WHERE login = $1", login).Scan(&password)
	if err != nil {
		return "", errwrap.Wrap(storage.ErrUpdatePassword, err)
	}

	return password, nil
}

func (s *Storage) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	cmdTag, err := s.pool.Exec(ctx, "UPDATE Users SET password = $1 WHERE id = $2",
		newPassword, userID)
	if err != nil {
		return errwrap.Wrap(storage.ErrUpdatePassword, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(storage.ErrUpdatePassword, storage.ErrNoRowsAffected)
	}

	return nil
}

func (s *Storage) GetJWTClaims(ctx context.Context, userCred *model.UserCredentials) (*model.JWTClaims, error) {
	claims := &model.JWTClaims{}

	err := s.pool.QueryRow(ctx, "SELECT id, access_level FROM Users WHERE login=$1 AND password=$2",
		userCred.Login, userCred.Password).Scan(&claims.UserID, &claims.AccessLvl)
	if err != nil {
		return nil, errwrap.Wrap(storage.ErrGetJWTClaims, err)
	}

	return claims, nil
}
