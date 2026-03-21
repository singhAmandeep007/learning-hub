package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"learninghub/db"
	"learninghub/errors"
	"learninghub/middleware"
	"learninghub/models"
	"learninghub/pkg/logger"
)

// GetTags handles GET /tags
func GetTags(c *gin.Context) {
	ctx := c.Request.Context()

	// Get product from context (validated by middleware)
	product, exists := middleware.GetProductFromContext(c)
	if !exists {
		errors.RespondWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
		return
	}

	database := db.New()
	tagService := db.NewTagService(database)

	// Get tags from product-specific collection
	docs, err := tagService.List(ctx, product)
	if err != nil {
		logger.Infof("Error fetching tags from database: %v\n", err)
		errors.RespondWithError(c, errors.ErrQueryFailed, "Failed to fetch tags")
		return
	}

	tags := make([]models.Tag, 0, len(docs))
	for _, doc := range docs {
		var tag models.Tag
		if err := doc.DataTo(&tag); err != nil {
			logger.Infof("Warning: Failed to unmarshal tag document ID %s: %v\n", doc.Ref.ID, err)
			continue
		}
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, tags)
}
