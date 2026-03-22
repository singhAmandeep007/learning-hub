package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"learninghub/config"
	"learninghub/constants"
)

const testValidProduct = "ecomm"

func TestGetProductFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		setupContext    func(*gin.Context)
		expectedProduct string
		expectedExists  bool
	}{
		{
			name: "Valid product in context",
			setupContext: func(c *gin.Context) {
				c.Set(constants.ProductContextKey, testValidProduct)
			},
			expectedProduct: testValidProduct,
			expectedExists:  true,
		},
		{
			name: "No product in context",
			setupContext: func(c *gin.Context) {
				// Don't set anything
			},
			expectedProduct: "",
			expectedExists:  false,
		},
		{
			name: "Wrong type in context",
			setupContext: func(c *gin.Context) {
				c.Set(constants.ProductContextKey, 123) // Set wrong type
			},
			expectedProduct: "",
			expectedExists:  false,
		},
		{
			name: "Nil value in context",
			setupContext: func(c *gin.Context) {
				c.Set(constants.ProductContextKey, nil)
			},
			expectedProduct: "",
			expectedExists:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Setup context
			tt.setupContext(c)

			// Test function
			product, exists := GetProductFromContext(c)

			assert.Equal(t, tt.expectedExists, exists)
			assert.Equal(t, tt.expectedProduct, product)
		})
	}
}

func TestProductValidationMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config.AppConfig = &config.EnvConfig{
		VALID_PRODUCTS: []string{testValidProduct},
	}

	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Track if next handler was called
	nextCalled := false

	// Setup route with middleware
	middleware := ProductValidationMiddleware()
	r.GET("/:product/resources", middleware, func(c *gin.Context) {
		nextCalled = true

		// Verify product is available in handler
		product, exists := GetProductFromContext(c)
		assert.True(t, exists)
		assert.Equal(t, testValidProduct, product)

		c.JSON(http.StatusOK, gin.H{"product": product})
	})

	// Create request with valid product
	req, _ := http.NewRequest("GET", "/"+testValidProduct+"/resources", nil)
	c.Request = req

	// Execute request
	r.ServeHTTP(w, req)

	// Verify middleware passed through to next handler
	assert.True(t, nextCalled)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify response contains product
	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, testValidProduct, response["product"])
}
