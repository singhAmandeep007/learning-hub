package handlers

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"learning-hub/config"
	"learning-hub/constants"
	"learning-hub/models"
)

// GetTags handles GET /api/tags
func GetTags(c *gin.Context) {
	ctx := c.Request.Context()

	docs, err := config.FirestoreClient.Collection(constants.CollectionTags).OrderBy("usageCount", firestore.Desc).Documents(ctx).GetAll()

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "query_failed",
			Message: "Failed to fetch tags",
		})
		return
	}

	tags := make([]models.Tag, 0, len(docs))
	for _, doc := range docs {
		var tag models.Tag
		if err := doc.DataTo(&tag); err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, tags)
}

