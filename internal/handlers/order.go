package handlers

import (
	"database/sql"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	"gophermart/internal/domain"
	"gophermart/internal/service"
)

// OrderHandler обрабатывает HTTP-запросы, связанные с заказами
type OrderHandler struct {
	orderService domain.OrderService
}

// NewOrderHandler создает новый экземпляр OrderHandler
func NewOrderHandler(orderService domain.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

// Register обрабатывает загрузку номера заказа
// @Summary Загрузка номера заказа
// @Description Загружает номер заказа для расчета начисления баллов лояльности
// @Tags orders
// @Accept text/plain
// @Produce json
// @Param number body string true "Номер заказа"
// @Success 202 "Новый номер заказа принят в обработку"
// @Success 200 "Номер заказа уже был загружен этим пользователем"
// @Failure 400 "Неверный формат запроса"
// @Failure 401 "Пользователь не аутентифицирован"
// @Failure 409 "Номер заказа уже был загружен другим пользователем"
// @Failure 422 "Неверный формат номера заказа"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /api/user/orders [post]
func (h *OrderHandler) Register(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user_id in context")
	}

	// Проверяем Content-Type
	if !strings.HasPrefix(c.Request().Header.Get("Content-Type"), "text/plain") {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса (Content-Type должен быть text/plain)")
	}

	// Читаем тело запроса
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса (не удалось прочитать тело запроса)")
	}
	defer c.Request().Body.Close()

	// Преобразуем байты в строку и убираем пробелы
	number := strings.TrimSpace(string(body))
	if number == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса (тело запроса не может быть пустым)")
	}

	err = h.orderService.Register(userID, number)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderExists):
			return c.NoContent(http.StatusOK)
		case errors.Is(err, service.ErrOrderRegisteredByOther):
			return echo.NewHTTPError(http.StatusConflict, "Номер заказа уже был загружен другим пользователем")
		case errors.Is(err, service.ErrInvalidOrderNumber):
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "Неверный формат номера заказа")
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
	}

	return c.NoContent(http.StatusAccepted)
}

// GetOrders возвращает список заказов пользователя
// @Summary Получение списка заказов
// @Description Возвращает список загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях
// @Tags orders
// @Produce json
// @Success 200 {array} domain.Order "Список заказов"
// @Success 204 "Нет данных для ответа"
// @Failure 401 "Пользователь не аутентифицирован"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /api/user/orders [get]
func (h *OrderHandler) GetOrders(c echo.Context) error {
	userIDRaw := c.Get("user_id")
	userID, ok := userIDRaw.(int)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "invalid user_id in context")
	}

	orders, err := h.orderService.GetOrders(userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusNoContent)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	if len(orders) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	return c.JSON(http.StatusOK, orders)
}
