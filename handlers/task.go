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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
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

func ViewTaskDetails(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskID := c.Param("task_id")

	var task models.Task
	if err := database.DBConn.Preload("Comments").First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var team models.Team
	if err := database.DBConn.First(&team, task.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	isAdmin := team.AdminID == userID
	isAssignedMember := task.AssignedTo == userID
	if !isAdmin && !isAssignedMember {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"task":     task,
		"comments": task.Comments,
	})
}

func UpdateTask(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	taskID := c.Param("task_id")

	var task models.Task
	if err := database.DBConn.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var team models.Team
	if err := database.DBConn.First(&team, task.TeamID).Error; err != nil || team.AdminID != adminID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		Title       string    `json:"title" binding:"required"`
		Description string    `json:"description"`
		Priority    string    `json:"priority" binding:"required"`
		Deadline    time.Time `json:"deadline" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	task.Title = req.Title
	task.Description = req.Description
	task.Priority = models.TaskPriority(req.Priority)
	task.Deadline = req.Deadline

	if err := database.DBConn.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully", "task": task})
}

func UpdateTaskStatus(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskID := c.Param("task_id")

	var task models.Task
	if err := database.DBConn.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	if task.AssignedTo != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not assigned to this task"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	task.Status = models.TaskStatus(req.Status)

	if err := database.DBConn.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task status updated successfully", "task": task})
}

func DeleteTask(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	taskID := c.Param("task_id")

	var task models.Task
	if err := database.DBConn.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
	}

	var team models.Team
	if err := database.DBConn.First(&team, task.TeamID).Error; err != nil || team.AdminID != adminID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := database.DBConn.Delete(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted succesfully"})
}

func ListTasks(c *gin.Context) {
	teamID := c.Param("team_id")

	var tasks []models.Task
	query := database.DBConn.Where("team_id = ?", teamID)

	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	if err := query.Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrivee tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

func ReassignTask(c *gin.Context) {
	adminID := c.MustGet("user_id").(uint)
	taskID := c.Param("task_id")

	var task models.Task
	if err := database.DBConn.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var team models.Team
	if err := database.DBConn.First(&team, task.TeamID).Error; err != nil || team.AdminID != adminID {
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

	var member models.User
	if err := database.DBConn.Model(&team).Where("id = ?", req.UserID).Association("Members").Find(&member); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "User is not a member of this team"})
		return
	}

	task.AssignedTo = req.UserID
	if err := database.DBConn.Save(&task).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reassign task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task reassigned sussesfully", "task": task})
}
