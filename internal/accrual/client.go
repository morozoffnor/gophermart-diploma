package accrual

import (
	"encoding/json"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"io"
	"log"
	"net/http"
	"strconv"
)

const (
	StatusRegistered = "REGISTERED"
	StatusInvalid    = "INVALID"
	StatusProcessing = "PROCESSING"
	StatusProcessed  = "PROCESSED"
)

type AccrualClient struct {
	cfg *config.Config
}

type OrderStatus struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func NewClient(cfg *config.Config) *AccrualClient {
	c := &AccrualClient{cfg: cfg}
	return c
}

func (c *AccrualClient) GetOrderStatus(orderID string) (*OrderStatus, error) {
	response, err := http.Get(c.cfg.AccrualSystemAddr + "/api/orders/" + orderID)
	if err != nil {
		log.Print(err)
	}
	switch response.StatusCode {
	case http.StatusTooManyRequests:
		log.Print("too many requests")
		timeout, _ := strconv.Atoi(response.Header.Get("Retry-After"))
		return nil, &ErrorTooManyRequests{
			Timeout: timeout,
		}
	case http.StatusNoContent:
		log.Print("order not registered in accrual system")
		return nil, &ErrorNotRegistered{}
	case http.StatusInternalServerError:
		log.Print("Internal error")
		return nil, &ErrorInternalError{}
	default:
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		var orderStatus OrderStatus

		err = json.Unmarshal(body, &orderStatus)
		if err != nil {
			return nil, err
		}
		return &orderStatus, nil
	}

}
