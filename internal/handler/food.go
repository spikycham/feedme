package handler

import (
	"context"
	"net/http"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
	"github.com/spikycham/feedme/pkg/random"
)

type FoodRepository interface {
	SelectAllFoods(ctx context.Context) error
	InsertFood(ctx context.Context, p repository.InsertFoodParams) error
	UpdateFoodByFoodID(ctx context.Context, foodId string, p *repository.UpdateFoodParams) error
	DeleteFoodByFoodID(ctx context.Context, foodId string) error
}

type FoodHandler struct {
	r *repository.FoodRepository
}

func NewFoodHandler(r *repository.FoodRepository) *FoodHandler {
	return &FoodHandler{r}
}

// Handlers.
// Get the food list.
type ResponseFoodList struct {
	Data []repository.Food `json:"data"`
}

func (h *FoodHandler) GetFoodList(w http.ResponseWriter, r *http.Request) error {
	foods, err := h.r.SelectAllFoods(r.Context())
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.Write(w, &ResponseFoodList{Data: foods})
	return nil
}

// Create a food handler.
type RequestCreateFood struct {
	Name         string                   `json:"name" validate:"required"`
	Detail       string                   `json:"detail" validate:"required"`
	Prize        *float32                 `json:"prize" validate:"required"`
	Rate         *float32                 `json:"rate" validate:"required"`
	RequiredTime int64                    `json:"required_time" validate:"required"`
	ImageURIs    []string                 `json:"image_uris" validate:"required"`
	Category     *repository.FoodCategory `json:"category" validate:"required"`
	// 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
}

func (h *FoodHandler) CreateFood(w http.ResponseWriter, r *http.Request) error {
	var body RequestCreateFood
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	// Validate the bussiness logic of the prize and rate.
	if *body.Prize < 0 || *body.Prize > 9999.99 {
		network.Error(w, http.StatusBadRequest)
		return constant.ErrOutOfRange
	}
	if *body.Rate < 0 || *body.Rate > 5 {
		network.Error(w, http.StatusBadRequest)
		return constant.ErrOutOfRange
	}

	id, err := random.RandID()
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	p := &repository.InsertFoodParams{
		FoodID:       id,
		Name:         body.Name,
		Detail:       body.Detail,
		Prize:        *body.Prize,
		Rate:         *body.Rate,
		RequiredTime: body.RequiredTime,
		ImageURIs:    body.ImageURIs,
		Category:     *body.Category,
	}
	if err := h.r.InsertFood(r.Context(), p); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusCreated)
	return nil
}

// Update a food handler.
type RequestUpdateFood struct {
	FoodID       string                   `json:"food_id"`
	Name         string                   `json:"name"`
	Detail       string                   `json:"detail"`
	Prize        *float32                 `json:"prize"`
	Rate         *float32                 `json:"rate"`
	RequiredTime int64                    `json:"required_time"`
	ImageURIs    []string                 `json:"image_uris"`
	Category     *repository.FoodCategory `json:"category"`
}

func (h *FoodHandler) UpdateFood(w http.ResponseWriter, r *http.Request) error {
	var body RequestUpdateFood
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	// Validate the bussiness logic of the prize and rate.
	if *body.Prize < 0 || *body.Prize > 9999.99 {
		network.Error(w, http.StatusBadRequest)
		return constant.ErrOutOfRange
	}
	if *body.Rate < 0 || *body.Rate > 5 {
		network.Error(w, http.StatusBadRequest)
		return constant.ErrOutOfRange
	}

	p := &repository.UpdateFoodParams{
		Name:         &body.Name,
		Detail:       &body.Detail,
		Prize:        body.Prize,
		Rate:         body.Rate,
		RequiredTime: &body.RequiredTime,
		ImageURIs:    &body.ImageURIs,
		Category:     body.Category,
	}
	if err := h.r.UpdateFoodByFoodID(r.Context(), body.FoodID, p); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusOK)
	return nil
}

// Delete food handler.
type RequestDeleteFood struct {
	FoodID string `json:"food_id" validate:"required"`
}

func (h *FoodHandler) DeleteFood(w http.ResponseWriter, r *http.Request) error {
	var body RequestDeleteFood
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if err := h.r.DeleteFoodByFoodID(r.Context(), body.FoodID); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusNoContent)
	return nil
}
