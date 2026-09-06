package handler

import (
	"net/http"

	"github.com/spikycham/feedme/internal/model"
	"github.com/spikycham/feedme/internal/repository"
	"github.com/spikycham/feedme/pkg/network"
	"github.com/spikycham/feedme/pkg/random"
)

type OrderHandler struct {
	r *repository.OrderRepository
}

func NewOrderHandler(r *repository.OrderRepository) *OrderHandler {
	return &OrderHandler{r}
}

type (
	OrderListItem struct {
		OrderID          string            `json:"order_id"`
		Status           model.OrderStatus `json:"status"`
		Amount           float64           `json:"amoun"`
		CreatedAt        int64             `json:"created_at"`
		DoneAt           int64             `json:"done_at"`
		Comment          string            `json:"comment"`
		CommentedAt      int64             `json:"commented_at"`
		CommentDeletedAt int64             `json:"comment_deleted_at"`
	}
	ResponseOrderList struct {
		List []OrderListItem `json:"list"`
	}
)

func (h *OrderHandler) GetOrderList(w http.ResponseWriter, r *http.Request) error {
	orders, err := h.r.SelectAllOrders(r.Context())
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	resp := make([]OrderListItem, 0)
	for _, o := range orders {
		resp = append(resp, OrderListItem{
			OrderID:          o.OrderID,
			Status:           o.Status,
			Amount:           o.Amount,
			CreatedAt:        o.CreatedAt,
			DoneAt:           o.DoneAt,
			Comment:          o.Comment,
			CommentedAt:      o.CommentedAt,
			CommentDeletedAt: o.CommentDeletedAt,
		})
	}

	network.Write(w, &ResponseOrderList{List: resp})
	return nil
}

type (
	CreateOrderFood struct {
		FoodID string `json:"food_id"`
		Count  int    `json:"count"`
	}
	RequestCreateOrder struct {
		Amount *float64          `json:"amount" validate:"required"`
		Foods  []CreateOrderFood `json:"foods" validate:"required"`
	}
)

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	var body RequestCreateOrder
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	id, err := random.RandID()
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	p := &repository.InsertOrderParams{
		OrderID: id,
		Amount:  *body.Amount,
	}
	for _, f := range body.Foods {
		p.Foods = append(p.Foods, repository.InsertOrderFoodParams{
			FoodID: f.FoodID,
			Count:  f.Count,
		})
	}
	if err := h.r.InsertOrder(r.Context(), p); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusCreated)
	return nil
}

type RequestUpdateOrderStatus struct {
	OrderID string `json:"order_id" validate:"required"`
	Status  int    `json:"status" validate:"required"`
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) error {
	var body RequestUpdateOrderStatus
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if err := h.r.UpdateOrderStatusByOrderID(r.Context(), body.OrderID, model.OrderStatus(body.Status)); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusOK)
	return nil
}

type RequestCreateOrderComment struct {
	OrderID string `json:"order_id" validate:"required"`
	Detail  string `json:"detail" validate:"required"`
}

func (h *OrderHandler) CreateOrderComment(w http.ResponseWriter, r *http.Request) error {
	var body RequestCreateOrderComment
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if err := h.r.UpdateOrderCommentByOrderID(r.Context(), body.OrderID, body.Detail); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusCreated)
	return nil
}

type RequestDeleteOrderComment struct {
	OrderID string `json:"order_id"`
}

func (h *OrderHandler) DeleteOrderComment(w http.ResponseWriter, r *http.Request) error {
	var body RequestDeleteOrderComment
	if err := network.Read(r, &body); err != nil {
		network.Error(w, http.StatusBadRequest)
		return err
	}

	if err := h.r.DeleteOrderCommentByOrderID(r.Context(), body.OrderID); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.WriteEmpty(w, http.StatusNoContent)
	return nil
}
