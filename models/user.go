package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserID       string `gorm:"unique;not null" json:"user_id"`      // Unique identifier for the user
	Name         string `gorm:"not null" json:"name"`                // Name of the user
	Email        string `gorm:"unique;not null" json:"email"`        // User's email, must be unique
	PasswordHash string `gorm:"not null" json:"-"`                   // Hashed password (don't expose it in JSON responses)
	Role         string `gorm:"not null" json:"role"`                // Role: Admin or Member
	Teams        []Team `gorm:"many2many:team_members" json:"teams"` // Many-to-many relationship with teams
	Tasks        []Task `gorm:"foreignKey:AssignedTo" json:"tasks"`  // One-to-many relationship with tasks (assigned to the user)
}
