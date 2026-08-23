package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
	"github.com/spikycham/feedme/pkg/token"
)

func NewAuthHandler(r AuthRepository) AuthHandler {
	return AuthHandler{r: r}
}

// Handlers type structs.
type (
	// Login.
	RequestLogin struct {
		Account  string `json:"account" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	ResponseToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	ResponseUser struct {
		UserID    string              `json:"user_id"`
		Name      string              `json:"name"`
		Account   string              `json:"account"`
		Role      repository.UserRole `json:"role"`
		AvatarURI *string             `json:"avatar_uri"`
	}
	ResponseLogin struct {
		Token ResponseToken `json:"token"`
		User  ResponseUser  `json:"user"`
	}
)

type AuthRepository interface {
	GetUserByAccount(ctx context.Context, account string) (*repository.User, error)
}

type AuthHandler struct {
	r AuthRepository
}

// Response the token and the user profiles.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var body RequestLogin
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	user, err := h.r.GetUserByAccount(r.Context(), body.Account)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	// Validate the password.
	// TODO: set the password to hashed string.
	// TODO: parse the password from the md5 hashed password in request body.
	if strings.Compare(body.Password, user.Password) != 0 {
		network.Error(w, http.StatusUnauthorized)
		return fmt.Errorf("account: %s logged with incorrect password", body.Account)
	}

	// Generate a new access token and a refresh token.
	at, err := token.Sign(user.UserID)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}
	// TODO: store the refresh token to postgres database.
	// Use psql is already enough for this project.
	rt, err := token.RandBase64(32)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.Write(w, &ResponseLogin{
		Token: ResponseToken{
			AccessToken:  at,
			RefreshToken: rt,
		},
		User: ResponseUser{
			UserID:    user.UserID,
			Name:      user.Name,
			Account:   user.Account,
			Role:      user.Role,
			AvatarURI: user.AvatarURI,
		},
	})
	return nil
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	return nil
}
