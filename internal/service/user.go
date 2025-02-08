package service

import (
	"fmt"
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
	existingUser, findErr := s.repo.FindByLogin(login)
	if findErr == nil && existingUser != nil {
		return nil, ErrUserExists
	}

	// Хешируем пароль
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return nil, fmt.Errorf("failed to hash password: %w", hashErr)
	}

	// Создаем нового пользователя
	user := &domain.User{
		Login:        login,
		PasswordHash: string(hashedPassword),
	}

	// Сохраняем пользователя в базу
	if createErr := s.repo.Create(user); createErr != nil {
		return nil, fmt.Errorf("failed to create user: %w", createErr)
	}

	// Генерируем JWT токен
	token, tokenErr := s.generateToken(user.ID, user.Login)
	if tokenErr != nil {
		return nil, fmt.Errorf("failed to generate token: %w", tokenErr)
	}

	return token, nil
}

// Authenticate проверяет учетные данные пользователя и возвращает токен, если данные верны.
func (s *UserService) Authenticate(login, password string) (*domain.AuthToken, error) {
	// Ищем пользователя по логину
	user, findErr := s.repo.FindByLogin(login)
	if findErr != nil {
		return nil, ErrInvalidLogin
	}

	// Проверяем пароль
	if compareErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); compareErr != nil {
		return nil, ErrInvalidLogin
	}

	// Генерируем JWT токен
	token, tokenErr := s.generateToken(user.ID, user.Login)
	if tokenErr != nil {
		return nil, fmt.Errorf("failed to generate token: %w", tokenErr)
	}

	return token, nil
}
