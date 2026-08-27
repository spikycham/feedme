package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/middleware"
	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
)

// The initial structure and interface.
type UserRepository interface {
	GetUserByUserID(ctx context.Context, userId string) (*repository.User, error)
	UpdateUserProfileByUserID(ctx context.Context, userId string, newAvatarUri, newUsername, newPassword *string) error
}

type UserHandler struct {
	r UserRepository
}

func NewUserHandler(r UserRepository) *UserHandler {
	return &UserHandler{r: r}
}

// Handlers type structs.
type (
	// Me.
	ResponseMe struct {
		UserID               string              `json:"user_id"`
		Name                 string              `json:"name"`
		Account              string              `json:"account"`
		Role                 repository.UserRole `json:"role"`
		AvatarURI            *string             `json:"avatar_uri"`
		ProfileBackgroundURI *string             `json:"profile_background_uri"`
		CreatedAt            int64               `json:"created_at"`
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
		network.Error(w, http.StatusUnauthorized)
		return constant.InvalidParsedToken
	}

	// TODO: hash the password and concate the salt to the hashed password string.
	if err := h.r.UpdateUserProfileByUserID(r.Context(), userId, body.NewAvatarURI, body.NewUsername, body.NewPassword); err != nil {
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
