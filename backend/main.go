package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"learning-hub/config"
	"learning-hub/constants"
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

func loadEnv() {
	envMode := getEnvMode()

	var envFile string
	switch envMode {
	case constants.EnvModeDev:
		envFile = ".env.local.dev"
	case constants.EnvModeProd:
		envFile = ".env.local.prod"
	default:
		log.Printf("Environment mode for local development %s not recognized, using system environment variables", envMode)
		return
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Local environment file %s not found. Using system environment variables instead", envFile)
		return
	}

	log.Printf("Loaded environment from %s", envFile)
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
	// Requires ENV_MODE to be set in docker-compose.yml or in system: "dev" or "prod"
	envMode := os.Getenv("ENV_MODE")
	// Default to dev if ENV_MODE is not set
	if envMode == "" {
		envMode = constants.EnvModeDev
	}
	return envMode
}
