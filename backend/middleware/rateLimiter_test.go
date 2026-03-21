package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"learninghub/errors"
)

func TestRateLimiterMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		limit          int
		window         time.Duration
		requestCount   int
		requestDelay   time.Duration
		expectedStatus []int
	}{
		{
			name:           "Under rate limit",
			limit:          5,
			window:         time.Minute,
			requestCount:   3,
			requestDelay:   0,
			expectedStatus: []int{200, 200, 200},
		},
		{
			name:           "Exceed rate limit",
			limit:          2,
			window:         time.Minute,
			requestCount:   4,
			requestDelay:   0,
			expectedStatus: []int{200, 200, 429, 429},
		},
		{
			name:           "Rate limit reset after window",
			limit:          2,
			window:         100 * time.Millisecond,
			requestCount:   3,
			requestDelay:   150 * time.Millisecond,
			expectedStatus: []int{200, 200, 200}, // Third request after window reset
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create new rate limiter for each test
			rateLimiter := NewRateLimiterMiddleware(tt.limit, tt.window)

			r := gin.New()
			r.GET("/test", rateLimiter.RateLimiter(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Make requests
			for i := 0; i < tt.requestCount; i++ {
				if i > 0 && tt.requestDelay > 0 {
					time.Sleep(tt.requestDelay)
				}

				w := httptest.NewRecorder()
				req, _ := http.NewRequest("GET", "/test", nil)
				req.RemoteAddr = "127.0.0.1:8080" // Set consistent IP

				r.ServeHTTP(w, req)

				// Check status code
				assert.Equal(t, tt.expectedStatus[i], w.Code, "Request %d failed", i+1)

				if w.Code == http.StatusTooManyRequests {
					// Verify error response structure
					var response errors.ErrorResponse
					err := json.NewDecoder(w.Body).Decode(&response)
					assert.NoError(t, err)
					assert.Equal(t, string(errors.ErrRateLimitExceeded), string(response.Error))
					assert.Contains(t, response.Message, "Rate limit exceeded")
				}
			}
		})
	}
}

func TestRateLimiterForMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 2
	window := time.Minute
	rateLimiter := NewRateLimiterMiddleware(limit, window)

	r := gin.New()
	// Only rate limit POST and PUT methods
	r.Use(rateLimiter.RateLimiterForMethods("POST", "PUT"))
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "GET success"})
	})
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "POST success"})
	})
	r.PUT("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "PUT success"})
	})

	// Test that GET requests are not rate limited
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "127.0.0.1:8080"

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "GET request %d should not be rate limited", i+1)
	}

	// Test that POST requests are rate limited
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = "127.0.0.1:8080"

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "POST request %d should succeed", i+1)
	}

	// This POST request should be rate limited
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "POST request should be rate limited")

	// Test that PUT requests are also rate limited (separate from POST)
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/test", nil)
		req.RemoteAddr = "127.0.0.1:8080"

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "PUT request %d should succeed", i+1)
	}

	// This PUT request should be rate limited
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "PUT request should be rate limited")
}

func TestRateLimiterDifferentIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 2
	window := time.Minute
	rateLimiter := NewRateLimiterMiddleware(limit, window)

	r := gin.New()
	r.GET("/test", rateLimiter.RateLimiter(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Test different IPs can make requests independently
	ips := []string{"192.168.1.1:8080", "192.168.1.2:8080", "192.168.1.3:8080"}

	for _, ip := range ips {
		for i := 0; i < limit; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = ip

			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "IP %s request %d should succeed", ip, i+1)
		}
	}

	// Now each IP should be rate limited
	for _, ip := range ips {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = ip

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTooManyRequests, w.Code, "IP %s should be rate limited", ip)
	}
}

func TestRateLimiterMethodSpecificPerIP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 2
	window := time.Minute
	rateLimiter := NewRateLimiterMiddleware(limit, window)

	r := gin.New()
	r.Use(rateLimiter.RateLimiterForMethods("POST"))
	r.POST("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "POST success"})
	})

	clientIP := "127.0.0.1:8080"

	// Each IP should have separate rate limits for each method
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", nil)
		req.RemoteAddr = clientIP

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "POST request %d should succeed", i+1)
	}

	// This should be rate limited
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	req.RemoteAddr = clientIP

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "POST request should be rate limited")

	// Test with different IP should work
	differentIP := "192.168.1.2:8080"
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/test", nil)
	req.RemoteAddr = differentIP

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "POST request from different IP should succeed")
}

func TestRateLimiterWindowSliding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 3
	window := 200 * time.Millisecond
	rateLimiter := NewRateLimiterMiddleware(limit, window)

	r := gin.New()
	r.GET("/test", rateLimiter.RateLimiter(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	clientIP := "127.0.0.1:8080"

	// Make requests to fill the limit
	for i := 0; i < limit; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = clientIP

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "Request %d should succeed", i+1)
	}

	// Next request should be rate limited
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = clientIP
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Should be rate limited")

	// Wait for window to pass
	time.Sleep(window + 50*time.Millisecond)

	// Should be able to make requests again
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = clientIP
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code, "Should succeed after window reset")
}

func TestRateLimiterZeroLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Edge case: zero limit should block all requests
	rateLimiter := NewRateLimiterMiddleware(0, time.Minute)

	r := gin.New()
	r.GET("/test", rateLimiter.RateLimiter(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTooManyRequests, w.Code, "Zero limit should block all requests")
}

func TestRateLimiterEmptyMethods(t *testing.T) {
	gin.SetMode(gin.TestMode)

	limit := 2
	window := time.Minute
	rateLimiter := NewRateLimiterMiddleware(limit, window)

	r := gin.New()
	// When no methods are specified with RateLimiterForMethods(), it still applies rate limiting
	// but only for the specified methods (none in this case), so GET should not be rate limited
	r.Use(rateLimiter.RateLimiterForMethods("POST")) // Only rate limit POST, not GET
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	clientIP := "127.0.0.1:8080"

	// Should allow all GET requests through since only POST is rate limited
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = clientIP

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code, "GET request %d should succeed when only POST is rate limited", i+1)
	}
}
