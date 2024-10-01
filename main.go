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

	r.POST("/login", handlers.LoginUser)

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())

	user := r.Group("/user")
	user.Use(middleware.AuthMiddleware())

	// admin
	admin.POST("/register", handlers.RegisterUser)
	admin.GET("/users", handlers.ListUsers)

	admin.POST("/teams", handlers.CreateTeam)
	admin.POST("/teams/:team_id/members", handlers.AddMemeberToTeam)
	admin.DELETE("/teams/:team_id/members", handlers.RemoveMemberFromTeam)
	admin.DELETE("/teams/:team_id", handlers.DeleteTeam)

	//member
	user.GET("/profile", handlers.GetUserProfile)
	user.GET("/teams/:team_id", handlers.ViewTeamDetails)
	user.GET("/teams", handlers.ListTeams)

	r.Run(port)

}
