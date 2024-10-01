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

func DeleteTeam(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	teamID := c.Param("team_id")

	var team models.Team
	if err := database.DBConn.First(&team, teamID).Error; err != nil || team.AdminID != adminID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := database.DBConn.Delete(&team).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

func ViewTeamDetails(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	teamID := c.Param("team_id")

	var team models.Team
	if err := database.DBConn.Preload("Memebers").Preload("Tasks").First(&team, teamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	isAdmin := team.AdminID == userID
	isMember := database.DBConn.Model(&team).Where("id = ?", userID).Association("Members").Find(&userID) == nil

	if !isAdmin && !isMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"team":    team,
		"members": team.Members,
		"tasks":   team.Tasks,
	})
}

func ListTeams(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var teams []models.Team
	if err := database.DBConn.Joins("JOIN team_members ON teams.id = team_members.team_id").
		Where("team_members.user_id = ? OR teams.admin_id = ?", userID, userID).
		Find(&teams).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrive teams"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"teams": teams})
}
