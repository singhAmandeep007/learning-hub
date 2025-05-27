package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	loadEnv()
	// Initialize Firebase
	app, err := initializeFirebase()
	if err != nil {
		log.Fatal("Failed to initialize Firebase:", err)
	}
	defer app.cleanup()

	// Setup Gin router
	r := gin.Default()
	
	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Configure appropriately for production
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API routes
	api := r.Group("/api")
	{
		api.GET("/resources", app.getResources)
		api.GET("/resources/:id", app.getResource)
		api.POST("/resources", app.adminAuth(), app.createResource)
		api.PUT("/resources/:id", app.adminAuth(), app.updateResource)
		api.DELETE("/resources/:id", app.adminAuth(), app.deleteResource)
		api.GET("/tags", app.getTags)
	}

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