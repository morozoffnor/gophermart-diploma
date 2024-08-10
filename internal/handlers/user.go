package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	bodyHelper "github.com/morozoffnor/go-url-shortener/pkg/body"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"log"
	"net/http"
)

func (h *Handlers) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// получаем тело, распаковываем его, если надо
	body, err := bodyHelper.GetBody(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed parsing body", http.StatusBadRequest)
		return
	}
	user := &auth.User{
		Id: uuid.New().String(),
	}

	// вытаскиваем джейсонку в структуру
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	// проверяем есть ли такой логин в базе
	exists, err := h.db.UserExists(r.Context(), user.Login)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// если есть, отдаём ошибку
	if exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// генерим новый JWT с айди юзера
	token, err := h.auth.Jwt.GenerateToken(user.Id)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// добавляем токен в куки
	ctx, err := h.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		log.Print(err)
		return
	}
	r = r.WithContext(ctx)

	// создаём запись о юзере в бд
	hashedPass := h.auth.HashPassword(user.Password)
	_, err = h.db.CreateUser(ctx, user.Id, user.Login, hashedPass)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}
