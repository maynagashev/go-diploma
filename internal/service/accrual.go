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

const (
	defaultClientTimeout = 10 * time.Second
	defaultRetryTimeout  = 60 * time.Second
)

// AccrualService сервис для взаимодействия с системой начислений.
type AccrualService struct {
	client  *http.Client
	baseURL string
}

// NewAccrualService создает новый экземпляр AccrualService.
func NewAccrualService(baseURL string) *AccrualService {
	return &AccrualService{
		client: &http.Client{
			Timeout: defaultClientTimeout,
		},
		baseURL: baseURL,
	}
}

// GetOrderAccrual получает информацию о начислении баллов за заказ.
func (s *AccrualService) GetOrderAccrual(ctx context.Context, orderNumber string) (*domain.OrderAccrual, error) {
	url := fmt.Sprintf("%s/api/orders/%s", s.baseURL, orderNumber)
	req, reqErr := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if reqErr != nil {
		return nil, fmt.Errorf("failed to create request: %w", reqErr)
	}

	resp, respErr := s.client.Do(req)
	if respErr != nil {
		return nil, fmt.Errorf("failed to do request: %w", respErr)
	}
	defer resp.Body.Close()

	// Если заказ не найден, возвращаем nil без ошибки
	if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	// Проверяем rate limit
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			seconds, durationErr := time.ParseDuration(retryAfter + "s")
			if durationErr != nil {
				return nil, fmt.Errorf("failed to parse retry after: %w", durationErr)
			}
			return nil, &RateLimitError{RetryAfter: seconds}
		}
	}

	// Проверяем успешность ответа
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var accrual domain.OrderAccrual
	if decodeErr := json.NewDecoder(resp.Body).Decode(&accrual); decodeErr != nil {
		return nil, fmt.Errorf("failed to decode response: %w", decodeErr)
	}

	return &accrual, nil
}

// RateLimitError ошибка превышения лимита запросов.
type RateLimitError struct {
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	return fmt.Sprintf("превышен лимит запросов, повторить через %v", e.RetryAfter)
}
