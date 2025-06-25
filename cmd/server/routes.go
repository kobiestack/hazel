package main

import (
	"github.com/freekobie/hazel/middlewares"
	"github.com/gin-gonic/gin"
)

func (app *application) routes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// users
	router.POST("/auth/register", app.handler.CreateUser)
	router.POST("/auth/login", app.handler.LoginUser)
	router.POST("/auth/access", app.handler.GetUserAccessToken)
	router.POST("/auth/verify", app.handler.VerifyUser)
	router.POST("/auth/verify/request", app.handler.RequestVerification)

	protected := router.Group("/")
	protected.Use(middlewares.Authentication())
	{
		//users
		protected.GET("/users/:id", app.handler.GetUser)
		protected.PATCH("/users/profile", app.handler.UpdateUserData)
		protected.DELETE("/users/:id", app.handler.DeleteUser)

		// workspaces
		protected.POST("/workspaces", app.handler.CreateWorkspace)
		protected.GET("/workspaces/:id", app.handler.GetWorkspace)
		protected.GET("/workspaces/me", app.handler.GetUserWorkspaces)
		protected.PATCH("/workspaces/:id", app.handler.UpdateWorkspace)
		protected.DELETE("/workspaces/:id", app.handler.DeleteWorkspace)
		protected.POST("/workspaces/:id/members", app.handler.AddWorkspaceMember)
		protected.GET("/workspaces/:id/members", app.handler.GetWorkspaceMembers)
		protected.DELETE("/workspaces/:id/members/:member_id", app.handler.DeleteWorkspaceMember)
		protected.GET("/workspaces/:id/projects", app.handler.GetProjectsInWorkspace)

		// projects
		protected.POST("/projects", app.handler.CreateProject)
		protected.GET("/projects/:id", app.handler.GetProject)
		protected.PATCH("/projects/:id", app.handler.UpdateProject)
		protected.DELETE("/projects/:id", app.handler.DeleteProject)
		protected.GET("/projects/:id/tasks", app.handler.GetProjectTasks)

		// Tasks
		protected.POST("/tasks", app.handler.CreateTask)
		protected.GET("/tasks/:id", app.handler.GetTask)
		protected.PATCH("/tasks/:id", app.handler.UpdateTask)
		protected.DELETE("/tasks/:id", app.handler.DeleteTask)
		protected.POST("/tasks/:id/assignments", app.handler.AssignTaskToUser)
		protected.GET("/tasks/:id/assignments", app.handler.GetAssignedUsers)
		protected.DELETE("/tasks/:id/assignments/:user_id", app.handler.RemoveAssignment)
	}

	return router
}
