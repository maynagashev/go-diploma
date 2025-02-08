package service

import (
	"database/sql"
	"errors"
	"log/slog"

	"gophermart/internal/domain"
	"gophermart/internal/utils"
)

// OrderService реализует интерфейс domain.OrderService.
type OrderService struct {
	repo domain.OrderRepository
}

// NewOrderService создает новый экземпляр OrderService.
func NewOrderService(repo domain.OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

// Register регистрирует новый заказ для пользователя.
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
	if !utils.ValidateLuhn(number) {
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

// GetOrders возвращает список заказов пользователя.
func (s *OrderService) GetOrders(userID int) ([]domain.Order, error) {
	orders, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Вычисляем сумму в рублях для каждого заказа
	for i := range orders {
		orders[i].CalculateAccrualRub()
	}

	return orders, nil
}
