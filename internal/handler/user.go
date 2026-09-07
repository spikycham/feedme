package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/middleware"
	"github.com/spikycham/feedme/internal/model"
	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
)

// The initial structure and interface.

type UserHandler struct {
	r *repository.UserRepository
}

func NewUserHandler(r *repository.UserRepository) *UserHandler {
	return &UserHandler{r: r}
}

// Handlers type structs.
type (
	// Me.
	ResponseMe struct {
		UserID               string         `json:"user_id"`
		Name                 string         `json:"name"`
		Account              string         `json:"account"`
		Role                 model.UserRole `json:"role"`
		AvatarURI            string         `json:"avatar_uri"`
		ProfileBackgroundURI string         `json:"profile_background_uri"`
		CreatedAt            int64          `json:"created_at"`
	}
	// Change avatar uri, username, password.
	RequestChangeProfile struct {
		NewAvatarURI *string `json:"new_avatar_uri"`
		NewUsername  *string `json:"new_username" validate:"min=2, max=10"`
		NewPassword  *string `json:"new_password" validate:"min=6, max=20"`
	}
)

// Response the user profiles through the token.
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) error {
	id := r.Context().Value(middleware.USER_ID_KEY).(string)

	me, err := h.r.GetUserByUserID(r.Context(), id)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	if err := network.Write(w, &ResponseMe{
		UserID:               me.UserID,
		Name:                 me.Name,
		Account:              me.Account,
		Role:                 me.Role,
		AvatarURI:            me.AvatarURI,
		ProfileBackgroundURI: me.ProfileBackgroundURI,
		CreatedAt:            me.CreatedAt,
	}); err != nil {
		return err
	}

	return nil
}

// Change avatar, username and password.
func (h *UserHandler) ChangeProfile(w http.ResponseWriter, r *http.Request) error {
	var body RequestChangeProfile
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if body.NewAvatarURI != nil && !strings.HasPrefix(strings.ToLower(*body.NewAvatarURI), "https://") {
		network.Error(w, http.StatusBadRequest)
		return constant.InvalidURIPrefix
	}

	userId, ok := r.Context().Value(middleware.USER_ID_KEY).(string)
	if !ok {
		network.Error(w, http.StatusInternalServerError)
		return constant.InvalidParsedToken
	}

	// Hash the password.
	secret := os.Getenv("HMAC_SECRET_KEY")
	if secret == "" {
		network.Error(w, http.StatusInternalServerError)
		return constant.MissingImportantEnv
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(*body.NewPassword))
	macPwd := hex.EncodeToString(mac.Sum(nil))

	if err := h.r.UpdateUserProfileByUserID(r.Context(), userId, body.NewAvatarURI, body.NewUsername, &macPwd); err != nil {
		// 404
		if errors.Is(err, constant.NoAffectedRows) {
			network.WriteEmpty(w, http.StatusNotFound)
			return err
		}
		return err
	}

	// 204
	network.WriteEmpty(w, http.StatusNoContent)
	return nil
}
