package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/task_manager/database"
	"github.com/shiibs/task_manager/models"
)

func CreateTeam(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)

	var req struct {
		Name string `json:"name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newTeam := models.Team{
		Name:    req.Name,
		AdminID: adminID,
	}

	if err := database.DBConn.Create(&newTeam).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Team created successfully", "team": newTeam})
}
