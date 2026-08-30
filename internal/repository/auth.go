package repository

import (
	"context"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db}
}

type (
	UserRole int

	User struct {
		ID                   int
		UserID               string
		Name                 string
		Account              string
		Password             string
		Role                 UserRole // 0 customer, 1 merchant
		AvatarURI            string
		ProfileBackgroundURI string
		CreatedAt            int64
	}
)

const (
	UserRoleCustomer UserRole = iota
	UserRoleMerchant
)

func (r *AuthRepository) GetUserByAccount(ctx context.Context, account string) (*User, error) {
	var user User
	if err := r.db.QueryRowContext(ctx, "SELECT id, user_id, name, account, password, role, avatar_uri, profile_background_uri, created_at FROM users WHERE account = ?", account).Scan(&user.ID,
		&user.UserID, &user.Name, &user.Account, &user.Password, &user.Role, &user.AvatarURI, &user.ProfileBackgroundURI, &user.CreatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
