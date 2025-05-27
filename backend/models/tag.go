package models

// Tag represents a tag with usage statistics
type Tag struct {
	Name       string `json:"name" firestore:"name"`
	UsageCount int    `json:"usageCount" firestore:"usageCount"`
}
