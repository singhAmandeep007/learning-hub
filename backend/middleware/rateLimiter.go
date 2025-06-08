package middleware

import (
	"learning-hub/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiterMiddleware interface {
	RateLimiter() gin.HandlerFunc
}

type rateLimiterMiddleware struct {
	limit   int
	window  time.Duration
	clients map[string][]time.Time
	mu      sync.Mutex
}

func NewRateLimiterMiddleware(limit int, window time.Duration) RateLimiterMiddleware {
	return &rateLimiterMiddleware{
		limit:   limit,
		window:  window,
		clients: make(map[string][]time.Time),
	}
}

func (m *rateLimiterMiddleware) RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		m.mu.Lock()
		defer m.mu.Unlock()

		now := time.Now()

		// Create new entry for this client if it doesn't exist
		if _, exists := m.clients[clientIP]; !exists {
			m.clients[clientIP] = []time.Time{}
		}

		// Remove timestamps outside the current window
		var validRequests []time.Time
		for _, timestamp := range m.clients[clientIP] {
			if now.Sub(timestamp) <= m.window {
				validRequests = append(validRequests, timestamp)
			}
		}

		m.clients[clientIP] = validRequests

		// Check if the client has exceeded the limit
		if len(m.clients[clientIP]) >= m.limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, models.ErrorResponse{Error: "too_many_request", Message: "Rate limit exceeded"})
			return
		}

		// Add current request timestamp
		m.clients[clientIP] = append(m.clients[clientIP], now)

		c.Next()
	}
}
