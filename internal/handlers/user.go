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
		ID: uuid.New().String(),
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
	token, err := h.auth.Jwt.GenerateToken(user.ID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	// добавляем токен в куки
	_, err = h.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		log.Print(err)
		return
	}

	// хешируем пароль, создаём запись о юзере в бд
	hashedPass := h.auth.HashPassword(user.Password)
	_, err = h.db.CreateUser(r.Context(), user.ID, user.Login, hashedPass)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	// получаем тело, распаковываем его, если надо
	body, err := bodyHelper.GetBody(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed parsing body", http.StatusBadRequest)
		return
	}
	user := &auth.User{}

	// вытаскиваем джейсонку в структуру
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	// находим такого юзера в бд
	dbUser, err := h.db.GetUserByLogin(r.Context(), user.Login)
	if err != nil {
		log.Print(err)
		http.Error(w, "No such user", http.StatusUnauthorized)
		return
	}

	// хешируем присланный пароль, сравниваем с хешем из бд
	hashedPass := h.auth.HashPassword(user.Password)
	if dbUser.Password != hashedPass {
		http.Error(w, "Wrong login/pass", http.StatusUnauthorized)
		return
	}

	// добавляем токен в куки
	token, err := h.auth.Jwt.GenerateToken(dbUser.ID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	_, err = h.auth.Jwt.AddTokenToCookies(&w, r, token)
	if err != nil {
		log.Print(err)
		return
	}

	w.Header().Set("Content-Type", "text/plain, utf-8")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ContextUserID).(string)
	balance, err := h.db.GetBalance(r.Context(), userID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(balance)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
