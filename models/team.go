package models

import (
	"gorm.io/gorm"
)

type Team struct {
	gorm.Model
	Name    string `gorm:"not null" json:"name"`                  // Team name
	AdminID uint   `gorm:"not null" json:"admin_id"`              // Foreign key to the User ID of the admin
	Members []User `gorm:"many2many:team_members" json:"members"` // Members of the team
	Tasks   []Task `gorm:"foreignKey:TeamID" json:"tasks"`        // Tasks associated with this team
}
