package storage

import (
	"context"
)

const (
	StatusNew = iota
	StatusProcessing
	StatusProcessed
	StatusInvalid
)

func (db *DB) OrderExists(ctx context.Context, orderID string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)"
	var exists bool
	err := db.pool.QueryRow(ctx, query, orderID).Scan(&exists)
	if err != nil {
		return true, err
	}
	return exists, err
}

func (db *DB) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	query := "SELECT status, accrual, user_id, uploaded_at FROM orders WHERE id=$1"
	order := &Order{
		Number: orderID,
	}
	err := db.pool.QueryRow(ctx, query, orderID).Scan(&order.Status, &order.Accrual, &order.UserID, &order.UploadedAt)
	if err != nil {
		return nil, err
	}
	return order, nil
}
func (db *DB) AddOrder(ctx context.Context, userID string, orderID string) error {
	query := "INSERT INTO orders (id, status, user_id, accrual) VALUES ($1, $2, $3)"
	_, err := db.pool.Exec(ctx, query, orderID, StatusNew, userID, 0)
	if err != nil {
		return err
	}
	return nil
}
