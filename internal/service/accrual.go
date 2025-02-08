package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gophermart/internal/domain"
)

// AccrualResponse представляет ответ от системы начислений.
type AccrualResponse struct {
	Order   string             `json:"order"`
	Status  domain.OrderStatus `json:"status"`
	Accrual *float64           `json:"accrual,omitempty"`
}

// AccrualService сервис для взаимодействия с системой начислений.
type AccrualService struct {
	client  *http.Client
	baseURL string
}

// NewAccrualService создает новый экземпляр AccrualService.
func NewAccrualService(baseURL string) *AccrualService {
	return &AccrualService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

// GetOrderAccrual получает информацию о начислении баллов за заказ.
func (s *AccrualService) GetOrderAccrual(ctx context.Context, orderNumber string) (*AccrualResponse, error) {
	url := fmt.Sprintf("%s/api/orders/%s", s.baseURL, orderNumber)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var accrual AccrualResponse
		if err := json.NewDecoder(resp.Body).Decode(&accrual); err != nil {
			return nil, fmt.Errorf("ошибка декодирования ответа: %w", err)
		}
		return &accrual, nil
	case http.StatusTooManyRequests:
		// Получаем время ожидания из заголовка
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			seconds, err := time.ParseDuration(retryAfter + "s")
			if err == nil {
				return nil, &RateLimitError{RetryAfter: seconds}
			}
		}
		return nil, &RateLimitError{RetryAfter: 60 * time.Second} // По умолчанию ждем 60 секунд
	case http.StatusNoContent:
		return nil, nil
	default:
		return nil, fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}
}

// RateLimitError ошибка превышения лимита запросов.
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("превышен лимит запросов, повторить через %v", e.RetryAfter)
}
