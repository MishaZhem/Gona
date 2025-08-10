package middleware

import (
	"net/http"

	"github.com/MishaZhem/Gona/src/gateway/internal/app"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no token provided"})
			return
		}

		userId, err := a.ValidateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("userID", userId)
		c.Next()
	}
}
