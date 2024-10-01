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

func AddMemeberToTeam(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	teamID := c.Param("team_id")

	var team models.Team
	if err := database.DBConn.First(&team, teamID).Error; err != nil || team.AdminID != adminID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := database.DBConn.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.DBConn.Model(&team).Association("Members").Append(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user to team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added to team successfully"})
}

func RemoveMemberFromTeam(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	teamID := c.Param("team_id")

	var team models.Team
	if err := database.DBConn.First(&team, teamID).Error; err != nil || team.AdminID != adminID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var user models.User
	if err := database.DBConn.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := database.DBConn.Model(&team).Association("Memebers").Delete(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove user from team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User removed from team successfully"})
}
