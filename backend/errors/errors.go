package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

// Error codes
const (
	// Validation errors (4xx)
	ErrInvalidParam       ErrorCode = "INVALID_PARAM"
	ErrInvalidContentType ErrorCode = "INVALID_CONTENT_TYPE"
	ErrInvalidPayload     ErrorCode = "INVALID_PAYLOAD"
	ErrFileTooLarge       ErrorCode = "FILE_TOO_LARGE"
	ErrUnsupportedType    ErrorCode = "UNSUPPORTED_TYPE"
	ErrMissingRequired    ErrorCode = "MISSING_REQUIRED"
	ErrInvalidFileType    ErrorCode = "INVALID_FILE_TYPE"

	// Authentication errors (4xx)
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrInvalidProduct    ErrorCode = "INVALID_PRODUCT"
	ErrForbidden         ErrorCode = "FORBIDDEN"
	ErrRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"

	// Resource errors (4xx)
	ErrResourceNotFound ErrorCode = "RESOURCE_NOT_FOUND"
	ErrResourceExists   ErrorCode = "RESOURCE_EXISTS"
	ErrTagNotFound      ErrorCode = "TAG_NOT_FOUND"

	// Database errors (5xx)
	ErrQueryFailed          ErrorCode = "QUERY_FAILED"
	ErrMutationFailed       ErrorCode = "MUTATION_FAILED"
	ErrDataConversionFailed ErrorCode = "DATA_CONVERSION_FAILED"

	// Storage errors (5xx)
	ErrUploadFailed ErrorCode = "UPLOAD_FAILED"
	ErrDeleteFailed ErrorCode = "DELETE_FAILED"

	// System errors (5xx)
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
)

// ErrorResponse represents the API error response structure
type ErrorResponse struct {
	Error   ErrorCode `json:"error"`             // Error code for frontend translation
	Message string    `json:"message,omitempty"` // Optional fallback message
	Details string    `json:"details,omitempty"` // Optional additional details
}

// errorMetadata maps error codes to HTTP status codes
var errorMetadata = map[ErrorCode]int{
	// Validation errors (4xx)
	ErrInvalidParam:       http.StatusBadRequest,
	ErrInvalidContentType: http.StatusBadRequest,
	ErrInvalidPayload:     http.StatusBadRequest,
	ErrFileTooLarge:       http.StatusBadRequest,
	ErrUnsupportedType:    http.StatusBadRequest,
	ErrMissingRequired:    http.StatusBadRequest,
	ErrInvalidFileType:    http.StatusBadRequest,

	// Authentication errors (4xx)
	ErrUnauthorized:      http.StatusUnauthorized,
	ErrInvalidProduct:    http.StatusBadRequest,
	ErrForbidden:         http.StatusForbidden,
	ErrRateLimitExceeded: http.StatusTooManyRequests,

	// Resource errors (4xx)
	ErrResourceNotFound: http.StatusNotFound,
	ErrResourceExists:   http.StatusConflict,
	ErrTagNotFound:      http.StatusNotFound,

	// Database errors (5xx)
	ErrQueryFailed:          http.StatusInternalServerError,
	ErrMutationFailed:       http.StatusInternalServerError,
	ErrDataConversionFailed: http.StatusInternalServerError,

	// Storage errors (5xx)
	ErrUploadFailed: http.StatusInternalServerError,
	ErrDeleteFailed: http.StatusInternalServerError,

	// System errors (5xx)
	ErrInternalServer: http.StatusInternalServerError,
}

// GetHTTPStatus returns the HTTP status code for an error code
func GetHTTPStatus(code ErrorCode) int {
	if status, exists := errorMetadata[code]; exists {
		return status
	}
	return http.StatusInternalServerError
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, code ErrorCode, message string) {
	status := GetHTTPStatus(code)
	response := ErrorResponse{
		Error:   code,
		Message: message,
	}
	c.JSON(status, response)
}

// RespondWithErrorDetails sends an error response with additional details
func RespondWithErrorDetails(c *gin.Context, code ErrorCode, message, details string) {
	status := GetHTTPStatus(code)
	response := ErrorResponse{
		Error:   code,
		Message: message,
		Details: details,
	}
	c.JSON(status, response)
}

// AbortWithError aborts the request with a standardized error response
func AbortWithError(c *gin.Context, code ErrorCode, message string) {
	status := GetHTTPStatus(code)
	response := ErrorResponse{
		Error:   code,
		Message: message,
	}
	c.AbortWithStatusJSON(status, response)
}
