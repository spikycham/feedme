package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/model"
	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
	"github.com/spikycham/feedme/pkg/random"
	"github.com/spikycham/feedme/pkg/token"
)

// The initial structure and interface.
type AuthHandler struct {
	r *repository.AuthRepository
}

func NewAuthHandler(r *repository.AuthRepository) AuthHandler {
	return AuthHandler{r}
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
		UserID               string         `json:"user_id"`
		Name                 string         `json:"name"`
		Account              string         `json:"account"`
		Role                 model.UserRole `json:"role"`
		AvatarURI            string         `json:"avatar_uri"`
		ProfileBackgroundURI string         `json:"profile_background_uri"`
		CreatedAt            int64          `json:"created_at"`
	}
	ResponseLogin struct {
		Token ResponseToken `json:"token"`
		User  ResponseUser  `json:"user"`
	}
)

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
	// Compare the password after hashing with a secret key from the environment.
	secret := os.Getenv("HMAC_SECRET_KEY")
	if secret == "" {
		network.Error(w, http.StatusInternalServerError)
		return constant.MissingImportantEnv
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body.Password))
	macPwd := hex.EncodeToString(mac.Sum(nil))

	if strings.Compare(macPwd, user.Password) != 0 {
		// We cannot use 401 as we need to validate
		// the 401 and resend request in the frontend.
		network.Error(w, http.StatusInternalServerError)
		return fmt.Errorf("account: %s logged with incorrect password", body.Account)
	}

	// Generate a new access token and a refresh token.
	at, err := token.Sign(user.UserID)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}
	// Store the refresh token to sqlite database.
	// Use sqlite is already enough for this project.
	rt, err := random.RandBase64(32)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	expiredAt := time.Now().Add(7 * 24 * time.Hour).Unix()
	if err := h.r.InsertRefreshToken(r.Context(), user.UserID, rt, expiredAt); err != nil {
		network.Error(w, http.StatusUnauthorized)
		return err
	}

	network.Write(w, &ResponseLogin{
		Token: ResponseToken{
			AccessToken:  at,
			RefreshToken: rt,
		},
		User: ResponseUser{
			UserID:               user.UserID,
			Name:                 user.Name,
			Account:              user.Account,
			Role:                 user.Role,
			AvatarURI:            user.AvatarURI,
			ProfileBackgroundURI: user.ProfileBackgroundURI,
			CreatedAt:            user.CreatedAt,
		},
	})
	return nil
}

// TODO: implement this.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) error {
	var body RequestRefreshToken
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	// Validate the refresh token and get the user_id from the same table.
	row, err := h.r.GetExpiredAtByRefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	if row.ExpiredAt < time.Now().Unix() {
		network.Error(w, http.StatusUnauthorized)
		return constant.RefreshTokenExpired
	}

	// Generate a new access token and a refresh token.
	at, err := token.Sign(row.UserID)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	rt, err := random.RandBase64(32)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	// Insert the new refresh token to databse.
	expiredAt := time.Now().Add(7 * 24 * time.Hour).Unix()
	if err := h.r.InsertRefreshToken(r.Context(), row.UserID, rt, expiredAt); err != nil {
		return err
	}

	network.Write(w, &ResponseToken{
		AccessToken:  at,
		RefreshToken: rt,
	})
	return nil
}
