package repository

import (
	"log/slog"

	"github.com/jmoiron/sqlx"

	"gophermart/internal/domain"
)

const (
	kopeksPerRuble = 100 // Количество копеек в рубле
)

// BalanceRepo реализует интерфейс domain.BalanceRepository
type BalanceRepo struct {
	db     *sqlx.DB
	logger *slog.Logger
}

// NewBalanceRepo создает новый экземпляр BalanceRepo
func NewBalanceRepo(db *sqlx.DB, logger *slog.Logger) *BalanceRepo {
	return &BalanceRepo{
		db: db,
		logger: logger.With(
			"package", "repository",
			"component", "BalanceRepo",
		),
	}
}

// GetBalance возвращает текущий баланс пользователя
func (r *BalanceRepo) GetBalance(userID int) (*domain.Balance, error) {
	var balance domain.Balance

	// Получаем сумму всех начислений (сразу в рублях, null значения заменяются на 0)
	err := r.db.Get(&balance.Current, `
		SELECT COALESCE(SUM(accrual), 0)::float / 100.0
		FROM orders 
		WHERE user_id = $1 AND status = 'PROCESSED'`, userID)
	if err != nil {
		return nil, err
	}

	// Получаем сумму всех списаний
	err = r.db.Get(&balance.Withdrawn, `
		SELECT COALESCE(SUM(amount_kop), 0)::float / 100.0
		FROM withdrawals 
		WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}

	// Вычитаем списания из начислений
	balance.Current -= balance.Withdrawn
	return &balance, nil
}

// CreateWithdrawal создает новую запись о списании средств
func (r *BalanceRepo) CreateWithdrawal(userID int, withdrawal *domain.Withdrawal) error {
	query := `
		INSERT INTO withdrawals (user_id, order_number, amount_kop)
		VALUES ($1, $2, $3)
		RETURNING processed_at`

	return r.db.QueryRow(
		query,
		userID,
		withdrawal.Order,
		withdrawal.AmountKop,
	).Scan(&withdrawal.ProcessedAt)
}

// GetWithdrawals возвращает историю списаний пользователя
func (r *BalanceRepo) GetWithdrawals(userID int) ([]domain.Withdrawal, error) {
	var withdrawals []domain.Withdrawal
	query := `
		SELECT order_number, amount_kop, processed_at
		FROM withdrawals 
		WHERE user_id = $1 
		ORDER BY processed_at DESC`

	if err := r.db.Select(&withdrawals, query, userID); err != nil {
		return nil, err
	}

	// Конвертируем копейки в рубли
	for i := range withdrawals {
		withdrawals[i].Sum = float64(withdrawals[i].AmountKop) / kopeksPerRuble
	}

	return withdrawals, nil
}
