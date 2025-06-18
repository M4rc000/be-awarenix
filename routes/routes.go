package routes

import (
	"be-awarenix/controllers"
	"be-awarenix/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Middleware global
	router.Use(gin.Logger(), gin.Recovery())

	// Public routes
	router.POST("/api/v1/auth/login", controllers.AuthLogin)
	router.POST("/api/v1/auth/logout", middlewares.JWTAuth(), controllers.AuthLogout)

	// Protected API routes (dengan JWT middleware,)
	api := router.Group("/api/v1")
	api.Use(middlewares.JWTAuth())
	{
		groups := api.Group("/groups")
		{
			groups.GET("/all", controllers.GetGroups)
			groups.POST("/register", controllers.RegisterGroup)
			groups.PUT("/:id", controllers.UpdateUser) // Edit Grooup
		}

		users := api.Group("/users")
		{
			users.POST("/register", controllers.RegisterUser)
			users.GET("/session", controllers.GetUserSession)
			users.GET("/all", controllers.GetUsers)      // Get all users with pagination, search, sorting
			users.PUT("/:id", controllers.UpdateUser)    // Edit User
			users.DELETE("/:id", controllers.DeleteUser) // Delete user
		}

		emailTemplate := api.Group("/email-template")
		{
			emailTemplate.GET("/all", controllers.GetEmailTemplates)
			emailTemplate.POST("/create", controllers.RegisterEmailTemplate)
			emailTemplate.PUT("/:id", controllers.UpdateEmailTemplate) // Edit User

		}

		landingPage := api.Group("/landing-page")
		{
			landingPage.GET("/all", controllers.GetLandingPages)
		}

		sendingprofiles := api.Group("/sending-profile")
		{
			sendingprofiles.GET("/all", controllers.GetSendingProfiles)
		}

		analytics := api.Group("/analytics")
		{
			analytics.GET("/growth-percentage", controllers.GetGrowthPercentage)
		}
	}
}
