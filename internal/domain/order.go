package domain

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// OrderStatus представляет статус обработки заказа
type OrderStatus string

const (
	// OrderStatusNew заказ загружен в систему, но не попал в обработку
	OrderStatusNew OrderStatus = "NEW"
	// OrderStatusProcessing вознаграждение за заказ рассчитывается
	OrderStatusProcessing OrderStatus = "PROCESSING"
	// OrderStatusInvalid система расчета вознаграждений отказала в расчете
	OrderStatusInvalid OrderStatus = "INVALID"
	// OrderStatusProcessed данные по заказу проверены и информация о расчете успешно получена
	OrderStatusProcessed OrderStatus = "PROCESSED"
)

// Value реализует интерфейс driver.Valuer для OrderStatus
func (s OrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

// Scan реализует интерфейс sql.Scanner для OrderStatus
func (s *OrderStatus) Scan(value interface{}) error {
	if value == nil {
		*s = ""
		return nil
	}
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to scan %T into OrderStatus", value)
	}
	*s = OrderStatus(strVal)
	return nil
}

// Order представляет заказ в системе
type Order struct {
	ID         int         `json:"-"          db:"id"`
	Number     string      `json:"number"     db:"number"`
	UserID     int         `json:"-"          db:"user_id"`
	Status     OrderStatus `json:"status"     db:"status"`
	Accrual    *int64      `json:"-"          db:"accrual,omitempty"` // сумма начисленных баллов в копейках
	AccrualRub *float64    `json:"accrual,omitempty" db:"-"`          // сумма начисленных баллов в рублях для JSON
	UploadedAt time.Time   `json:"uploaded_at" db:"uploaded_at"`
}

// SetAccrual устанавливает сумму начисления в копейках и автоматически обновляет сумму в рублях
func (o *Order) SetAccrual(kop int64) {
	o.Accrual = &kop
	if kop > 0 {
		rub := float64(kop) / 100.0
		o.AccrualRub = &rub
	}
}

// GetAccrualRub возвращает сумму начисленных баллов в рублях
func (o *Order) GetAccrualRub() float64 {
	if o.Accrual == nil {
		return 0
	}
	return float64(*o.Accrual) / 100.0
}

// OrderRepository определяет интерфейс для доступа к данным заказов
type OrderRepository interface {
	// Create создает новый заказ
	Create(order *Order) error
	// FindByNumber ищет заказ по номеру
	FindByNumber(number string) (*Order, error)
	// FindByUserID возвращает все заказы пользователя
	FindByUserID(userID int) ([]Order, error)
	// FindByStatus возвращает заказы с указанными статусами
	FindByStatus(statuses []OrderStatus) ([]Order, error)
	// UpdateStatus обновляет статус заказа
	UpdateStatus(orderID int, status OrderStatus) error
	// UpdateAccrual обновляет сумму начисленных баллов за заказ
	UpdateAccrual(orderID int, accrualKop int64) error
}

// OrderService определяет интерфейс для бизнес-логики работы с заказами
type OrderService interface {
	// Register регистрирует новый заказ для пользователя
	Register(userID int, number string) error
	// GetOrders возвращает список заказов пользователя
	GetOrders(userID int) ([]Order, error)
}

// OrderRequest представляет данные запроса на регистрацию заказа
type OrderRequest struct {
	Number string `json:"number" validate:"required"`
}
