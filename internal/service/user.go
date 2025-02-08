package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"gophermart/internal/domain"
)

// UserService реализует интерфейс domain.UserService.
type UserService struct {
	repo           domain.UserRepository
	jwtSecret      []byte
	jwtExpiryHours int
}

// NewUserService создает новый экземпляр UserService.
func NewUserService(repo domain.UserRepository, jwtSecret string, jwtExpirationTime time.Duration) *UserService {
	return &UserService{
		repo:           repo,
		jwtSecret:      []byte(jwtSecret),
		jwtExpiryHours: int(jwtExpirationTime.Hours()),
	}
}

// generateToken создает новый JWT токен для пользователя.
func (s *UserService) generateToken(userID int, login string) (*domain.AuthToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"login":   login,
		"exp":     time.Now().Add(time.Duration(s.jwtExpiryHours) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.AuthToken{Token: tokenString}, nil
}

// Register создает нового пользователя с указанными учетными данными.
func (s *UserService) Register(login, password string) (*domain.AuthToken, error) {
	// Проверяем, существует ли пользователь
	existingUser, err := s.repo.FindByLogin(login)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &domain.User{
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Генерируем JWT токен
	return s.generateToken(user.ID, user.Login)
}

// Authenticate проверяет учетные данные пользователя и возвращает токен, если данные верны.
func (s *UserService) Authenticate(login, password string) (*domain.AuthToken, error) {
	user, err := s.repo.FindByLogin(login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrInvalidLogin
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidLogin
	}

	// Генерируем JWT токен
	return s.generateToken(user.ID, user.Login)
}
