package service

import (
	"errors"
	"log/slog"

	"gophermart/internal/domain"
	"gophermart/internal/utils"
)

var (
	// ErrInsufficientFunds ошибка недостаточно средств.
	ErrInsufficientFunds = errors.New("недостаточно средств")
)

// BalanceService реализует интерфейс domain.BalanceService.
type BalanceService struct {
	repo   domain.BalanceRepository
	logger *slog.Logger
}

// NewBalanceService создает новый экземпляр BalanceService.
func NewBalanceService(repo domain.BalanceRepository, logger *slog.Logger) *BalanceService {
	return &BalanceService{
		repo: repo,
		logger: logger.With(
			"package", "service",
			"component", "BalanceService",
		),
	}
}

// GetBalance возвращает текущий баланс пользователя.
func (s *BalanceService) GetBalance(userID int) (*domain.Balance, error) {
	return s.repo.GetBalance(userID)
}

// Withdraw списывает средства с баланса пользователя.
func (s *BalanceService) Withdraw(userID int, req *domain.WithdrawalRequest) error {
	// Проверяем номер заказа по алгоритму Луна
	if !utils.ValidateLuhn(req.Order) {
		return domain.ErrInvalidOrderNumber
	}

	// Получаем текущий баланс
	balance, err := s.repo.GetBalance(userID)
	if err != nil {
		return err
	}

	// Проверяем достаточно ли средств
	if balance.Current < req.Sum {
		return ErrInsufficientFunds
	}

	// Создаем запись о списании
	withdrawal := &domain.Withdrawal{
		Order:     req.Order,
		AmountKop: int64(req.Sum * domain.KopPerRuble),
	}

	return s.repo.CreateWithdrawal(userID, withdrawal)
}

// GetWithdrawals возвращает историю списаний пользователя.
func (s *BalanceService) GetWithdrawals(userID int) ([]domain.Withdrawal, error) {
	return s.repo.GetWithdrawals(userID)
}
