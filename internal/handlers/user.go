package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"errors"
	"gophermart/internal/domain"
	"gophermart/internal/service"
)

// UserHandler обрабатывает HTTP-запросы, связанные с пользователями.
type UserHandler struct {
	userService domain.UserService
}

// NewUserHandler создает новый экземпляр UserHandler.
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register обрабатывает регистрацию пользователя.
// @Summary Регистрация нового пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RegisterRequest true "Учетные данные для регистрации"
// @Success 200 {object} domain.AuthToken "Пользователь успешно зарегистрирован"
// @Failure 400 "Неверный формат запроса"
// @Failure 409 "Логин уже занят"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /api/user/register [post]
// @Description Регистрирует нового пользователя с логином и паролем.
func (h *UserHandler) Register(c echo.Context) error {
	var req domain.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Ошибка валидации")
	}

	token, err := h.userService.Register(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) {
			return echo.NewHTTPError(http.StatusConflict, "Пользователь уже существует")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	// Устанавливаем токен в заголовок Authorization
	c.Response().Header().Set("Authorization", "Bearer "+token.Token)
	return c.JSON(http.StatusOK, token)
}

// LoginRequest представляет данные запроса на вход.
type LoginRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Authenticate обрабатывает аутентификацию пользователя.
// @Summary Аутентификация пользователя.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Учетные данные для входа"
// @Success 200 {object} domain.AuthToken "Пользователь успешно аутентифицирован"
// @Failure 400 "Неверный формат запроса"
// @Failure 401 "Неверная пара логин/пароль"
// @Failure 500 "Внутренняя ошибка сервера"
// @Router /api/user/login [post]
// @Description Аутентифицирует пользователя по логину и паролю.
func (h *UserHandler) Authenticate(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Ошибка валидации")
	}

	token, err := h.userService.Authenticate(req.Login, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidLogin) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Неверный логин или пароль")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	// Устанавливаем токен в заголовок Authorization
	c.Response().Header().Set("Authorization", "Bearer "+token.Token)
	return c.JSON(http.StatusOK, token)
}
