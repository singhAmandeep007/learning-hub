package middleware

import (
	"net/url"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"learninghub/config"
	"learninghub/constants"
	"learninghub/pkg/logger"
)

// CORSMiddleware to enable CORS support
func CORSMiddleware() gin.HandlerFunc {
	allowOrigins := getValidCORSOrigins(config.AppConfig.CORS_ORIGINS)

	return cors.New(cors.Config{
		AllowOrigins:     allowOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

func getValidCORSOrigins(origins string) []string {
	originList := strings.Split(origins, ",")
	var validOrigins []string

	for _, origin := range originList {
		origin = strings.TrimSpace(origin)

		// Validate URL format
		if _, err := url.Parse(origin); err != nil {
			logger.Infof("Invalid CORS origin format: %s", origin)
			continue
		}

		// Ensure HTTPS in production
		if config.AppConfig.ENV_MODE == constants.EnvModeProd && !strings.HasPrefix(origin, "https://") {
			logger.Infof("Only HTTPS origins should be allowed in production: %s", origin)
			continue
		}

		validOrigins = append(validOrigins, origin)
	}

	return validOrigins
}
