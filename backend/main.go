package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"learning-hub/config"
	"learning-hub/firebase"
	"learning-hub/handlers"
	"learning-hub/middleware"
)

func main() {
	// Load environment variables
	loadEnv()
	// Populate AppConfig with env variables
	err := config.LoadConfig() 
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
		return
	}
	
	// Initialize Firebase
	err = firebase.InitializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer func() {
		if err := firebase.CloseFirebase(); err != nil {
				log.Printf("Error during Firebase cleanup: %v", err)
		}
	}()

	// Setup Gin router
	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		api.GET("/resources", handlers.GetResources)
		api.GET("/resources/:id", handlers.GetResource)
		api.POST("/resources", middleware.AdminAuthMiddleware(), handlers.CreateResource)
		api.PATCH("/resources/:id", middleware.AdminAuthMiddleware(), handlers.UpdateResource)
		api.DELETE("/resources/:id",middleware.AdminAuthMiddleware(), handlers.DeleteResource)

		api.GET("/tags", handlers.GetTags)
	}

	r.Use(middleware.CORSMiddleware())

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func loadEnv() {
	// Load environment variables from .env file
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}