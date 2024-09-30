package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/task_manager/database"
	"github.com/shiibs/task_manager/models"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(uint)

		var user models.User
		if err := database.DBConn.First(&user, userID).Error; err != nil || user.Role != "Admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied. Admins only"})
			c.Abort()
			return
		}

		c.Next()
	}
}
