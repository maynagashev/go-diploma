package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gophermart/internal/domain"
	"gophermart/internal/service"
)

// contextKey используется для ключей контекста
type contextKey string

const (
	defaultWorkerCount  = 5
	defaultPollInterval = 1 * time.Second
	defaultRetryTimeout = 1 * time.Minute

	workerIDKey = contextKey("worker_id")
)

// AccrualWorker обработчик заказов для получения информации о начислениях
type AccrualWorker struct {
	orderRepo      domain.OrderRepository
	accrualService *service.AccrualService
	workerCount    int
	pollInterval   time.Duration
	retryTimeout   time.Duration
}

// NewAccrualWorker создает новый экземпляр AccrualWorker
func NewAccrualWorker(
	orderRepo domain.OrderRepository,
	accrualService *service.AccrualService,
	workerCount int,
	pollInterval time.Duration,
	retryTimeout time.Duration,
) *AccrualWorker {
	if workerCount <= 0 {
		workerCount = defaultWorkerCount
	}
	if pollInterval <= 0 {
		pollInterval = defaultPollInterval
	}
	if retryTimeout <= 0 {
		retryTimeout = defaultRetryTimeout
	}

	return &AccrualWorker{
		orderRepo:      orderRepo,
		accrualService: accrualService,
		workerCount:    workerCount,
		pollInterval:   pollInterval,
		retryTimeout:   retryTimeout,
	}
}

// Start запускает обработку заказов
func (w *AccrualWorker) Start(ctx context.Context) {
	var wg sync.WaitGroup

	// Запускаем пул воркеров
	for i := 0; i < w.workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.worker(ctx, workerID)
		}(i)
	}

	// Ждем завершения всех воркеров
	wg.Wait()
}

// worker обрабатывает заказы
func (w *AccrualWorker) worker(ctx context.Context, id int) {
	logger := slog.With("worker_id", id)
	logger.Info("воркер начал работу")

	// Добавляем worker_id в контекст
	ctx = context.WithValue(ctx, workerIDKey, id)

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info("воркер завершил работу")
			return
		case <-ticker.C:
			if err := w.processOrders(ctx); err != nil {
				logger.Error("ошибка обработки заказов", "error", err)
				// Увеличиваем интервал опроса при ошибках
				ticker.Reset(w.retryTimeout)
			} else {
				// Возвращаем нормальный интервал опроса
				ticker.Reset(w.pollInterval)
			}
		}
	}
}

// processOrders обрабатывает заказы, ожидающие обновления статуса
func (w *AccrualWorker) processOrders(ctx context.Context) error {
	logger := slog.With(
		"func", "processOrders",
		"worker_id", ctx.Value(workerIDKey),
	)
	logger.Debug("starting to process orders")

	// Получаем заказы для обработки (NEW или PROCESSING)
	statuses := []domain.OrderStatus{domain.OrderStatusNew, domain.OrderStatusProcessing}
	logger.Debug("querying orders", "statuses", statuses)

	orders, err := w.orderRepo.FindByStatus(statuses)
	if err != nil {
		logger.Error("failed to find orders by status",
			"error", err,
			"error_type", fmt.Sprintf("%T", err))
		return err
	}

	logger.Debug("found orders to process", "count", len(orders))

	for _, order := range orders {
		// Проверяем контекст перед каждым заказом
		if ctx.Err() != nil {
			logger.Debug("context cancelled", "error", ctx.Err())
			return ctx.Err()
		}

		logger.Debug("processing order",
			"order_id", order.ID,
			"order_number", order.Number,
			"current_status", order.Status)

		// Получаем информацию о начислении
		accrual, err := w.accrualService.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			// Если превышен лимит запросов, ждем и пропускаем остальные заказы
			if rateLimitErr, ok := err.(*service.RateLimitError); ok {
				logger.Warn("rate limit exceeded",
					"retry_after", rateLimitErr.RetryAfter,
					"order_number", order.Number)
				time.Sleep(rateLimitErr.RetryAfter)
				return nil
			}
			logger.Error("failed to get accrual info",
				"order_number", order.Number,
				"error", err)
			continue
		}

		// Если заказ не найден, пропускаем
		if accrual == nil {
			logger.Debug("order not found in accrual system",
				"order_number", order.Number)
			continue
		}

		logger.Debug("received accrual info",
			"order_number", order.Number,
			"status", accrual.Status,
			"accrual", accrual.Accrual)

		// Обновляем статус заказа
		if err := w.orderRepo.UpdateStatus(order.ID, accrual.Status); err != nil {
			logger.Error("failed to update order status",
				"order_number", order.Number,
				"status", accrual.Status,
				"error", err)
			continue
		}

		// Если есть начисление, обновляем сумму
		if accrual.Status == domain.OrderStatusProcessed && accrual.Accrual != nil {
			accrualKop := int64(*accrual.Accrual * 100) // конвертируем рубли в копейки
			logger.Debug("updating order accrual",
				"order_number", order.Number,
				"accrual_rub", *accrual.Accrual,
				"accrual_kop", accrualKop)

			if err := w.orderRepo.UpdateAccrual(order.ID, accrualKop); err != nil {
				logger.Error("failed to update order accrual",
					"order_number", order.Number,
					"accrual_kop", accrualKop,
					"error", err)
			}
		}
	}

	return nil
}
