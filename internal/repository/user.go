package repository

import (
	"context"
	"database/sql"

	"github.com/spikycham/feedme/internal/constant"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// ID        int
// UserID    string
// Name      string
// Account   string
// Password  string
// Role      UserRole
// AvatarURI *string
func (r *UserRepository) GetUserByUserID(ctx context.Context, userId string) (*User, error) {
	var user User
	if err := r.db.QueryRowContext(ctx, "SELECT user_id, name, account, role, avatar_uri FROM users WHERE user_id = $1", userId).Scan(
		&user.UserID,
		&user.Name,
		&user.Account,
		&user.Role,
		&user.AvatarURI,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

// Upload the image to the assets then change the avatar uri.
func (r *UserRepository) UpdateUserAvatarByUserID(ctx context.Context, userId, newUri string) error {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET avatar_uri = $1 WHERE user_id = $2", newUri, userId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}

func (r *UserRepository) UpdateUsernameByUserID(ctx context.Context, userId, newName string) error {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET name = $1 WHERE user_id = $2", newName, userId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}

func (r *UserRepository) UpdateUserPasswordByUserID(ctx context.Context, userId, newPassword string) error {
	res, err := r.db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE user_id = $2", newPassword, userId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}
