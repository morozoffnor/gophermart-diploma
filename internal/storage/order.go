package storage

import (
	"context"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
	StatusInvalid    = "INVALID"
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

func (db *DB) GetOrdersList(ctx context.Context, userID string) ([]*Order, error) {
	query := "SELECT id, status, accrual, uploaded_at FROM orders WHERE user_id=$1"
	rows, err := db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []*Order
	for rows.Next() {
		var order Order
		order.UserID = userID
		err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, &order)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}

func (db *DB) AddOrder(ctx context.Context, userID string, orderID string) error {
	query := "INSERT INTO orders (id, status, user_id, accrual) VALUES ($1, $2, $3, $4)"
	_, err := db.pool.Exec(ctx, query, orderID, StatusNew, userID, 0)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) UpdateOrderFromAccrual(ctx context.Context, orderID string, status string, accrual float64) error {
	query := "UPDATE orders SET status = $1, accrual = $2 where id = $3"
	_, err := db.pool.Exec(ctx, query, status, accrual, orderID)
	if err != nil {
		return err
	}
	return nil
}
