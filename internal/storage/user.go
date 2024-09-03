package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"log"
)

func (db *DB) UserExists(ctx context.Context, login string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)"
	var exists bool
	err := db.pool.QueryRow(ctx, query, login).Scan(&exists)
	if err != nil {
		// решил тут возвращать именно true, чтобы случайно не создать
		// ещё одного пользователя. Человеческий фактор, все дела ;)
		return true, err
	}
	return exists, err
}

func (db *DB) GetUserByID(ctx context.Context, userID string) (*auth.User, error) {
	user := &auth.User{
		ID: userID,
	}

	query := "SELECT login, password FROM users where id=$1"
	err := db.pool.QueryRow(ctx, query, userID).Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) GetUserByLogin(ctx context.Context, login string) (*auth.User, error) {
	user := &auth.User{}

	query := "SELECT id, login, password FROM users where login=$1"
	err := db.pool.QueryRow(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) CreateUser(ctx context.Context, userID string, login string, passwordHash string) (*auth.User, error) {
	newUser := &auth.User{
		ID:       userID,
		Login:    login,
		Password: passwordHash,
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	query := "INSERT INTO users (id, login, password) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, query, newUser.ID, newUser.Login, newUser.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Print("user with this id already exists, trying again")

			// TODO: хз надо ли это вообще
			newUser.ID = uuid.New().String()
			_, err = tx.Exec(ctx, query, newUser.ID, newUser.Login, newUser.Password)
			if err != nil {
				log.Print("there is no way there are two identical uuids generated in a row, returning err: ", err)
				return nil, err
			}

			err = tx.Commit(ctx)
			if err != nil {
				log.Print("cannot commit: ", err)
			}

		}
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	newUser.Authenticated = true
	return newUser, nil
}

func (db *DB) GetBalance(ctx context.Context, userID string) (BalanceInfo, error) {
	query := "SELECT balance, withdrawn FROM users WHERE id=$1"
	var balance BalanceInfo
	err := db.pool.QueryRow(ctx, query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		log.Printf("error while getting balance, userid %s", userID)
		return balance, err
	}
	return balance, nil
}

func (db *DB) UpdateBalance(ctx context.Context, userID string, value float64) error {
	query := "UPDATE users SET balance = balance + $2 where id =$1"
	_, err := db.pool.Exec(ctx, query, userID, value)
	if err != nil {
		log.Printf("error while updating balance, userid %s", userID)
		return err
	}
	return err
}

func (db *DB) UpdateWithdrawals(ctx context.Context, userID string, value float64) error {
	query := "UPDATE users SET withdrawn = withdrawn + $2 where id =$1"
	_, err := db.pool.Exec(ctx, query, userID, value)
	if err != nil {
		log.Printf("error while updating withdrawn, userid %s", userID)
		return err
	}
	return err
}
