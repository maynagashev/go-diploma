package service

import (
	"database/sql"
	"errors"
	"log/slog"

	"gophermart/internal/domain"
)

// OrderService реализует интерфейс domain.OrderService
type OrderService struct {
	repo domain.OrderRepository
}

// NewOrderService создает новый экземпляр OrderService
func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// Register регистрирует новый заказ для пользователя
func (s *OrderService) Register(userID int, number string) error {
	// Проверяем, существует ли заказ
	existingOrder, err := s.repo.FindByNumber(number)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	// Если заказ существует
	if existingOrder != nil {
		// Если заказ принадлежит текущему пользователю
		if existingOrder.UserID == userID {
			return ErrOrderExists
		}
		// Если заказ принадлежит другому пользователю
		return ErrOrderRegisteredByOther
	}

	// Проверяем номер заказа по алгоритму Луна
	if !isValidLuhn(number) {
		return ErrInvalidOrderNumber
	}

	// Создаем новый заказ
	order := &domain.Order{
		Number: number,
		UserID: userID,
		Status: domain.OrderStatusNew,
	}

	err = s.repo.Create(order)

	slog.Info("created new order",
		"id", order.ID,
		"uploaded_at", order.UploadedAt,
		"user_id", userID,
		"order_number", number,
		"status", domain.OrderStatusNew)

	return err
}

// GetOrders возвращает список заказов пользователя
func (s *OrderService) GetOrders(userID int) ([]domain.Order, error) {
	return s.repo.FindByUserID(userID)
}

// isValidLuhn проверяет номер заказа по алгоритму Луна
func isValidLuhn(number string) bool {
	// Преобразуем строку в слайс цифр
	digits := make([]int, len(number))
	for i, r := range number {
		if r < '0' || r > '9' {
			return false
		}
		digits[i] = int(r - '0')
	}

	// Алгоритм Луна
	sum := 0
	parity := len(digits) % 2
	for i, digit := range digits {
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}
