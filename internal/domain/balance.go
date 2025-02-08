package domain

import "time"

// Balance представляет баланс пользователя.
type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

// Withdrawal представляет списание средств.
type Withdrawal struct {
	Order       string    `json:"order" db:"order_number"`
	Sum         float64   `json:"sum" db:"-"`
	AmountKop   int64     `json:"-" db:"amount_kop"`
	ProcessedAt time.Time `json:"processed_at" db:"processed_at"`
}

// WithdrawalRequest представляет запрос на списание средств.
type WithdrawalRequest struct {
	Order string  `json:"order" validate:"required"`
	Sum   float64 `json:"sum" validate:"required,gt=0"`
}

// BalanceRepository определяет интерфейс для работы с балансом.
type BalanceRepository interface {
	GetBalance(userID int) (*Balance, error)
	CreateWithdrawal(userID int, withdrawal *Withdrawal) error
	GetWithdrawals(userID int) ([]Withdrawal, error)
}

// BalanceService определяет интерфейс для бизнес-логики работы с балансом.
type BalanceService interface {
	GetBalance(userID int) (*Balance, error)
	Withdraw(userID int, req *WithdrawalRequest) error
	GetWithdrawals(userID int) ([]Withdrawal, error)
}
