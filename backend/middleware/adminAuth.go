package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"learning-hub/config"
	"learning-hub/constants"
	"learning-hub/models"
)

// AdminAuthMiddleware middleware for admin-only routes
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Expected format: "Bearer SECRET_KEY"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" && parts[1] == config.AdminSecret {
				c.Next()
				return
			}
		}

		// Check query parameter as fallback
		secret := c.Query(constants.AdminSecretQueryParamKey)
		if secret == config.AdminSecret {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "Admin authentication required",
		})
		c.Abort()
	}
}