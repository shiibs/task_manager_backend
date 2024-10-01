package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/task_manager/database"
	"github.com/shiibs/task_manager/models"
)

func AddComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	taskID := c.Param("task_id")

	var task models.Task
	if err := database.DBConn.First(&task, taskID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var team models.Team
	if err := database.DBConn.First(&team, task.TeamID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}

	isAdmin := userID == team.AdminID
	isAssigned := task.AssignedTo == userID

	if !isAdmin && !isAssigned {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	newComment := models.Comment{
		TaskID:  task.ID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := database.DBConn.Create(&newComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully", "comment": newComment})
}

func ViewComments(c *gin.Context) {
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
	isAssigned := task.AssignedTo == userID

	if !isAdmin && !isAssigned {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": task.Comments})
}

func EditComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	commentID := c.Param("comment_id")

	var comment models.Comment
	if err := database.DBConn.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this comment"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	comment.Content = req.Content
	if err := database.DBConn.Save(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to edit comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment edited successfully", "comment": comment})
}

func DeleteComment(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	commentID := c.Param("comment_id")

	var comment models.Comment
	if err := database.DBConn.First(&comment, commentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this comment"})
		return
	}

	if err := database.DBConn.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
