package handler

import (
	"net/http"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/model"
	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
	"github.com/spikycham/feedme/pkg/random"
)

type FoodHandler struct {
	r *repository.FoodRepository
}

func NewFoodHandler(r *repository.FoodRepository) *FoodHandler {
	return &FoodHandler{r}
}

// Handlers.
// Get the food list.
type (
	FoodListItemStep struct {
		Sort   int    `json:"sort"`
		Detail string `json:"detail"`
	}
	FoodListItemComment struct {
		CommentID string `json:"comment_id"`
		Detail    string `json:"detail"`
		CreatedAt int64  `json:"created_at"`
		DeletedAt int64  `json:"deleted_at"`
	}
	FoodListItem struct {
		FoodID       string                `json:"food_id"`
		Name         string                `json:"name"`
		Detail       string                `json:"detail"`
		Prize        float32               `json:"prize"`
		Rate         float32               `json:"rate"`
		RequiredTime int64                 `json:"required_time"`
		SoldCount    int                   `json:"sold_count"`
		ImageURIs    []string              `json:"image_uris"`
		Ingredients  []string              `json:"ingredients"`
		Category     model.FoodCategory    `json:"category"`
		CreatedAt    int64                 `json:"created_at"`
		DeletedAt    int64                 `json:"deleted_at"`
		Steps        []FoodListItemStep    `json:"steps"`
		Comments     []FoodListItemComment `json:"comments"`
		// 0 staple food, 1 vegetable, 2 meat, 3 seafood, 4 soup, 5 dessert, 6 drink, 7 other
	}
	ResponseFoodList struct {
		List []FoodListItem `json:"list"`
	}
)

func (h *FoodHandler) GetFoodList(w http.ResponseWriter, r *http.Request) error {
	foods, err := h.r.SelectAllFoods(r.Context())
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	resp := make([]FoodListItem, 0)
	for _, f := range foods {
		item := FoodListItem{
			FoodID:       f.Food.FoodID,
			Name:         f.Food.Name,
			Detail:       f.Food.Detail,
			Prize:        f.Food.Prize,
			Rate:         f.Food.Rate,
			RequiredTime: f.Food.RequiredTime,
			SoldCount:    f.Food.SoldCount,
			ImageURIs:    f.Food.ImageURIs,
			Ingredients:  f.Food.Ingredients,
			Category:     f.Food.Category,
			CreatedAt:    f.Food.CreatedAt,
			DeletedAt:    f.Food.DeletedAt,
		}

		item.Steps = make([]FoodListItemStep, 0)
		for _, s := range f.Steps {
			item.Steps = append(item.Steps, FoodListItemStep{
				Sort:   s.Sort,
				Detail: s.Detail,
			})
		}

		item.Comments = make([]FoodListItemComment, 0)
		for _, c := range f.Comments {
			item.Comments = append(item.Comments, FoodListItemComment{
				CommentID: c.CommentID,
				Detail:    c.Detail,
				CreatedAt: c.CreatedAt,
				DeletedAt: c.DeletedAt,
			})
		}

		resp = append(resp, item)
	}

	network.Write(w, &ResponseFoodList{List: resp})
	return nil
}

// Create a food handler.
type RequestCreateFood struct {
	Name         string              `json:"name" validate:"required"`
	Detail       string              `json:"detail" validate:"required"`
	Prize        *float32            `json:"prize" validate:"required"`
	Rate         *float32            `json:"rate" validate:"required"`
	RequiredTime int64               `json:"required_time" validate:"required"`
	ImageURIs    []string            `json:"image_uris" validate:"required"`
	Ingredients  []string            `json:"ingredients" validate:"required"`
	Category     *model.FoodCategory `json:"category" validate:"required"`
	Steps        []FoodListItemStep  `json:"steps" validate:"required"`
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
	if len(body.ImageURIs) > 3 {
		network.Error(w, http.StatusBadRequest)
		return constant.ErrOutOfRange
	}

	fid, err := random.RandID()
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	p := &repository.InsertFoodParams{
		FoodID:       fid,
		Name:         body.Name,
		Detail:       body.Detail,
		Prize:        *body.Prize,
		Rate:         *body.Rate,
		RequiredTime: body.RequiredTime,
		ImageURIs:    body.ImageURIs,
		Ingredients:  body.Ingredients,
		Category:     *body.Category,
	}

	for _, s := range body.Steps {
		p.Steps = append(p.Steps, model.FoodStep{
			FoodID: fid,
			Sort:   s.Sort,
			Detail: s.Detail,
		})
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
	FoodID       string              `json:"food_id" validate:"required"`
	Name         string              `json:"name"`
	Detail       string              `json:"detail"`
	Prize        *float32            `json:"prize"`
	Rate         *float32            `json:"rate"`
	RequiredTime int64               `json:"required_time"`
	ImageURIs    []string            `json:"image_uris"`
	Ingredients  []string            `json:"ingredients"`
	Category     *model.FoodCategory `json:"category"`
}

func (h *FoodHandler) UpdateFood(w http.ResponseWriter, r *http.Request) error {
	var body RequestUpdateFood
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	// Validate the bussiness logic of the prize and rate.
	if body.Prize != nil {
		if *body.Prize < 0 || *body.Prize > 9999.99 {
			network.Error(w, http.StatusBadRequest)
			return constant.ErrOutOfRange
		}
	}
	if body.Rate != nil {
		if *body.Rate < 0 || *body.Rate > 5 {
			network.Error(w, http.StatusBadRequest)
			return constant.ErrOutOfRange
		}
	}
	if len(body.ImageURIs) > 3 {
		network.Error(w, http.StatusBadRequest)
		return constant.ErrOutOfRange
	}

	p := &repository.UpdateFoodParams{
		Name:         &body.Name,
		Detail:       &body.Detail,
		Prize:        body.Prize,
		Rate:         body.Rate,
		RequiredTime: &body.RequiredTime,
		ImageURIs:    body.ImageURIs,
		Ingredients:  body.Ingredients,
		Category:     body.Category,
	}
	if err := h.r.UpdateFoodByFoodID(r.Context(), body.FoodID, p); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusOK)
	return nil
}

// Update a food steps handler.
type RequestUpdateFoodStep struct {
	FoodID string             `json:"food_id" validate:"required"`
	Steps  []FoodListItemStep `json:"steps" validate:"required"`
}

func (h *FoodHandler) UpdateFoodSteps(w http.ResponseWriter, r *http.Request) error {
	var body RequestUpdateFoodStep
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	steps := make([]model.FoodStep, 0)
	for _, s := range body.Steps {
		steps = append(steps, model.FoodStep{
			Sort:   s.Sort,
			Detail: s.Detail,
		})
	}

	if err := h.r.UpdateFoodStepsByFoodID(r.Context(), body.FoodID, steps); err != nil {
		return err
	}

	network.WriteEmpty(w, http.StatusCreated)
	return nil
}

// Delete food handler.
type RequestDeleteFood struct {
	FoodID string `json:"food_id" validate:"required"`
}

// NOTE: it is not necessary to delete the food steps until now.
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

// Food comment handlers.
type RequestCreateFoodComment struct {
	FoodID string `json:"food_id" validate:"required"`
	Detail string `json:"detail" validate:"required"`
}

// Only allowed to create a food comment when the user ordered the food in a order in the past.
func (h *FoodHandler) CreateFoodComment(w http.ResponseWriter, r *http.Request) error {
	var body RequestCreateFoodComment
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	cid, err := random.RandID()
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	if err := h.r.InsertFoodCommentByFoodID(r.Context(), cid, body.FoodID, body.Detail); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusCreated)
	return nil
}

type RequestDeleteFoodComment struct {
	CommentID string `json:"comment_id" validate:"required"`
}

func (h *FoodHandler) DeleteFoodComment(w http.ResponseWriter, r *http.Request) error {
	var body RequestDeleteFoodComment
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if err := h.r.DeleteFoodCommentByCommentID(r.Context(), body.CommentID); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusNoContent)
	return nil
}
