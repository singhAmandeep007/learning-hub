package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"learning-hub/constants"
	"learning-hub/models"
)

func TestHandleMultipartFormError(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedError  string
		expectedMsg    string
	}{
		{
			name:           "File too large error",
			err:            errors.New("multipart: file too large"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "form_parse_error",
			expectedMsg:    "File too large. Maximum size is 100 MB",
		},
		{
			name:           "No multipart boundary error",
			err:            errors.New("multipart: no multipart boundary"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "form_parse_error",
			expectedMsg:    "Invalid multipart form data - no boundary found",
		},
		{
			name:           "Generic error",
			err:            errors.New("some other error"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "form_parse_error",
			expectedMsg:    "Failed to parse form data: some other error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test request
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte{}))
			req.Header.Set("Content-Type", "multipart/form-data")
			
			// Create a new Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call the function
			handleMultipartFormError(c, tt.err)

			// Assert response status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response body
			var response models.ErrorResponse
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			// Assert error response fields
			assert.Equal(t, tt.expectedError, response.Error)
			assert.Equal(t, tt.expectedMsg, response.Message)
		})
	}
}

// TestHandleMultipartFormErrorWithMaxFileSize tests the file size message formatting
func TestHandleMultipartFormErrorWithMaxFileSize(t *testing.T) {
	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create a test request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer([]byte{}))
	req.Header.Set("Content-Type", "multipart/form-data")
	
	// Create a new Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call the function with a file too large error
	handleMultipartFormError(c, errors.New("multipart: file too large"))

	// Parse response body
	var response models.ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// Assert that the message contains the correct file size
	expectedSize := constants.MaxFileSize / (1 << 20) // Convert to MB
	assert.Contains(t, response.Message, fmt.Sprintf("%d MB", expectedSize))
} 