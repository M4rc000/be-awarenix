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
			users.GET("/session", controllers.GetUserSession)
			users.GET("/all", controllers.GetUsers) // Get all users with pagination, search, sorting
			// users.GET("/:id", controllers.GetUserByID)    // Get single user by ID
			// users.POST("/", controllers.CreateUser)       // Create new user
			// users.PUT("/:id", controllers.UpdateUser)     // Update user
			// users.DELETE("/:id", controllers.DeleteUser)  // Delete user
			// users.GET("/stats", controllers.GetUserStats) // Get user statistics
		}
	}
}
