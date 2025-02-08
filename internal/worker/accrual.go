package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"errors"
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
	logger         *slog.Logger
	orderRepo      domain.OrderRepository
	accrualService *service.AccrualService
	workerCount    int
	pollInterval   time.Duration
	retryTimeout   time.Duration
}

// NewAccrualWorker создает новый экземпляр AccrualWorker
func NewAccrualWorker(
	logger *slog.Logger,
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
		logger: logger.With(
			"package", "worker",
			"component", "AccrualWorker",
		),
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
	for workerID := range w.workerCount {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.worker(ctx, workerID)
		}(workerID)
	}

	// Ждем завершения всех воркеров
	wg.Wait()
}

// worker обрабатывает заказы
func (w *AccrualWorker) worker(ctx context.Context, id int) {
	// Создаем отдельный логгер для этого воркера
	workerLogger := w.logger.With("worker_id", id)
	workerLogger.Info("воркер начал работу")

	// Добавляем worker_id в контекст
	ctx = context.WithValue(ctx, workerIDKey, id)

	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			workerLogger.Info("воркер завершил работу")
			return
		case <-ticker.C:
			if err := w.processOrders(ctx, workerLogger); err != nil {
				workerLogger.Error("ошибка обработки заказов", "error", err)
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
func (w *AccrualWorker) processOrders(ctx context.Context, logger *slog.Logger) error {
	logger = logger.With("method", "processOrders")
	logger.Debug("начало обработки заказов")

	// Получаем заказы для обработки (NEW или PROCESSING)
	statuses := []domain.OrderStatus{domain.OrderStatusNew, domain.OrderStatusProcessing}
	logger.Debug("запрос заказов", "статусы", statuses)

	orders, err := w.orderRepo.FindByStatus(statuses)
	if err != nil {
		logger.Error("ошибка при поиске заказов",
			"error", err,
			"error_type", fmt.Sprintf("%T", err))
		return err
	}

	logger.Debug("кол-во заказов для обработки воркером", "количество", len(orders))

	for _, order := range orders {
		// Проверяем контекст перед каждым заказом
		if ctx.Err() != nil {
			logger.Debug("контекст отменен", "error", ctx.Err())
			return ctx.Err()
		}

		logger.Debug("обработка заказа",
			"id заказа", order.ID,
			"номер заказа", order.Number,
			"текущий статус", order.Status)

		// Получаем информацию о начислении
		accrual, err := w.accrualService.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			var rateLimitErr *service.RateLimitError
			if errors.As(err, &rateLimitErr) {
				w.logger.Info("rate limit exceeded, waiting",
					"order_number", order.Number,
					"retry_after", rateLimitErr.RetryAfter)
				time.Sleep(rateLimitErr.RetryAfter)
				continue
			}
			w.logger.Error("failed to get order accrual",
				"order_number", order.Number,
				"error", err)
			continue
		}

		// Если заказ не найден, пропускаем
		if accrual == nil {
			logger.Debug("заказ не найден в системе начислений",
				"номер заказа", order.Number)
			continue
		}

		logger.Debug("получена информация о начислении",
			"номер заказа", order.Number,
			"статус", accrual.Status,
			"начисление", accrual.Accrual)

		// Обновляем статус заказа
		if err := w.orderRepo.UpdateStatus(order.ID, accrual.Status); err != nil {
			logger.Error("ошибка обновления статуса заказа",
				"номер заказа", order.Number,
				"статус", accrual.Status,
				"error", err)
			continue
		}

		// Если есть начисление, обновляем сумму
		if accrual.Status == domain.OrderStatusProcessed && accrual.Accrual != nil {
			accrualKop := int64(*accrual.Accrual * 100) // конвертируем рубли в копейки
			logger.Debug("обновление суммы начисления",
				"номер заказа", order.Number,
				"начисление (руб)", *accrual.Accrual,
				"начисление (коп)", accrualKop)

			if err := w.orderRepo.UpdateAccrual(order.ID, accrualKop); err != nil {
				logger.Error("ошибка обновления суммы начисления",
					"номер заказа", order.Number,
					"начисление (коп)", accrualKop,
					"error", err)
			}
		}
	}

	return nil
}
