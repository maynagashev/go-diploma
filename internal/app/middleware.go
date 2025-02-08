package app

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// extractTokenFromHeader извлекает JWT токен из заголовка Authorization.
func extractTokenFromHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Отсутствует токен авторизации")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Неверный формат токена")
	}

	return parts[1], nil
}

// validateToken проверяет JWT токен и возвращает claims.
func validateToken(tokenString string, secret string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неверный метод подписи токена")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неверный токен")
	}

	if !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Токен недействителен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "Неверный формат данных токена")
	}

	return claims, nil
}

// extractUserData извлекает данные пользователя из claims.
func extractUserData(claims jwt.MapClaims) (int, string, error) {
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", echo.NewHTTPError(http.StatusUnauthorized, "Отсутствует или неверный формат user_id")
	}

	login, ok := claims["login"].(string)
	if !ok {
		return 0, "", echo.NewHTTPError(http.StatusUnauthorized, "Отсутствует или неверный формат login")
	}

	return int(userIDFloat), login, nil
}

// JWTMiddleware создает middleware для проверки JWT токена.
func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из заголовка
			tokenString, err := extractTokenFromHeader(c)
			if err != nil {
				return err
			}

			// Проверяем токен и получаем claims
			claims, err := validateToken(tokenString, secret)
			if err != nil {
				return err
			}

			// Извлекаем данные пользователя
			userID, login, err := extractUserData(claims)
			if err != nil {
				return err
			}

			// Устанавливаем данные в контекст
			c.Set("user_id", userID)
			c.Set("login", login)

			return next(c)
		}
	}
}
