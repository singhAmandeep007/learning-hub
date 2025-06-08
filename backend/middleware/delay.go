package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

func DelayMiddleware(delay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		time.Sleep(delay) // Pause execution for the specified duration
		c.Next()          // Continue to the next handler in the chain
	}
}
