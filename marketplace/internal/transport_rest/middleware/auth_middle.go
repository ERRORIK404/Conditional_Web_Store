package middleware

import (
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/localJWT"
	"github.com/ERRORIK404/Conditional_Web_Store/marketplace/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Уникальный ключ для контекста (строка)
const authStatusKey = "authStatus"

// AuthStatus хранит информацию об аутентификации
type AuthStatus struct {
	IsAuthenticated bool
	UserID          uuid.UUID
	// Дополнительные поля при необходимости
}

// AuthMiddleware создает middleware для проверки аутентификации
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		status := AuthStatus{IsAuthenticated: false} // Статус по умолчанию

		// Проверяем авторизацию
		token := c.GetHeader("Authorization")
		if token != "" {
			// Реальная валидация токена
			if userIDString, isValid := localjwt.CheckJWT(token); isValid {
				userID, err := uuid.Parse(userIDString)
				if err == nil {
					status = AuthStatus{
						IsAuthenticated: true,
						UserID:          userID,
					}
					logger.Logger.Info("")
				} else {
					logger.Logger.Debug("ERROR in auth(can't parse uuid")
				}
			}
		}

		// Сохраняем в контексте Gin с использованием строкового ключа
		c.Set(authStatusKey, status)
		c.Next()
	}
}

// GetAuthStatus извлекает статус аутентификации из контекста
func GetAuthStatus(c *gin.Context) AuthStatus {
	// Извлекаем значение по ключу
	val, exists := c.Get(authStatusKey)
	if !exists {
		return AuthStatus{IsAuthenticated: false}
	}

	// Приводим к правильному типу
	status, ok := val.(AuthStatus)
	if !ok {
		return AuthStatus{IsAuthenticated: false}
	}

	return status
}
