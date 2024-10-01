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
	admin.POST("/teams/:team_id/tasks", handlers.CreateTask)
	admin.PUT("/tasks/:task_id", handlers.UpdateTask)
	admin.DELETE("/tasks/:task_id", handlers.DeleteTask)
	admin.PATCH("/tasks/:task_id/reassign", handlers.ReassignTask)

	//member
	user.GET("/profile", handlers.GetUserProfile)
	user.GET("/teams/:team_id", handlers.ViewTeamDetails)
	user.GET("/teams", handlers.ListTeams)
	user.GET("/task/:task_id", handlers.ViewTaskDetails)
	user.PATCH("/tasks/:task_id/status", handlers.UpdateTaskStatus)
	user.GET("/teams/:team_id/tasks", handlers.ListTasks)
	user.POST("/tasks/:task_id/comments", handlers.AddComment)
	user.GET("/tasks/:task_id/comments", handlers.ViewComments)
	user.PUT("/comments/:comment_id", handlers.EditComment)
	user.DELETE("/comments/:comment_id", handlers.DeleteComment)

	r.Run(port)

}
