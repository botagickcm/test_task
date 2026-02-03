package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"test_task/internal/models"
	"test_task/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := slog.Default()
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Отсутствует Authorization header")
			c.JSON(http.StatusUnauthorized, models.ServerResponse{
				Status: "error",
				Error:  "Authorization header required",
			})
			c.Abort()
			return
		}

		var tokenString string

		authHeader = strings.TrimSpace(authHeader)

		if strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			tokenString = strings.TrimSpace(authHeader[7:])
		} else {

			if strings.HasPrefix(strings.ToLower(authHeader), "bearer") && len(authHeader) > 6 {
				tokenString = strings.TrimSpace(authHeader[6:])
			} else {

				tokenString = authHeader
			}
		}

		tokenString = strings.TrimSpace(tokenString)

		if tokenString == "" {
			logger.Warn("Пустой токен", "original_header", authHeader)
			c.JSON(http.StatusUnauthorized, models.ServerResponse{
				Status: "error",
				Error:  "Token is empty",
			})
			c.Abort()
			return
		}

		// Логирование для отладки
		logger.Debug("Проверка токена",
			"original_header", authHeader,
			"extracted_token_length", len(tokenString),
			"token_prefix", "***"+tokenString[len(tokenString)-10:], // последние 10 символов
		)

		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			logger.Warn("Невалидный токен", "error", err.Error(), "token_length", len(tokenString))
			c.JSON(http.StatusUnauthorized, models.ServerResponse{
				Status: "error",
				Error:  "Invalid or expired token",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		logger.Debug("Токен валиден", "user_id", claims.UserID)
		c.Next()
	}
}
