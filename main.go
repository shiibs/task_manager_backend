package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shiibs/task_manager/database"
	"github.com/shiibs/task_manager/handlers"
	"github.com/shiibs/task_manager/middleware"
)

func init() {
	database.ConnectDB()
}

func main() {
	port := os.Getenv("PORT")

	psqlDB, err := database.DBConn.DB()
	if err != nil {
		panic("Error in DB connection")
	}
	defer psqlDB.Close()

	r := gin.New()

	// User routes
	r.POST("/login", handlers.LoginUser)
	r.POST("/register", middleware.AuthMiddleware(), middleware.AdminAuthMiddleware(), handlers.RegisterUser)
	r.GET("/profile", middleware.AuthMiddleware(), handlers.GetUserProfile)
	r.GET("/users", middleware.AuthMiddleware(), middleware.AdminAuthMiddleware(), handlers.ListUsers)

	r.Run(port)

}
