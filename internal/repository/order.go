package repository

import (
	"context"
	"database/sql"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/internal/model"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db}
}

func (r *OrderRepository) SelectAllOrders(ctx context.Context) ([]model.Order, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT order_id, status, amount, created_at, done_at, comment, commented_at, comment_deleted_at FROM orders")
	if err != nil {
		return nil, err
	}

	orders := make([]model.Order, 0)
	for rows.Next() {
		var order model.Order
		if err := rows.Scan(
			&order.OrderID,
			&order.Status,
			&order.Amount,
			&order.CreatedAt,
			&order.DoneAt,
			&order.Comment,
			&order.CommentedAt,
			&order.CommentDeletedAt,
		); err != nil {
			return nil, err
		}
		order.ID = -1

		orders = append(orders, order)
	}

	return orders, nil
}

type InsertOrderFoodParams struct {
	FoodID string
	Count  int
}
type InsertOrderParams struct {
	OrderID string
	Amount  float64
	Foods   []InsertOrderFoodParams
}

func (r *OrderRepository) InsertOrder(ctx context.Context, p *InsertOrderParams) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert into orders and order_foods tables.
	if _, err := tx.ExecContext(ctx, "INSERT INTO orders (order_id, amount) VALUES (?, ?)", p.OrderID, p.Amount); err != nil {
		return err
	}
	for _, f := range p.Foods {
		// Validate the existance of foods.
		var isExist bool
		if err := tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM foods WHERE food_id = ?)", f.FoodID).Scan(&isExist); err != nil {
			return err
		}
		if !isExist {
			return constant.NoDataFound
		}

		// Insert the real record.
		if _, err := tx.ExecContext(ctx, "INSERT INTO order_foods (order_id, food_id, food_count) VALUES (?, ?, ?)", p.OrderID, f.FoodID, f.Count); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrderStatusByOrderID(ctx context.Context, orderId string, status model.OrderStatus) error {
	res, err := r.db.ExecContext(ctx, "UPDATE orders SET status = ?, done_at = (unixepoch()) WHERE order_id = ? AND status = 0", status, orderId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}

func (r *OrderRepository) UpdateOrderCommentByOrderID(ctx context.Context, orderId, detail string) error {
	res, err := r.db.ExecContext(
		ctx,
		"UPDATE orders SET comment = ?, commented_at = (unixepoch()) WHERE order_id = ? AND commented_at = -1 AND status = 2",
		detail,
		orderId,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}

func (r *OrderRepository) DeleteOrderCommentByOrderID(ctx context.Context, orderId string) error {
	res, err := r.db.ExecContext(ctx, "UPDATE orders SET comment_deleted_at = (unixepoch()) WHERE order_id = ? AND comment_deleted_at = -1", orderId)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return constant.NoAffectedRows
	}

	return nil
}
