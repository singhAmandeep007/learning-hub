package models

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       []Resource `json:"data"`
	NextCursor string     `json:"nextCursor,omitempty"`
	HasMore    bool       `json:"hasMore"`
	Total      int        `json:"total,omitempty"`
}
