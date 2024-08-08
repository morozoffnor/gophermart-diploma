package storage

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"log"
)

func (db *DB) GetUserCreds(userID string, ctx context.Context) (*auth.User, error) {
	user := &auth.User{
		Id: userID,
	}

	query := "SELECT login, password FROM users where id=$1"
	err := db.pool.QueryRow(ctx, query, userID).Scan(&user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (db *DB) CreateUser(ctx context.Context, login string, passwordHash string) (*auth.User, error) {
	newUser := &auth.User{
		Id:       uuid.New().String(),
		Login:    login,
		Password: passwordHash,
	}

	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	query := "INSERT INTO users (id, login, password) VALUES ($1, $2, $3)"
	_, err = tx.Exec(ctx, query, newUser.Id, newUser.Login, newUser.Password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Print("user with this id already exists, trying again")

			newUser.Id = uuid.New().String()
			_, err = tx.Exec(ctx, query, newUser.Id, newUser.Login, newUser.Password)
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

	return newUser, nil
}
