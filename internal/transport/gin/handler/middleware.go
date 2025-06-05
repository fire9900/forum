package handler

import (
	"errors"
	"github.com/fire9900/auth/pkg/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var (
	ErrMissingAuthHeader = errors.New("отсутствует заголовок авторизации")
	ErrInvalidToken      = errors.New("невалидный токен")
)

func AuthMiddleware(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": ErrMissingAuthHeader,
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		valid, err := authClient.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   ErrInvalidToken,
				"details": err.Error(),
			})
			return
		}

		if !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		userID, err := authClient.GetUserID(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Failed to get user ID",
				"details": err.Error(),
			})
			return
		}

		c.Set("userID", int(userID))
		c.Next()
	}
}
