package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"learninghub/config"
	"learninghub/constants"
	"learninghub/firebase"
	"learninghub/handlers"
	"learninghub/middleware"
	logger "learninghub/pkg/logger"
	"learninghub/utils"
)

func main() {
	// Create context that listens for the interrupt signal from the OS.
	signalCtx, signalCtxStop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl+\
		syscall.SIGTERM, // ask a program to terminate
	)
	defer signalCtxStop()

	// Initialize logger
	logger.InitGlobal(
		logger.WithServiceName("learninghub-server"),
		logger.WithDefaultDestinations(logger.FileLogger, logger.ConsoleLogger),
		logger.WithConsoleDestination(),
		logger.WithFileDestination(utils.ResolvePathFromProjectRoot("logs/learninghub-server.log"), 10, 5, 30, true),
		logger.WithMinLevel(logger.DebugLevel),
	)
	defer logger.CloseGlobal()

	// Populate AppConfig with env variables
	err := config.LoadConfig()
	if err != nil {
		logger.Infof("Error loading configuration: %v", err)
	}

	// Initialize Firebase
	err = firebase.InitializeFirebase()
	if err != nil {
		logger.Infof("Failed to initialize Firebase: %v", err)
	}

	// Setup Gin router
	r := setupRouter()
	port := config.AppConfig.PORT

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Initializing the server in a goroutine so that it won't block the graceful shutdown handling below
	go func() {
		logger.Infof("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// listen for the interrupt signal
	<-signalCtx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	signalCtxStop()
	logger.Infof("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	shutdownCtx, shutdownCtxStop := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCtxStop()

	// Shutdown server
	logger.Infof("Shutting down HTTP server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Infof("Closing Firebase connections...")
		firebase.CloseFirebase()
		logger.Fatalf("Server forced to shutdown: %v\n", err)
	}

	logger.Infof("Closing Firebase connections...")
	firebase.CloseFirebase()

	logger.Infof("Server exiting")
}

func setupRouter() *gin.Engine {
	envMode := config.AppConfig.ENV_MODE
	if envMode == constants.EnvModeProd {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup Gin router
	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	// if envMode == constants.EnvModeDev {
	// 	r.Use(middleware.DelayMiddleware(1000 * time.Millisecond))
	// }

	r.Use(middleware.NewRateLimiterMiddleware(20, time.Minute).RateLimiterForMethods("POST", "PUT", "PATCH", "DELETE"))

	// API routes
	// /api/v1/:product/resources
	api := r.Group("/api/v1")
	{
		// Product-specific routes
		productGroup := api.Group("/:product", middleware.ProductValidationMiddleware())
		{
			productGroup.GET("/resources", handlers.GetResources)
			productGroup.GET("/resources/:id", handlers.GetResource)
			productGroup.POST("/resources", handlers.CreateResource)
			productGroup.PATCH("/resources/:id", handlers.UpdateResource)
			productGroup.DELETE("/resources/:id", handlers.DeleteResource)

			productGroup.GET("/tags", handlers.GetTags)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	return r
}
