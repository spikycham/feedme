package model

type OrderStatus int

const (
	OrderStatusPending = iota
	OrderStatusRejected
	OrderStatusDone
)

type Order struct {
	ID               int
	OrderID          string
	Status           OrderStatus
	Amount           float64
	CreatedAt        int64
	DoneAt           int64
	Comment          string
	CommentedAt      int64
	CommentDeletedAt int64
}
