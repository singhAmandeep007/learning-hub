package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"learning-hub/constants"
	"learning-hub/firebase"
	"learning-hub/models"
	"learning-hub/utils"
)

// GetResources handles GET /api/resources
func GetResources(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	typeFilter := c.DefaultQuery("type", "all") // "all" | "video" | "pdf" | "article"
	tagsParam := c.Query("tags") // "onboarding,tutorial" | "onboarding"
	search := c.Query("search") // "getting%20started" | "v1.2"
	cursor := c.Query("cursor")
	limitStr := c.DefaultQuery("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > constants.MaxPageSize {
		limit = constants.DefaultPageSize
	}

	// Build Firestore query
	query := firebase.FirestoreClient.Collection(constants.CollectionResources).OrderBy("createdAt", firestore.Desc)

	// Apply type filter
	if typeFilter != "all" {
		query = query.Where("type", "==", typeFilter)
	}

	// Apply tags filter
	if tagsParam != "" {
		tags := utils.NormalizeTags(strings.Split(tagsParam, ","))
		if len(tags) > 0 {
			query = query.Where("tags", "array-contains-any", tags)
		}
	}

	// Apply cursor for pagination
	if cursor != "" {
		// In a real implementation, you'd decode the cursor to get the document snapshot
		// For simplicity, we'll use offset here (not recommended for large datasets)
		offset, _ := strconv.Atoi(cursor)
		query = query.Offset(offset)
	}

	// Execute query
	docs, err := query.Limit(limit + 1).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "query_failed",
			Message: "Failed to fetch resources",
		})
		return
	}

	// Process results
	resources := make([]models.Resource, 0, len(docs))
	for i, doc := range docs {
		if i >= limit { // We fetched one extra to check if there are more
			break
		}

		var resource models.Resource
		if err := doc.DataTo(&resource); err != nil {
			log.Printf("Error converting document %s: %v", doc.Ref.ID, err)
			continue
		}
		resource.ID = doc.Ref.ID

		// Apply search filter
		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(resource.Title), searchLower) &&
				!strings.Contains(strings.ToLower(resource.Description), searchLower) {
				continue
			}
		}

		resources = append(resources, resource)
	}

	// Prepare response
	response := models.PaginatedResponse{
		Data:    resources,
		HasMore: len(docs) > limit,
	}

	if response.HasMore {
		// Simple cursor implementation using offset
		currentOffset, _ := strconv.Atoi(cursor)
		response.NextCursor = strconv.Itoa(currentOffset + limit)
	}

	c.JSON(http.StatusOK, response)
}

// GetResource handles GET /api/resources/:id
func GetResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	doc, err := firebase.FirestoreClient.Collection(constants.CollectionResources).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound,  models.ErrorResponse{
			Error:   "resource_not_found",
			Message: "Resource not found",
		})
		return
	}

	var resource models.Resource
	if err := doc.DataTo(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "data_conversion_failed",
			Message: "Failed to process resource data",
		})
		return
	}
	resource.ID = doc.Ref.ID

	c.JSON(http.StatusOK, resource)
}

// CreateResource handles POST /api/resources
func CreateResource(c *gin.Context) {
	ctx := c.Request.Context()

	// Log the content type for debugging
	log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))
	log.Printf("Content-Length: %s", c.GetHeader("Content-Length"))

	// Check if it's actually a multipart form
	contentType := c.GetHeader("Content-Type")
	if contentType == "" || !strings.Contains(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_content_type",
			Message: "Request must be multipart/form-data",
		})
		return
	}

	// error handling for max memory
	if err := c.Request.ParseMultipartForm(constants.MaxFileSize); err != nil {
		handleMultipartFormError(c, err)
		return
	}

	log.Printf("Form values: %+v", c.Request.MultipartForm.Value)
	
	if c.Request.MultipartForm.File != nil {
		for key, files := range c.Request.MultipartForm.File {
			log.Printf("File field '%s': %d files", key, len(files))
			for i, file := range files {
				log.Printf("  File %d: %s (size: %d)", i, file.Filename, file.Size)
			}
		}
	}


	// Extract form fields
	resource := models.Resource{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Type:        c.PostForm("type"),
		URL:         c.PostForm("url"),
		ThumbnailURL: c.PostForm("thumbnailUrl"),
		Tags:        utils.NormalizeTags(strings.Split(c.PostForm("tags"), ",")),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validate required fields
	if resource.Title == "" || resource.Description == "" || resource.Type == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_required_fields",
			Message: "Title, description, and type are required",
		})
		return
	}

	// Validate resource type
	if resource.Type != constants.ResourceTypeVideo && resource.Type != constants.ResourceTypePDF && resource.Type != constants.ResourceTypeArticle {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_type",
			Message: "Type must be 'video', 'pdf', or 'article'",
		})
		return
	}

	// Check if resource type is article and url is not provided
	if resource.Type == constants.ResourceTypeArticle && resource.URL == ""  {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "missing_url",
			Message: "Url must be provided for 'article'",
		})
		return
	}

	// Handle file uploads for video and pdf types and check if url is not provided
	if (resource.Type == constants.ResourceTypeVideo || resource.Type == constants.ResourceTypePDF) && resource.URL == "" {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "file_required",
				Message: fmt.Sprintf("File is required for %s resources", resource.Type),
			})
			return
		}
		defer file.Close()

		// Upload file to Cloud Storage
		url, err := utils.UploadFile(ctx, file, header, resource.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "upload_failed",
				Message: "Failed to upload file",
			})
			return
		}
		resource.URL = url.PublicURL
	}

	// Handle thumbnail upload and check if thumbnail url not provided
	if resource.ThumbnailURL == "" {
		// Handle thumbnail upload (optional)
		thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail")
		if err == nil {
			defer thumbnailFile.Close()
			thumbnailURL, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, "image")
			if err != nil {
				log.Printf("Failed to upload thumbnail: %v", err)
			} else {
				resource.ThumbnailURL = thumbnailURL.PublicURL
			}
		}
	}


	// Save to Firestore
	docRef, _, err := firebase.FirestoreClient.Collection(constants.CollectionResources).Add(ctx, resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "save_failed",
			Message: "Failed to save resource",
		})
		return
	}

	// Update tag usage counts
	utils.UpdateTagUsage(ctx, resource.Tags, 1)

	resource.ID = docRef.ID
	c.JSON(http.StatusCreated, resource)
}

// UpdateResource handles PUT /api/resources/:id
func UpdateResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get existing resource
	doc, err := firebase.FirestoreClient.Collection(constants.CollectionResources).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "resource_not_found",
			Message: "Resource not found",
		})
		return
	}

	var existingResource models.Resource
	if err := doc.DataTo(&existingResource); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "data_conversion_failed",
			Message: "Failed to process existing resource data",
		})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(constants.MaxFileSize); err != nil {
		handleMultipartFormError(c, err)
		return
	}

	// Update fields
	updatedResource := existingResource
	if title := c.PostForm("title"); title != "" {
		updatedResource.Title = title
	}
	if description := c.PostForm("description"); description != "" {
		updatedResource.Description = description
	}
	if resourceType := c.PostForm("type"); resourceType != "" {
		updatedResource.Type = resourceType
	}
	if url := c.PostForm("url"); url != "" {
		updatedResource.URL = url
	}
	if thumbnailUrl := c.PostForm("thumbnailUrl"); thumbnailUrl != ""{
		updatedResource.ThumbnailURL = thumbnailUrl
	}
	if tagsStr := c.PostForm("tags"); tagsStr != "" {
		oldTags := existingResource.Tags
		newTags := utils.NormalizeTags(strings.Split(tagsStr, ","))
		updatedResource.Tags = newTags
		
		// Update tag usage counts
		utils.UpdateTagUsage(ctx, oldTags, -1)
		utils.UpdateTagUsage(ctx, newTags, 1)
	}
	updatedResource.UpdatedAt = time.Now()

	// Handle file replacement if resource is of type video or pdf
	if (updatedResource.Type == constants.ResourceTypeVideo || updatedResource.Type == constants.ResourceTypePDF) && updatedResource.URL == "" {
		// Handle file replacement
		if file, header, err := c.Request.FormFile("file"); err == nil {
			defer file.Close()

			// Delete old file if it exists
			if existingResource.URL != "" {
				if err = utils.DeleteFileFromURL(ctx, existingResource.URL); err != nil {
					log.Printf("Failed to delete file: %v", err)
				}
			}

			// Upload new file
			url, err := utils.UploadFile(ctx, file, header, updatedResource.Type)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error:   "upload_failed",
					Message: "Failed to upload new file",
				})
				return
			}
			updatedResource.URL = url.PublicURL
		}
	}

	if updatedResource.ThumbnailURL == "" {
		// Handle thumbnail replacement
		if thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail"); err == nil {
			defer thumbnailFile.Close()

			// Delete old thumbnail
			if existingResource.ThumbnailURL != "" {
				log.Printf("existingResource.ThumbnailURL: %v", existingResource.ThumbnailURL)
				if err =  utils.DeleteFileFromURL(ctx, existingResource.ThumbnailURL); err != nil {
					log.Printf("Failed to delete thumbnail: %v", err)
				}
			}

			// Upload new thumbnail
			thumbnailURL, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, "image")
			if err != nil {
				log.Printf("Failed to upload thumbnail: %v", err)
			} else {
				updatedResource.ThumbnailURL = thumbnailURL.PublicURL
			}
		}
	}

	// Save updated resource
	_, err = firebase.FirestoreClient.Collection(constants.CollectionResources).Doc(id).Set(ctx, updatedResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "update_failed",
			Message: "Failed to update resource",
		})
		return
	}

	updatedResource.ID = id
	c.JSON(http.StatusOK, updatedResource)
}

// DeleteResource handles DELETE /api/resource/:id
func DeleteResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get existing resource to clean up files
	doc, err := firebase.FirestoreClient.Collection(constants.CollectionResources).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "resource_not_found",
			Message: "Resource not found",
		})
		return
	}

	var resource models.Resource
	if err := doc.DataTo(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "data_conversion_failed",
			Message: "Failed to process resource data",
		})
		return
	}

	// Delete files from Cloud Storage
	if resource.URL != "" {
		utils.DeleteFileFromURL(ctx, resource.URL)
	}
	if resource.ThumbnailURL != "" {
		utils.DeleteFileFromURL(ctx, resource.ThumbnailURL)
	}

	// Update tag usage counts
	utils.UpdateTagUsage(ctx, resource.Tags, -1)

	// Delete from Firestore
	_, err = firebase.FirestoreClient.Collection(constants.CollectionResources).Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "delete_failed",
			Message: "Failed to delete resource",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// handleMultipartFormError handles errors from ParseMultipartForm and returns appropriate error response
func handleMultipartFormError(c *gin.Context, err error) {
	log.Printf("ParseMultipartForm error: %v", err)
	
	var message string
	if strings.Contains(err.Error(), "too large") {
		message = fmt.Sprintf("File too large. Maximum size is %d MB", constants.MaxFileSize/(1<<20))
	} else if strings.Contains(err.Error(), "no multipart boundary") {
		message = "Invalid multipart form data - no boundary found"
	} else {
		message = fmt.Sprintf("Failed to parse form data: %v", err)
	}
	
	c.JSON(http.StatusBadRequest, models.ErrorResponse{
		Error:   "form_parse_error",
		Message: message,
	})
}