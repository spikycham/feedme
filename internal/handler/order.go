package handler

import (
	"net/http"

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

type ResponseOrderList struct {
	Data []repository.Order `json:"data"`
}

func (h *OrderHandler) GetOrderList(w http.ResponseWriter, r *http.Request) error {
	orders, err := h.r.SelectAllOrders(r.Context())
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.Write(w, &ResponseOrderList{Data: orders})
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

	if err := h.r.UpdateOrderStatusByOrderID(r.Context(), body.OrderID, repository.OrderStatus(body.Status)); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	return nil
}
