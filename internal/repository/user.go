package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetUserByUserID(ctx context.Context, userId string) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRowContext(ctx, "SELECT user_id, name, account, role, avatar_uri, profile_background_uri, created_at FROM users WHERE user_id = ?", userId).Scan(
		&user.UserID,
		&user.Name,
		&user.Account,
		&user.Role,
		&user.AvatarURI,
		&user.ProfileBackgroundURI,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserProfileByUserID(ctx context.Context, userId string, newAvatarUri, newUsername, newPassword, newProfileBackgroundUri *string) error {
	var sets []string
	var args []any

	if newAvatarUri != nil {
		sets = append(sets, "avatar_uri = ?")
		args = append(args, *newAvatarUri)
	}
	if newUsername != nil {
		sets = append(sets, "name = ?")
		args = append(args, *newUsername)
	}
	if newPassword != nil {
		sets = append(sets, "password = ?")
		args = append(args, *newPassword)
	}
	if newProfileBackgroundUri != nil {
		sets = append(sets, "profile_background_uri = ?")
		args = append(args, *newProfileBackgroundUri)
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, userId)
	query := fmt.Sprintf("UPDATE users SET %s WHERE user_id = ?", strings.Join(sets, ", "))

	res, err := r.db.ExecContext(ctx, query, args...)
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
