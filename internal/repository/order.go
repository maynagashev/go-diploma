package repository

import (
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"gophermart/internal/domain"
)

// OrderRepo реализует интерфейс domain.OrderRepository.
type OrderRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewOrderRepo создает новый экземпляр OrderRepo.
func NewOrderRepo(db *sqlx.DB, logger *slog.Logger) *OrderRepo {
	return &OrderRepo{
		db: db,
		logger: logger.With(
			"package", "repository",
			"component", "OrderRepo",
		),
	}
}

// Create создает новый заказ.
func (r *OrderRepo) Create(order *domain.Order) error {
	query := `
		INSERT INTO orders (number, user_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, uploaded_at`

	return r.db.QueryRow(
		query,
		order.Number,
		order.UserID,
		order.Status,
	).Scan(&order.ID, &order.UploadedAt)
}

// FindByNumber ищет заказ по номеру.
func (r *OrderRepo) FindByNumber(number string) (*domain.Order, error) {
	var order domain.Order
	query := `SELECT * FROM orders WHERE number = $1`
	err := r.db.Get(&order, query, number)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// FindByUserID возвращает все заказы пользователя.
func (r *OrderRepo) FindByUserID(userID int) ([]domain.Order, error) {
	var orders []domain.Order
	query := `
		SELECT * FROM orders 
		WHERE user_id = $1 
		ORDER BY uploaded_at DESC`
	err := r.db.Select(&orders, query, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateStatus обновляет статус заказа.
func (r *OrderRepo) UpdateStatus(orderID int, status domain.OrderStatus) error {
	query := `
		UPDATE orders 
		SET status = $1 
		WHERE id = $2`
	_, err := r.db.Exec(query, status, orderID)
	return err
}

// UpdateAccrual обновляет сумму начисленных баллов за заказ.
func (r *OrderRepo) UpdateAccrual(orderID int, accrualKop int64) error {
	logger := r.logger.With("method", "UpdateAccrual")
	logger.Info("обновление статуса на PROCESSED",
		"id заказа", orderID,
		"начисление (коп)", accrualKop)

	query := `
		UPDATE orders 
		SET accrual = $1, status = $2 
		WHERE id = $3`
	_, err := r.db.Exec(query, accrualKop, domain.OrderStatusProcessed, orderID)
	return err
}

// FindByStatus возвращает заказы с указанными статусами.
func (r *OrderRepo) FindByStatus(statuses []domain.OrderStatus) ([]domain.Order, error) {
	logger := r.logger.With("method", "FindByStatus")

	// Преобразуем OrderStatus в []string для запроса к БД и логирования
	statusStrings := make([]string, len(statuses))
	for i, s := range statuses {
		statusStrings[i] = string(s)
	}

	query := `
		SELECT * FROM orders 
		WHERE status = ANY($1)
		ORDER BY uploaded_at ASC`

	var orders []domain.Order
	err := r.db.Select(&orders, query, statusStrings)
	if err != nil {
		logger.Error("ошибка при поиске заказов",
			"error", err,
			"тип ошибки", fmt.Sprintf("%T", err),
		)
		return nil, err
	}

	logger.Debug("поиск заказов по статусам", "статусы", statusStrings, "количество", len(orders))
	return orders, nil
}
