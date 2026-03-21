package errors

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetHTTPStatus(t *testing.T) {
	tests := []struct {
		name         string
		errorCode    ErrorCode
		expectedCode int
	}{
		{
			name:         "Invalid param error",
			errorCode:    ErrInvalidParam,
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "Unauthorized error",
			errorCode:    ErrUnauthorized,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Resource not found error",
			errorCode:    ErrResourceNotFound,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "Rate limit exceeded error",
			errorCode:    ErrRateLimitExceeded,
			expectedCode: http.StatusTooManyRequests,
		},
		{
			name:         "Query failed error",
			errorCode:    ErrQueryFailed,
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:         "Unknown error code defaults to 500",
			errorCode:    ErrorCode("UNKNOWN_ERROR"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := GetHTTPStatus(tt.errorCode)
			assert.Equal(t, tt.expectedCode, status)
		})
	}
}

func TestRespondWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		errorCode     ErrorCode
		message       string
		expectedCode  int
		expectedError ErrorCode
		expectedMsg   string
	}{
		{
			name:          "Bad request error",
			errorCode:     ErrInvalidParam,
			message:       "Invalid parameter provided",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidParam,
			expectedMsg:   "Invalid parameter provided",
		},
		{
			name:          "Not found error",
			errorCode:     ErrResourceNotFound,
			message:       "Resource not found",
			expectedCode:  http.StatusNotFound,
			expectedError: ErrResourceNotFound,
			expectedMsg:   "Resource not found",
		},
		{
			name:          "Internal server error",
			errorCode:     ErrQueryFailed,
			message:       "Database query failed",
			expectedCode:  http.StatusInternalServerError,
			expectedError: ErrQueryFailed,
			expectedMsg:   "Database query failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			RespondWithError(c, tt.errorCode, tt.message)

			assert.Equal(t, tt.expectedCode, w.Code)

			var response ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedError, response.Error)
			assert.Equal(t, tt.expectedMsg, response.Message)
			assert.Empty(t, response.Details)
		})
	}
}

func TestRespondWithErrorDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	errorCode := ErrUploadFailed
	message := "File upload failed"
	details := "File size exceeds maximum limit of 500MB"

	RespondWithErrorDetails(c, errorCode, message, details)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, errorCode, response.Error)
	assert.Equal(t, message, response.Message)
	assert.Equal(t, details, response.Details)
}

func TestAbortWithError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		errorCode     ErrorCode
		message       string
		expectedCode  int
		expectedError ErrorCode
		expectedMsg   string
	}{
		{
			name:          "Rate limit exceeded",
			errorCode:     ErrRateLimitExceeded,
			message:       "Too many requests",
			expectedCode:  http.StatusTooManyRequests,
			expectedError: ErrRateLimitExceeded,
			expectedMsg:   "Too many requests",
		},
		{
			name:          "Invalid product",
			errorCode:     ErrInvalidProduct,
			message:       "Invalid product parameter",
			expectedCode:  http.StatusBadRequest,
			expectedError: ErrInvalidProduct,
			expectedMsg:   "Invalid product parameter",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			AbortWithError(c, tt.errorCode, tt.message)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.True(t, c.IsAborted())

			var response ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedError, response.Error)
			assert.Equal(t, tt.expectedMsg, response.Message)
			assert.Empty(t, response.Details)
		})
	}
}

func TestErrorResponseJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Test that the JSON structure is correct
	RespondWithError(c, ErrInvalidParam, "Test message")

	contentType := w.Header().Get("Content-Type")
	assert.Contains(t, contentType, "application/json")

	var response map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify JSON structure
	assert.Contains(t, response, "error")
	assert.Contains(t, response, "message")
	assert.Equal(t, "INVALID_PARAM", response["error"])
	assert.Equal(t, "Test message", response["message"])
}
