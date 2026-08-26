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
	UpdateUserAvatarByUserID(ctx context.Context, userId, newUri string) error
	UpdateUsernameByUserID(ctx context.Context, userId, newName string) error
	UpdateUserPasswordByUserID(ctx context.Context, userId, newPassword string) error
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
		UserID    string              `json:"user_id"`
		Name      string              `json:"name"`
		Account   string              `json:"account"`
		Role      repository.UserRole `json:"role"`
		AvatarURI *string             `json:"avatar_uri"`
	}
	// Change avatar uri.
	RequestChangeAvatar struct {
		NewURI string `json:"new_uri" validate:"required"`
	}
	// Change username.
	RequestChangeUsername struct {
		NewName string `json:"new_name" validate:"required, min=2, max=10"`
	}
	// Change password.
	RequestChangePassword struct {
		NewPassword string `json:"new_password" validate:"required, min=6, max=20"`
	}
)

// Response the user profiles through the token.
func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) error {
	id := r.URL.Query().Get("user_id")

	me, err := h.r.GetUserByUserID(r.Context(), id)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	if err := network.Write(w, &ResponseMe{
		UserID:    me.UserID,
		Name:      me.Name,
		Account:   me.Account,
		Role:      me.Role,
		AvatarURI: me.AvatarURI,
	}); err != nil {
		return err
	}

	return nil
}

// Change avatar, username and password.
func (h *UserHandler) ChangeAvatarURI(w http.ResponseWriter, r *http.Request) error {
	var body RequestChangeAvatar
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if !strings.HasPrefix(strings.ToLower(body.NewURI), "https://") {
		network.Error(w, http.StatusBadRequest)
		return constant.InvalidURIPrefix
	}

	userId, ok := r.Context().Value(middleware.USER_ID_KEY).(string)
	if !ok {
		network.Error(w, http.StatusUnauthorized)
		return constant.InvalidParsedToken
	}

	if err := h.r.UpdateUserAvatarByUserID(r.Context(), userId, body.NewURI); err != nil {
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

func (h *UserHandler) ChangeUsername(w http.ResponseWriter, r *http.Request) error {
	var body RequestChangeUsername
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	userId, ok := r.Context().Value(middleware.USER_ID_KEY).(string)
	if !ok {
		network.Error(w, http.StatusUnauthorized)
		return constant.InvalidParsedToken
	}

	if err := h.r.UpdateUsernameByUserID(r.Context(), userId, body.NewName); err != nil {
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

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) error {
	var body RequestChangePassword
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	userId, ok := r.Context().Value(middleware.USER_ID_KEY).(string)
	if !ok {
		network.Error(w, http.StatusUnauthorized)
		return constant.InvalidParsedToken
	}

	// TODO: hash the password and concate the salt to the hashed password string.

	if err := h.r.UpdateUserPasswordByUserID(r.Context(), userId, body.NewPassword); err != nil {
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
