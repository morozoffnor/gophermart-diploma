package handlers

import (
	"encoding/json"
	"github.com/morozoffnor/go-url-shortener/pkg/body"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"github.com/morozoffnor/gophermart-diploma/pkg/luhn"
	"log"
	"net/http"
)

func (h *Handlers) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ContextUserID).(string)

	orders, err := h.db.GetOrdersList(r.Context(), userID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(orders)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (h *Handlers) UploadOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ContextUserID).(string)

	reqBody, err := body.GetBody(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed parsing body", http.StatusBadRequest)
		return
	}
	// проверяем валиден ли номер заказа, возвращаем ошибку если нет 422
	orderNumber := string(reqBody)
	if !luhn.Valid(orderNumber) {
		log.Print("Invalid order number")
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	// проверяем существует ли уже такой заказ
	exists, err := h.db.OrderExists(r.Context(), orderNumber)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		order, err := h.db.GetOrder(r.Context(), orderNumber)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		// если айдишник пользователя не совпадает, возвращаем ошибку 409
		if order.UserID != userID {
			log.Print("This order belongs to another user")
			http.Error(w, "This order is already created by another user", http.StatusConflict)
			return
		}
		// если айдишник пользователя совпадает, то говорим, что всё ок 200
		log.Print("Order already created by this user")
		w.WriteHeader(http.StatusOK)
		return
	}
	// если не существует, то создаём новый заказ
	err = h.db.AddOrder(r.Context(), userID, orderNumber)
	h.worker.AddToQueue(orderNumber)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
