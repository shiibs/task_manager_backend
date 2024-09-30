package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/shiibs/task_manager/utils"
)

// AuthMiddleware validates the JWT token and authenticates the user
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the token from the Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token not provided"})
			c.Abort()
			return
		}

		// Validate the token
		token, err := utils.ValidateJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Extract user ID from the token claims
		claims := token.Claims.(jwt.MapClaims)
		userID, ok := claims["user_id"].(float64) // JWT numeric values are floats
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Store the user ID in the context for use in the next handlers
		c.Set("user_id", uint(userID))

		// Proceed to the next handler
		c.Next()
	}
}
