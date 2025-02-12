package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

const (
	statusBadRequest = 400
)

// CustomValidator пользовательский валидатор для фреймворка Echo.
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator создает новый экземпляр валидатора.
func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

// Validate проверяет переданную структуру.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(statusBadRequest, err.Error())
	}
	return nil
}
