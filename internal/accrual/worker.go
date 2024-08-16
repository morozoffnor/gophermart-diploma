package accrual

import (
	"context"
	"errors"
	"github.com/morozoffnor/gophermart-diploma/internal/config"
	"github.com/morozoffnor/gophermart-diploma/internal/storage"
	"log"
	"time"
)

type Worker struct {
	cfg     *config.Config
	queueCh chan string
	db      *storage.DB
	client  *AccrualClient
}

func NewWorker(cfg *config.Config, db *storage.DB, c *AccrualClient) *Worker {
	w := &Worker{
		cfg:     cfg,
		queueCh: make(chan string),
		db:      db,
		client:  c,
	}
	return w
}

func (w *Worker) Start(ctx context.Context) {
	log.Print("Started worker")
	for order := range w.queueCh {
		go w.ProcessOrder(order)
	}
}

func (w *Worker) AddToQueue(order string) {
	log.Print("Adding to processing queue")
	w.queueCh <- order
}

func (w *Worker) ProcessStaleOrders(ctx context.Context) {
	orders, err := w.db.GetUnprocessedOrders(ctx)
	if err != nil {
		return
	}
	for _, order := range orders {
		w.AddToQueue(order)
	}
}

func (w *Worker) ProcessOrder(order string) {
	orderStatus, err := w.client.GetOrderStatus(order)
	if err != nil {
		var terr *ErrorTooManyRequests
		ok := errors.As(err, &terr)
		if ok {
			// ждём столько секунд, сколько нам сказано в ответе и пробуем снова
			time.Sleep(time.Duration(terr.Timeout*1000) * time.Second)
			w.AddToQueue(order)
			return
		}

		var nerr *ErrorNotRegistered
		ok = errors.As(err, &nerr)
		if ok {
			log.Printf("order %s is not registered in accrual system", order)
			return
		}

		var ierr *ErrorInternalError
		ok = errors.As(err, &ierr)
		if ok {
			log.Printf("an internal server error in accrual system occured with order %s", order)
			return
		}
		// если что-то другое, то выводим в консоль
		log.Print(err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// достаём заказ из бд
	dbOrder, err := w.db.GetOrder(ctx, order)
	if err != nil {
		// если ошибка с бд, то пробуем ещё раз через 5 секунд
		log.Print(err)
		time.Sleep(5 * time.Second)
		w.AddToQueue(order)
		return
	}

	// проверяем статус заказа в бд
	if dbOrder.Status == storage.StatusProcessed || dbOrder.Status == storage.StatusInvalid {
		// если заказ уже обработан, то просто выходим
		log.Printf("order %s already processed", order)
		return
	}

	// если система начислений вернула invalid, помечаем у себя тоже
	if orderStatus.Status == StatusInvalid {
		log.Printf("order %s in invalid", order)
		err := w.db.UpdateOrderFromAccrual(ctx, orderStatus.Order, orderStatus.Status, orderStatus.Accrual)
		if err != nil {
			// если ошибка с бд, пробуем ещё раз
			log.Print(err)
			time.Sleep(5 * time.Second)
			w.AddToQueue(order)
			return
		}
		return
	}

	// если заказ есть в системе начислений, то меняем статус у себя на PROCESSING
	if orderStatus.Status == StatusRegistered || orderStatus.Status == StatusProcessing {
		err := w.db.UpdateOrderFromAccrual(ctx, orderStatus.Order, storage.StatusProcessing, orderStatus.Accrual)
		if err != nil {
			// если ошибка с бд, пробуем ещё раз
			log.Print(err)
			time.Sleep(5 * time.Second)
			w.AddToQueue(order)
			return
		}
		return
	}

	// меняем статус на PROCESSED, начисляем баллы
	err = w.db.UpdateOrderFromAccrual(ctx, orderStatus.Order, orderStatus.Status, orderStatus.Accrual)
	if err != nil {
		// если ошибка с бд, пробуем ещё раз
		log.Print(err)
		time.Sleep(5 * time.Second)
		w.AddToQueue(order)
		return
	}
	err = w.db.UpdateBalance(ctx, dbOrder.UserID, orderStatus.Accrual)
	if err != nil {
		// если ошибка с бд, пробуем ещё раз
		log.Print(err)
		time.Sleep(5 * time.Second)
		w.AddToQueue(order)
		return
	}
	return
}
