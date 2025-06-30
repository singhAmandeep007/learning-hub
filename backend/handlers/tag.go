package handlers

import (
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"learning-hub/constants"
	"learning-hub/firebase"
	"learning-hub/models"
)

// GetTags handles GET /tags
func GetTags(c *gin.Context) {
	ctx := c.Request.Context()

	docs, err := firebase.FirestoreClient.Collection(constants.CollectionTags).OrderBy("usageCount", firestore.Desc).Documents(ctx).GetAll()

	if err != nil {
		log.Printf("Error fetching tags from Firestore: %v\n", err)

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   constants.QueryFailed,
			Message: "Failed to fetch tags",
		})
		return
	}

	tags := make([]models.Tag, 0, len(docs))
	for _, doc := range docs {
		var tag models.Tag
		if err := doc.DataTo(&tag); err != nil {
			log.Printf("Warning: Failed to unmarshal tag document ID %s: %v\n", doc.Ref.ID, err)
			continue
		}
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, tags)
}
