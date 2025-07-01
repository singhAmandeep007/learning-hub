package handlers

import (
	"encoding/json"
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
	"learning-hub/middleware"
	"learning-hub/models"
	"learning-hub/utils"
)

// GetResources handles GET /resources
// Supports filtering by type, tags, and search, as well as pagination using cursor and limit.
// Query Params:
//   - type: Filter by resource type ("video", "pdf", "article")
//   - tags: Comma-separated list of tags to filter by
//   - search: Search string for title/description
//   - cursor: Offset for pagination (as stringified int)
//   - limit: Number of items per page (default 20, max 100)
func GetResources(c *gin.Context) {
	ctx := c.Request.Context()

	// Get product from context (validated by middleware)
	product := middleware.GetProductFromContext(c)

	// Parse query parameters
	typeFilter := c.Query("type") // "video" | "pdf" | "article"
	tagsParam := c.Query("tags")  // "onboarding,tutorial" | "onboarding"
	search := c.Query("search")   // "getting%20started" | "v1.2"
	cursor := c.Query("cursor")
	limitStr := c.DefaultQuery("limit", "20")

	// set to default page size if error in conversion or limit <= 0 or greator than max page size
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > constants.MaxPageSize {
		limit = constants.DefaultPageSize
	}

	// Build Firestore query with product-specific collection
	query := firebase.FirestoreClient.Collection(constants.GetResourcesCollectionName(product)).OrderBy("createdAt", firestore.Desc)

	// Apply type filter
	if utils.IsValidResourceType(typeFilter) {
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
		offset, err := strconv.Atoi(cursor)
		if err == nil && offset >= 0 {
			// skips the records
			query = query.Offset(offset)
		}
	}

	// Execute query with limit
	docs, err := query.Limit(limit + 1).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   fmt.Sprintf("%s:%v", constants.QueryFailed, err),
			Message: "Failed to fetch resources",
		})
		return
	}

	// Process results
	resources := make([]models.Resource, 0, len(docs))

	for _, doc := range docs {
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

		// Only add if we haven't reached the limit
		if len(resources) < limit {
			resources = append(resources, resource)
		}
	}

	// We have more if we fetched more documents than our limit
	// AND we have exactly 'limit' resources after filtering
	hasMore := len(docs) > limit && len(resources) == limit

	response := models.PaginatedResponse{
		Data:    resources,
		HasMore: hasMore,
	}

	// Set next cursor only if there are more items
	if hasMore {
		currentOffset := 0
		if cursor != "" {
			currentOffset, _ = strconv.Atoi(cursor)
		}
		response.NextCursor = strconv.Itoa(currentOffset + limit)
	}

	c.JSON(http.StatusOK, response)
}

// GetResource handles GET /resources/:id
func GetResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get product from context (validated by middleware)
	product := middleware.GetProductFromContext(c)

	// Get document from product-specific collection
	collectionName := constants.GetResourcesCollectionName(product)
	doc, err := firebase.FirestoreClient.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   fmt.Sprintf("%s:%v", constants.NotFound, err),
			Message: "Resource not found",
		})
		return
	}

	var resource models.Resource
	if err := doc.DataTo(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   fmt.Sprintf("%s:%v", err, constants.DataConversionFailed),
			Message: "Failed to process resource data",
		})
		return
	}
	resource.ID = doc.Ref.ID

	c.JSON(http.StatusOK, resource)
}

// CreateResource handles POST /resources
//   - Accepts multipart/form-data for resource creation.
//   - Required fields: title, description, type.
//   - For "video" and "pdf" types, if url provided in the request, it will be prioritized and used as the resource's URL.
//     even if a file is uploaded.
//   - For "article" type, url is required.
func CreateResource(c *gin.Context) {
	ctx := c.Request.Context()

	// Get product from context (validated by middleware)
	product := middleware.GetProductFromContext(c)

	contentType := c.GetHeader("Content-Type")

	// Check if it's actually a multipart form
	if contentType == "" || !strings.Contains(contentType, "multipart/form-data") {
		c.JSON(http.StatusBadGateway, models.ErrorResponse{
			Error:   constants.InvalidContentType,
			Message: "Request must be multipart/form-data",
		})
		return
	}

	// error handling for max memory
	if err := c.Request.ParseMultipartForm(constants.MaxFileSize); err != nil {
		handleMultipartFormError(c, err)
		return
	}

	// log the files for debugging
	// if c.Request.MultipartForm.File != nil {
	// 	for key, files := range c.Request.MultipartForm.File {
	// 		log.Printf("File field '%s': %d files", key, len(files))
	// 		for i, file := range files {
	// 			log.Printf("  File %d: %s (size: %d)", i, file.Filename, file.Size)
	// 		}
	// 	}
	// }

	// Extract form fields
	resource := models.Resource{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Type:        c.PostForm("type"),

		URL:          c.PostForm("url"),
		ThumbnailURL: c.PostForm("thumbnailUrl"),
		Tags:         utils.NormalizeTags(strings.Split(c.PostForm("tags"), ",")),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Validate required fields
	if resource.Title == "" || resource.Description == "" || resource.Type == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.InvalidPayload,
			Message: "Title, description, and type are required",
		})
		return
	}

	// Validate resource type
	if !utils.IsValidResourceType(resource.Type) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.InvalidPayload,
			Message: "Type must be 'video', 'pdf', or 'article'",
		})
		return
	}

	// Check if resource type is article AND url is not provided
	if resource.Type == constants.ResourceTypeArticle && resource.URL == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.InvalidPayload,
			Message: "Url must be provided for 'article'",
		})
		return
	}

	// Handle file uploads for video and pdf types if url is not provided
	if (resource.Type == constants.ResourceTypeVideo || resource.Type == constants.ResourceTypePDF) && resource.URL == "" {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   constants.InvalidPayload,
				Message: fmt.Sprintf("File is required for %s resources", resource.Type),
			})
			return
		}
		// successfully opened the file
		defer file.Close()

		// Upload file to Cloud Storage
		url, err := utils.UploadFile(ctx, file, header, product, resource.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   constants.UploadFailed,
				Message: "Failed to upload file",
			})
			return
		}
		resource.URL = url.PublicURL
	}

	// Handle thumbnail upload if thumbnail url not provided
	if resource.ThumbnailURL == "" {
		// Handle thumbnail upload (optional)
		thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail")
		if err == nil {
			defer thumbnailFile.Close()

			thumbnailURL, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, product, "image")

			if err != nil {
				log.Printf("Failed to upload thumbnail: %v", err)
			} else {
				resource.ThumbnailURL = thumbnailURL.PublicURL
			}
		}
	}

	// Save to Firestore in product-specific collection
	collectionName := constants.GetResourcesCollectionName(product)
	docRef, _, err := firebase.FirestoreClient.Collection(collectionName).Add(ctx, resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   fmt.Sprintf("%s:%v", constants.MutationFailed, err),
			Message: "Failed to save resource",
		})
		return
	}

	// Update tag usage counts
	utils.UpdateTagUsage(ctx, product, resource.Tags, 1)

	resource.ID = docRef.ID
	c.JSON(http.StatusCreated, resource)
}

// UpdateResource handles PATCH /resources/:id
//   - Accepts multipart/form-data for resource update.
//   - Only allows updating fields except for resource type (cannot be changed).
//   - Handles file and thumbnail replacement if provided.
func UpdateResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get product from context (validated by middleware)
	product := middleware.GetProductFromContext(c)

	// Get existing resource from product-specific collection
	collectionName := constants.GetResourcesCollectionName(product)
	doc, err := firebase.FirestoreClient.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   constants.NotFound,
			Message: "Resource not found",
		})
		return
	}

	var existingResource models.Resource
	if err := doc.DataTo(&existingResource); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   constants.DataConversionFailed,
			Message: "Failed to process existing resource data",
		})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(constants.MaxFileSize); err != nil {
		handleMultipartFormError(c, err)
		return
	}

	var oldTags []string
	var newTags []string

	var updatedResource models.Resource
	bytes, _ := json.Marshal(existingResource)
	json.Unmarshal(bytes, &updatedResource)

	if title, titleExists := c.GetPostForm("title"); titleExists {
		updatedResource.Title = title
	}
	if description, descriptionExists := c.GetPostForm("description"); descriptionExists {
		updatedResource.Description = description
	}
	if resourceType, typeExists := c.GetPostForm("type"); typeExists {
		updatedResource.Type = resourceType
		// Validate resource type
		if !utils.IsValidResourceType(updatedResource.Type) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   constants.InvalidPayload,
				Message: "Type must be 'video', 'pdf', or 'article'",
			})
			return
		}

		// Check if trying to change resource type
		if existingResource.Type != updatedResource.Type {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   constants.InvalidPayload,
				Message: "Resource type cannot be changed",
			})
			return
		}
	}
	if tagsStr, tagsStrExists := c.GetPostForm("tags"); tagsStrExists {
		oldTags = existingResource.Tags
		newTags = utils.NormalizeTags(strings.Split(tagsStr, ","))

		updatedResource.Tags = newTags
	}
	updatedResource.UpdatedAt = time.Now()

	// Handle URL and file updates
	urlFromForm, urlFromFormExists := c.GetPostForm("url")
	_, fileExists := c.Request.MultipartForm.File["file"]

	// If both URL and file are provided, return error
	if urlFromFormExists && fileExists {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.InvalidPayload,
			Message: "Either provide url or file",
		})
		return
	}

	if urlFromFormExists {
		if existingResource.URL != "" {
			if err = utils.DeleteFileFromURL(ctx, existingResource.URL); err != nil {
				log.Printf("Failed to delete old file: %v", err)
			}
		}

		updatedResource.URL = urlFromForm
	}

	if fileExists && (existingResource.Type == constants.ResourceTypeVideo || existingResource.Type == constants.ResourceTypePDF) {
		// User provided a new file to upload
		if file, header, err := c.Request.FormFile("file"); err == nil {
			defer file.Close()

			// Delete old file if it was stored in our storage
			if existingResource.URL != "" {
				if err = utils.DeleteFileFromURL(ctx, existingResource.URL); err != nil {
					log.Printf("Failed to delete old file: %v", err)
				}
			}

			// Upload new file
			uploadResult, err := utils.UploadFile(ctx, file, header, product, existingResource.Type)
			if err != nil {
				c.JSON(http.StatusInternalServerError, models.ErrorResponse{
					Error:   constants.UploadFailed,
					Message: "Failed to upload new file",
				})
				return
			}
			updatedResource.URL = uploadResult.PublicURL
		}
	}

	// Handle thumbnail URL and file updates
	thumbnailURLFromForm, thumbnailURLFromFormExists := c.GetPostForm("thumbnailUrl")
	_, thumbnailFileExists := c.Request.MultipartForm.File["thumbnail"]

	// If both thumbnail URL and thumbnail file are provided, return error
	if thumbnailURLFromFormExists && thumbnailFileExists {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   constants.InvalidPayload,
			Message: "Either provide url or thumbnail",
		})
		return
	}

	if thumbnailURLFromFormExists {
		// User provided a new thumbnail URL
		// Delete old thumbnail if it was stored in our storage
		if existingResource.ThumbnailURL != "" {
			if err = utils.DeleteFileFromURL(ctx, existingResource.ThumbnailURL); err != nil {
				log.Printf("Failed to delete old thumbnail: %v", err)
			}
		}
		updatedResource.ThumbnailURL = thumbnailURLFromForm
	}

	if thumbnailFileExists {
		// User provided a new thumbnail file to upload
		if thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail"); err == nil {
			defer thumbnailFile.Close()

			// Delete old thumbnail if it was stored in our storage
			if existingResource.ThumbnailURL != "" {
				if err = utils.DeleteFileFromURL(ctx, existingResource.ThumbnailURL); err != nil {
					log.Printf("Failed to delete old thumbnail: %v", err)
				}
			}

			// Upload new thumbnail
			thumbnailResult, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, product, "image")
			if err != nil {
				log.Printf("Failed to upload thumbnail: %v", err)
			} else {
				updatedResource.ThumbnailURL = thumbnailResult.PublicURL
			}
		}
	}

	// Save updated resource to product-specific collection
	collectionName = constants.GetResourcesCollectionName(product)
	_, err = firebase.FirestoreClient.Collection(collectionName).Doc(id).Set(ctx, updatedResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   fmt.Sprintf("%s,%v", constants.MutationFailed, err),
			Message: "Failed to update resource",
		})
		return
	}

	// Update tag usage counts
	if len(oldTags) > 0 || len(newTags) > 0 {
		utils.UpdateTagUsage(ctx, product, oldTags, -1)
		utils.UpdateTagUsage(ctx, product, newTags, 1)
	}

	updatedResource.ID = id
	c.JSON(http.StatusOK, updatedResource)
}

// DeleteResource handles DELETE /resource/:id
//   - Deletes a resource and associated files from storage.
func DeleteResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get product from context (validated by middleware)
	product := middleware.GetProductFromContext(c)

	// Get existing resource to clean up files from product-specific collection
	collectionName := constants.GetResourcesCollectionName(product)
	doc, err := firebase.FirestoreClient.Collection(collectionName).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   constants.NotFound,
			Message: "Resource not found",
		})
		return
	}

	var resource models.Resource
	if err := doc.DataTo(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   fmt.Sprintf("%s,%v", constants.DataConversionFailed, err),
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
	utils.UpdateTagUsage(ctx, product, resource.Tags, -1)

	// Delete from Firestore product-specific collection
	collectionName = constants.GetResourcesCollectionName(product)
	_, err = firebase.FirestoreClient.Collection(collectionName).Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   fmt.Sprintf("%s:%v", constants.MutationFailed, err),
			Message: "Failed to delete resource",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// handleMultipartFormError handles errors from ParseMultipartForm
//   - returns appropriate error response
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
		Error:   constants.InvalidPayload,
		Message: message,
	})
}
