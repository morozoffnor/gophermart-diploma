package handlers

import (
	"encoding/json"
	bodyHelper "github.com/morozoffnor/go-url-shortener/pkg/body"
	"github.com/morozoffnor/gophermart-diploma/internal/auth"
	"github.com/morozoffnor/gophermart-diploma/pkg/luhn"
	"log"
	"net/http"
)

func (h *Handlers) Withdraw(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ContextUserID).(string)

	body, err := bodyHelper.GetBody(r)
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed parsing body", http.StatusBadRequest)
		return
	}

	var withdrawal WithdrawRequest
	err = json.Unmarshal(body, &withdrawal)
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid json", http.StatusUnprocessableEntity)
		return
	}

	// получаем текущий баланс из бд
	balance, err := h.db.GetBalance(r.Context(), userID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// проверяем хватает ли баллов на балансе
	if withdrawal.Sum > balance.Current {
		log.Print("insufficient funds")
		http.Error(w, "insufficient funds", http.StatusPaymentRequired)
		return
	}

	// проверяем заказ на валидность и существует ли он вообще
	if !luhn.Valid(withdrawal.OrderNumber) {
		log.Print("Invalid order number")
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}
	exists, err := h.db.OrderExists(r.Context(), withdrawal.OrderNumber)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !exists {
		log.Print("Invalid order number")
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	// списываем баллы
	err = h.db.UpdateBalance(r.Context(), userID, -withdrawal.Sum)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// сохраняем кол-во списанных баллов
	err = h.db.UpdateWithdrawals(r.Context(), userID, withdrawal.Sum)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// сохраняем списание
	err = h.db.AddWithdrawal(r.Context(), withdrawal.OrderNumber, userID, withdrawal.Sum)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

}

func (h *Handlers) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.ContextUserID).(string)

	// забираем списания из бд
	withdrawals, err := h.db.GetUserWithdrawals(r.Context(), userID)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// маршаллим в джейсонку
	resp, err := json.Marshal(withdrawals)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
