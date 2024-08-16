package storage

import "context"

func (db *DB) AddWithdrawal(ctx context.Context, orderID string, userID string, sum float64) error {
	query := "INSERT INTO withdrawals(order_id, sum, user_id) VALUES ($1, $2, $3)"
	_, err := db.pool.Exec(ctx, query, orderID, sum, userID)
	return err
}

func (db *DB) GetUserWithdrawals(ctx context.Context, userID string) ([]WithdrawalInfo, error) {
	query := "SELECT order_id, sum, withdrawn_at FROM withdrawals WHERE user_id=$1"
	rows, err := db.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	var result []WithdrawalInfo
	for rows.Next() {
		var item WithdrawalInfo
		err = rows.Scan(&item.OrderNumber, &item.Sum, &item.ProcessedAt)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, err
}
