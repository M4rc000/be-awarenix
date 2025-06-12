package routes

import (
	"be-awarenix/controllers"
	"be-awarenix/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Middleware global
	// router.Use(gin.Logger(), gin.Recovery())

	// Public routes
	router.POST("/api/v1/auth/login", controllers.AuthLogin)
	router.POST("/api/v1/auth/logout", middlewares.JWTAuth(), controllers.AuthLogout)

	// Protected API routes (dengan JWT middleware,)
	api := router.Group("/api/v1")
	api.Use(middlewares.JWTAuth())
	{
		users := api.Group("/users")
		{
			users.POST("/register", controllers.RegisterUser)
			users.GET("/session", controllers.GetUserSession)
			users.GET("/all", controllers.GetUsers)      // Get all users with pagination, search, sorting
			users.DELETE("/:id", controllers.DeleteUser) // Delete user
		}
	}
}
