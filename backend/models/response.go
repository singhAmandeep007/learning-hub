package models

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data        []Resource `json:"data"`
	NextCursor  string     `json:"nextCursor,omitempty"`
	HasMore     bool       `json:"hasMore"`
	Total       int        `json:"total,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}