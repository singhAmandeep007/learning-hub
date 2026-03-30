package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"learninghub/models"
)

func TestMapResourceToContract_WithThumbnail(t *testing.T) {
	createdAt := time.Date(2026, 3, 30, 10, 0, 0, 0, time.UTC)
	updatedAt := createdAt.Add(2 * time.Hour)

	resource := models.Resource{
		ID:           "res-1",
		Title:        "Test Resource",
		Description:  "Resource description",
		Type:         "article",
		URL:          "https://example.com/article",
		ThumbnailURL: "https://example.com/thumb.png",
		Tags:         []string{"api", "contract"},
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	mapped := mapResourceToContract(resource)

	assert.Equal(t, "res-1", mapped.Id)
	assert.Equal(t, "Test Resource", mapped.Title)
	assert.Equal(t, "Resource description", mapped.Description)
	assert.Equal(t, "article", string(mapped.Type))
	assert.Equal(t, "https://example.com/article", mapped.Url)
	assert.NotNil(t, mapped.ThumbnailUrl)
	assert.Equal(t, "https://example.com/thumb.png", *mapped.ThumbnailUrl)
	assert.Equal(t, []string{"api", "contract"}, mapped.Tags)
	assert.True(t, createdAt.Equal(mapped.CreatedAt))
	assert.True(t, updatedAt.Equal(mapped.UpdatedAt))
}

func TestMapResourceToContract_WithoutThumbnailOmitsJSONField(t *testing.T) {
	resource := models.Resource{
		ID:          "res-2",
		Title:       "No Thumb",
		Description: "No thumbnail resource",
		Type:        "video",
		URL:         "https://example.com/video",
		Tags:        []string{"video"},
	}

	mapped := mapResourceToContract(resource)
	assert.Nil(t, mapped.ThumbnailUrl)

	encoded, err := json.Marshal(mapped)
	assert.NoError(t, err)

	var payload map[string]any
	err = json.Unmarshal(encoded, &payload)
	assert.NoError(t, err)

	_, hasThumbnail := payload["thumbnailUrl"]
	assert.False(t, hasThumbnail)
}

func TestMapTagToContract(t *testing.T) {
	tag := models.Tag{Name: "golang", UsageCount: 8}
	mapped := mapTagToContract(tag)

	assert.Equal(t, "golang", mapped.Name)
	assert.Equal(t, 8, mapped.UsageCount)
}
