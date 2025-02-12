package domain

import "errors"

var (
	// ErrInvalidOrderNumber ошибка неверный номер заказа.
	ErrInvalidOrderNumber = errors.New("неверный номер заказа")
)
