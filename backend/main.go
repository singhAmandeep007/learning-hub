package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"learning-hub/config"
	"learning-hub/constants"
	"learning-hub/firebase"
	"learning-hub/handlers"
	"learning-hub/middleware"
)

func main() {
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
		return
	}
	defer func() {
		if err := firebase.CloseFirebase(); err != nil {
			log.Printf("Error during Firebase cleanup: %v", err)
		}
	}()

	// Setup Gin router
	r := setupRouter()

	port := config.AppConfig.PORT

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
		return
	}
}

func setupRouter() *gin.Engine {
	envMode := getEnvMode()
	if envMode == constants.EnvModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup Gin router
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// if envMode == constants.EnvModeDev {
	// 	r.Use(middleware.DelayMiddleware(1000 * time.Millisecond))
	// }

	r.Use(middleware.NewRateLimiterMiddleware(100, time.Minute).RateLimiter())

	// API routes
	api := r.Group("/api")
	{
		api.GET("/resources", handlers.GetResources)
		api.GET("/resources/:id", handlers.GetResource)
		api.POST("/resources", middleware.AdminAuthMiddleware(), handlers.CreateResource)
		api.PATCH("/resources/:id", middleware.AdminAuthMiddleware(), handlers.UpdateResource)
		api.DELETE("/resources/:id", middleware.AdminAuthMiddleware(), handlers.DeleteResource)

		api.GET("/tags", handlers.GetTags)
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return r
}

func getEnvMode() string {
	envMode := os.Getenv("ENV_MODE")
	// Requires ENV_MODE to be set in docker-compose.yml or in system: "dev" or "prod"
	if envMode != constants.EnvModeDev && envMode != constants.EnvModeProd {
		log.Fatalf("ENV_MODE environment variable is not set. Please set it to 'dev' or 'prod'.")
		return ""
	}
	return envMode
}
