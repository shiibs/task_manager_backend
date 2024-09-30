package database

import (
	"log"
	"os"

	"github.com/shiibs/task_manager/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func ConnectDB() {
	dns := os.Getenv("DATABASE_URL")

	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		panic("Database connection failed!")
	}

	log.Println("DB Connected")
	db.AutoMigrate(&models.User{}, &models.Team{}, &models.Task{}, &models.Comment{})
}
