package domain

import (
	"time"
)

// User представляет пользователя в системе.
type User struct {
	ID           int       `json:"-"     db:"id"`
	Login        string    `json:"login" db:"login"`
	PasswordHash string    `json:"-"     db:"password_hash"`
	CreatedAt    time.Time `json:"-"     db:"created_at"`
	UpdatedAt    time.Time `json:"-"     db:"updated_at"`
}

// AuthToken представляет данные авторизационного токена.
type AuthToken struct {
	Token string `json:"token"`
}

// UserRepository определяет интерфейс для доступа к данным пользователей.
type UserRepository interface {
	Create(user *User) error
	FindByLogin(login string) (*User, error)
}

// UserService определяет интерфейс для бизнес-логики работы с пользователями.
type UserService interface {
	Register(login, password string) (*AuthToken, error)
	Authenticate(login, password string) (*AuthToken, error)
}

// RegisterRequest представляет данные запроса на регистрацию.
type RegisterRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}
