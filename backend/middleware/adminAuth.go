package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learning-hub/config"
	"learning-hub/constants"
	"learning-hub/models"
)

// AdminAuthMiddleware middleware for admin-only routes
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check header
		authHeaderValue := c.GetHeader(constants.AdminSecretHeader)

		adminSecret := config.AppConfig.ADMIN_SECRET

		if authHeaderValue != "" {
			// Expected format: "SECRET_KEY"
			if authHeaderValue == adminSecret {
				c.Next()
				return
			}
		}

		// Check query parameter as fallback
		secret := c.Query(constants.AdminSecretQueryParamKey)
		if secret == adminSecret {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   constants.Unauthorized,
			Message: "Admin authentication required",
		})
	}
}
