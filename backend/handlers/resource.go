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

	"learning-hub/config"
	"learning-hub/constants"
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
	query := config.FirestoreClient.Collection(constants.CollectionResources).OrderBy("createdAt", firestore.Desc)

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

		// Apply search filter (simple implementation)
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

	doc, err := config.FirestoreClient.Collection(constants.CollectionResources).Doc(id).Get(ctx)
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


	// Set max multipart memory (this replaces ParseMultipartForm)
	// c.Request.ParseMultipartForm(MaxFileSize)
	

		// Check if it's actually a multipart form
		contentType := c.GetHeader("Content-Type")
		if contentType == "" || !strings.Contains(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "invalid_content_type",
				Message: "Request must be multipart/form-data",
			})
			return
		}

		if err := c.Request.ParseMultipartForm(constants.MaxFileSize); err != nil {
			log.Printf("ParseMultipartForm error: %v", err)
			
			// Provide more specific error messages
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
	if resource.Type != "video" && resource.Type != "pdf" && resource.Type != "article" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_type",
			Message: "Type must be 'video', 'pdf', or 'article'",
		})
		return
	}

	// Handle file uploads for video and pdf types
	if resource.Type == "video" || resource.Type == "pdf" {
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
		resource.URL = url
	}

	// Handle thumbnail upload (optional)
	thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail")
	if err == nil {
		defer thumbnailFile.Close()
		thumbnailURL, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, "image")
		if err != nil {
			log.Printf("Failed to upload thumbnail: %v", err)
		} else {
			resource.ThumbnailURL = thumbnailURL
		}
	}

	// Save to Firestore
	docRef, _, err := config.FirestoreClient.Collection(constants.CollectionResources).Add(ctx, resource)
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