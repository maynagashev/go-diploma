package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gophermart/internal/domain"
	"gophermart/internal/service"
)

// BalanceHandler обработчик запросов для работы с балансом
type BalanceHandler struct {
	balanceService domain.BalanceService
}

// NewBalanceHandler создает новый экземпляр BalanceHandler
func NewBalanceHandler(balanceService domain.BalanceService) *BalanceHandler {
	return &BalanceHandler{
		balanceService: balanceService,
	}
}

// Register регистрирует обработчики в Echo
func (h *BalanceHandler) Register(e *echo.Echo) {
	e.GET("/api/user/balance", h.GetBalance)
	e.POST("/api/user/balance/withdraw", h.Withdraw)
	e.GET("/api/user/withdrawals", h.GetWithdrawals)
}

// GetBalance возвращает текущий баланс пользователя
func (h *BalanceHandler) GetBalance(c echo.Context) error {
	userID := c.Get("user_id").(int)

	balance, err := h.balanceService.GetBalance(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	return c.JSON(http.StatusOK, balance)
}

// Withdraw обрабатывает запрос на списание средств
func (h *BalanceHandler) Withdraw(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var req domain.WithdrawalRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Ошибка валидации")
	}

	err := h.balanceService.Withdraw(userID, &req)
	if err != nil {
		switch err {
		case domain.ErrInvalidOrderNumber:
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Неверный номер заказа")
		case service.ErrInsufficientFunds:
			return echo.NewHTTPError(http.StatusPaymentRequired, "Недостаточно средств")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	return c.NoContent(http.StatusOK)
}

// GetWithdrawals возвращает историю списаний пользователя
func (h *BalanceHandler) GetWithdrawals(c echo.Context) error {
	userID := c.Get("user_id").(int)

	withdrawals, err := h.balanceService.GetWithdrawals(userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	if len(withdrawals) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, withdrawals)
}
