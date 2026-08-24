package repository

import (
	"context"
	"database/sql"
)

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

type (
	UserRole int

	User struct {
		ID        int
		UserID    string
		Name      string
		Account   string
		Password  string
		Role      UserRole
		AvatarURI *string
	}
)

const (
	RoleCustomer UserRole = iota
	RoleMerchant
)

type AuthRepository struct {
	db *sql.DB
}

func (r *AuthRepository) GetUserByAccount(ctx context.Context, account string) (*User, error) {
	var user User
	if err := r.db.QueryRowContext(ctx, "SELECT id, user_id, name, account, password, role, avatar_uri FROM users WHERE account = $1", account).Scan(&user.ID,
		&user.UserID, &user.Name, &user.Account, &user.Password, &user.Role, &user.AvatarURI); err != nil {
		return nil, err
	}
	return &user, nil
}
