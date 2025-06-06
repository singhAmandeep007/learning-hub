package models

import "time"

// Resource represents a learning resource
type Resource struct {
	ID           string    `json:"id" firestore:"-"`
	Title        string    `json:"title" firestore:"title" binding:"required"`
	Description  string    `json:"description" firestore:"description" binding:"required"`
	Type         string    `json:"type" firestore:"type" binding:"required,oneof=video pdf article"`
	URL          string    `json:"url" firestore:"url"`
	ThumbnailURL string    `json:"thumbnailUrl,omitempty" firestore:"thumbnailUrl,omitempty"`
	Tags         []string  `json:"tags" firestore:"tags"`
	CreatedAt    time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" firestore:"updatedAt"`
}
