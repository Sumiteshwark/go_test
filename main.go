package main

import (
	"backend/config"
	"backend/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func main() {
	db := config.SetupDatabaseConnection()

	router := gin.Default()

	// Define CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"*"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}

	// Apply CORS middleware to the Gin router
	router.Use(cors.New(corsConfig))

	routes.UserRoutes(router, db)
	routes.AuthRoutes(router)
	routes.SystemRoutes(router)

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the backend",
		})
	})

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
