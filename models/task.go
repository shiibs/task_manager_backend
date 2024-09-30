package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	NotStarted  TaskStatus = "Not Started"
	InProgress  TaskStatus = "In Progress"
	UnderReview TaskStatus = "Under Review"
	Completed   TaskStatus = "Completed"
)

type TaskPriority string

const (
	Low    TaskPriority = "Low"
	Medium TaskPriority = "Medium"
	High   TaskPriority = "High"
	Urgent TaskPriority = "Urgent"
)

type Task struct {
	gorm.Model
	Title       string       `gorm:"not null" json:"title"`
	Description string       `json:"description"`
	Status      TaskStatus   `gorm:"type:enum;not null;default:'Not Started'" json:"status"`
	Priority    TaskPriority `gorm:"type:enum;not null;default:'Medium'" json:"priority"`
	AssignedTo  uint         `json:"assigned_to"` // Foreign key to the User assigned to this task
	TeamID      uint         `json:"team_id"`     // Foreign key to the Team the task belongs to
	AssignDate  time.Time    `gorm:"not null" json:"assign_date"`
	Deadline    time.Time    `gorm:"not null" json:"deadline"`
	Comments    []Comment    `gorm:"foreignKey:TaskID" json:"comments"`
	CreatedBy   uint         `json:"created_by"`
}
