package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"learninghub/errors"
)

type RateLimiterMiddleware interface {
	RateLimiter() gin.HandlerFunc
	RateLimiterForMethods(methods ...string) gin.HandlerFunc
}

type rateLimiterMiddleware struct {
	limit   int
	window  time.Duration
	clients map[string][]time.Time
	mu      sync.RWMutex
}

func NewRateLimiterMiddleware(limit int, window time.Duration) RateLimiterMiddleware {
	return &rateLimiterMiddleware{
		limit:   limit,
		window:  window,
		clients: make(map[string][]time.Time),
	}
}

// RateLimiter creates a rate limiter middleware for all HTTP methods
func (m *rateLimiterMiddleware) RateLimiter() gin.HandlerFunc {
	return m.createRateLimiter(nil)
}

// RateLimiterForMethods creates a rate limiter middleware for specific HTTP methods
func (m *rateLimiterMiddleware) RateLimiterForMethods(methods ...string) gin.HandlerFunc {
	return m.createRateLimiter(methods)
}

func (m *rateLimiterMiddleware) createRateLimiter(methods []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// If specific methods are provided but none match, skip rate limiting
		if len(methods) > 0 {
			shouldLimit := false
			for _, method := range methods {
				if c.Request.Method == method {
					shouldLimit = true
					break
				}
			}
			if !shouldLimit {
				c.Next()
				return
			}
		}

		// Create unique key based on IP and method (if methods are specified)
		var key string
		if len(methods) > 0 {
			key = fmt.Sprintf("%s:%s", c.ClientIP(), c.Request.Method)
		} else {
			key = c.ClientIP()
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		now := time.Now()

		// Create new entry for this key if it doesn't exist
		if _, exists := m.clients[key]; !exists {
			m.clients[key] = []time.Time{}
		}

		// Remove timestamps outside the current window
		var validRequests []time.Time
		for _, timestamp := range m.clients[key] {
			if now.Sub(timestamp) <= m.window {
				validRequests = append(validRequests, timestamp)
			}
		}

		m.clients[key] = validRequests

		// Check if the client has exceeded the limit
		if len(m.clients[key]) >= m.limit {
			errors.AbortWithError(c, errors.ErrRateLimitExceeded, "Rate limit exceeded")
			return
		}

		// Add current request timestamp
		m.clients[key] = append(m.clients[key], now)

		c.Next()
	}
}
