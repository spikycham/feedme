package repository

import (
	"context"
	"database/sql"

	"github.com/spikycham/feedme/internal/model"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) GetUserByAccount(ctx context.Context, account string) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRowContext(ctx, "SELECT id, user_id, name, account, password, role, avatar_uri, profile_background_uri, created_at FROM users WHERE account = ?", account).Scan(&user.ID,
		&user.UserID, &user.Name, &user.Account, &user.Password, &user.Role, &user.AvatarURI, &user.ProfileBackgroundURI, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) GetExpiredAtByRefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenRow, error) {
	var row model.RefreshTokenRow
	if err := r.db.QueryRowContext(
		ctx,
		"SELECT id, user_id, token, created_at, expired_at FROM tokens WHERE token = ?",
		refreshToken,
	).Scan(
		&row.ID,
		&row.UserID,
		&row.Token,
		&row.CreatedAt,
		&row.ExpiredAt,
	); err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *AuthRepository) InsertRefreshToken(ctx context.Context, userId, token string, expiredAt int64) error {
	if _, err := r.db.ExecContext(
		ctx,
		"INSERT INTO tokens (user_id, token, expired_at) VALUES (?, ?, ?)",
		userId,
		token,
		expiredAt,
	); err != nil {
		return err
	}
	return nil
}
