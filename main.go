package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"be-awarenix/config"
	"be-awarenix/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	// Init DB
	config.InitDatabase()

	// Setup Gin engine
	app := gin.Default()
	// app.Use(middlewares.CORSMiddleware())
	// app.Use(cors.Default())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Load routes
	routes.SetupRoutes(app)

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	app.Run(fmt.Sprintf(":%s", port))
}
