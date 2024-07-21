package postgres

import (
	"context"

	"github.com/OsipyanG/market/services/auth-msv/internal/model"
	"github.com/OsipyanG/market/services/auth-msv/internal/storage"
	"github.com/OsipyanG/market/services/auth-msv/pkg/errwrap"
	"github.com/google/uuid"
)

func (s *Storage) GetAllUsersWithLevel(ctx context.Context, lvl model.AccessLevel) ([]*model.UserInfo, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, login, access_level FROM Users WHERE access_level = $1", lvl)
	if err != nil {
		return nil, errwrap.Wrap(storage.ErrAllUsersWithLevel, err)
	}
	defer rows.Close()

	users := []*model.UserInfo{}

	for rows.Next() {
		user := &model.UserInfo{}

		err := rows.Scan(&user.ID, &user.Login, &user.AccessLvl)
		if err != nil {
			return nil, errwrap.Wrap(storage.ErrAllUsersWithLevel, err)
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, errwrap.Wrap(storage.ErrAllUsersWithLevel, err)
	}

	return users, nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	cmdTag, err := s.pool.Exec(ctx, "DELETE FROM Users WHERE id = $1", userID)
	if err != nil {
		return errwrap.Wrap(storage.ErrDeleteUser, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(storage.ErrDeleteUser, storage.ErrNoRowsAffected)
	}

	return nil
}

func (s *Storage) SetAccessLevel(ctx context.Context, userID uuid.UUID, lvl model.AccessLevel) error {
	cmdTag, err := s.pool.Exec(ctx, "UPDATE Users SET access_level = $1 WHERE id = $2", lvl, userID)
	if err != nil {
		return errwrap.Wrap(storage.ErrUpdatePassword, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errwrap.Wrap(storage.ErrUpdatePassword, storage.ErrNoRowsAffected)
	}

	return nil
}
