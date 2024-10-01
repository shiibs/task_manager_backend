package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/task_manager/database"
	"github.com/shiibs/task_manager/models"
)

func CreateTask(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	teamID := c.Param("team_id")

	var team models.Team
	if err := database.DBConn.First(&team, teamID).Error; err != nil || team.AdminID != adminID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied"})
		return
	}

	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		AssignedTo  uint      `json:"assigned_to" binding:"required"`
		Priority    string    `json:"priority" binding:"required"`
		Deadline    time.Time `json:"deadline" binding:"required"`
	}

	var member models.User
	if err := database.DBConn.Model(&team).Where("id = ?", req, req.AssignedTo).Association("Members").Find(&member); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member of this team"})
		return
	}

	newTask := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      models.NotStarted,
		Priority:    models.TaskPriority(req.Priority),
		AssignedTo:  req.AssignedTo,
		TeamID:      team.ID,
		AssignDate:  time.Now(),
		Deadline:    req.Deadline,
		CreatedBy:   adminID,
	}

	if err := database.DBConn.Create(&newTask).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": newTask})
}
