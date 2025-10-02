package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WebSocketAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from query parameter for WebSocket upgrade
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token required"})
			c.Abort()
			return
		}

		// TODO: Validate JWT token and extract user_id
		// For now, placeholder validation
		// validated, userID := validateJWT(token)
		// if !validated {
		//     c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		//     c.Abort()
		//     return
		// }
		// c.Set("user_id", userID)

		c.Next()
	}
}