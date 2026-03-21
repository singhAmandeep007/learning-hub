package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"learninghub/constants"
	"learninghub/db"
	"learninghub/errors"
	"learninghub/middleware"
	"learninghub/models"
	"learninghub/pkg/logger"
	"learninghub/utils"
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
	product, exists := middleware.GetProductFromContext(c)

	if !exists {
		errors.RespondWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
		return
	}

	// Parse query parameters
	typeFilter := c.Query(constants.QueryParamType) // "video" | "pdf" | "article"
	tagsParam := c.Query(constants.QueryParamTags)  // "onboarding,tutorial" | "onboarding"
	search := c.Query(constants.QueryParamSearch)   // "getting%20started" | "v1.2"
	cursor := c.Query(constants.QueryParamCursor)
	limitStr := c.DefaultQuery(constants.QueryParamLimit, constants.DefaultLimitValue)

	// set to default page size if error in conversion or limit <= 0 or greator than max page size
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > constants.MaxPageSize {
		limit = constants.DefaultPageSize
	}

	// Create database services
	database := db.New()
	resourceService := db.NewResourceService(database)

	// Prepare query parameters
	var tags []string
	if tagsParam != "" {
		tags = utils.NormalizeTags(strings.Split(tagsParam, ","))
	}

	var validTypeFilter string
	if utils.IsValidResourceType(typeFilter) {
		validTypeFilter = typeFilter
	}

	// Execute query with limit + 1 to check for more results
	docs, err := resourceService.List(ctx, db.ResourceQuery{
		Product: product,
		Type:    validTypeFilter,
		Tags:    tags,
		Cursor:  cursor,
		Limit:   limit + 1,
	})
	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrQueryFailed, "Failed to fetch resources", err.Error())
		return
	}

	// Process results
	resources := make([]models.Resource, 0, len(docs))

	for _, doc := range docs {
		var resource models.Resource
		if err := doc.DataTo(&resource); err != nil {
			logger.Infof("Error converting document %s: %v", doc.Ref.ID, err)
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

		// Convert URLs to signed URLs before adding to response
		signedURL, signedThumbnailURL, err := utils.ConvertResourceURLsToSigned(
			ctx,
			resource.URL,
			resource.ThumbnailURL,
			constants.DefaultSignedURLExpiration,
		)
		if err != nil {
			logger.Infof("Error generating signed URLs for resource %s: %v", resource.ID, err)
			// Continue with original URLs if signing fails
		} else {
			resource.URL = signedURL
			resource.ThumbnailURL = signedThumbnailURL
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
	product, exists := middleware.GetProductFromContext(c)

	if !exists {
		errors.RespondWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
		return
	}

	// Create database services
	database := db.New()
	resourceService := db.NewResourceService(database)

	// Get document from product-specific collection
	doc, err := resourceService.GetByID(ctx, product, id)

	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrResourceNotFound, "Resource not found", err.Error())
		return
	}

	var resource models.Resource
	if err := doc.DataTo(&resource); err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrDataConversionFailed, "Failed to process resource data", err.Error())
		return
	}
	resource.ID = doc.Ref.ID

	// Convert URLs to signed URLs before returning
	signedURL, signedThumbnailURL, err := utils.ConvertResourceURLsToSigned(
		ctx,
		resource.URL,
		resource.ThumbnailURL,
		constants.DefaultSignedURLExpiration,
	)
	if err != nil {
		logger.Infof("Error generating signed URLs for resource %s: %v", resource.ID, err)
		// Continue with original URLs if signing fails
	} else {
		resource.URL = signedURL
		resource.ThumbnailURL = signedThumbnailURL
	}

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
	product, exists := middleware.GetProductFromContext(c)
	if !exists {
		errors.RespondWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
		return
	}

	contentType := c.GetHeader("Content-Type")

	// Check if it's actually a multipart form
	if contentType == "" || !strings.Contains(contentType, "multipart/form-data") {
		errors.RespondWithError(c, errors.ErrInvalidContentType, "Request must be multipart/form-data")
		return
	}

	// Check Content-Length before parsing multipart form
	contentLength := c.Request.ContentLength
	if contentLength > constants.MaxFileSize {
		errors.RespondWithError(c, errors.ErrFileTooLarge, fmt.Sprintf("Request size too large. Maximum size is %d MB", constants.MaxFileSize/(1<<20)))
		return
	}

	// Parse multipart form with MaxFileSize limit
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
		Title:       c.PostForm(constants.FormFieldTitle),
		Description: c.PostForm(constants.FormFieldDescription),
		Type:        c.PostForm(constants.FormFieldType),

		URL:          c.PostForm(constants.FormFieldURL),
		ThumbnailURL: c.PostForm(constants.FormFieldThumbnailURL),
		Tags:         utils.NormalizeTags(strings.Split(c.PostForm(constants.FormFieldTags), ",")),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Validate required fields
	if resource.Title == "" || resource.Description == "" || resource.Type == "" {
		errors.RespondWithError(c, errors.ErrMissingRequired, "Title, description, and type are required")
		return
	}

	// Validate resource type
	if !utils.IsValidResourceType(resource.Type) {
		errors.RespondWithError(c, errors.ErrUnsupportedType, "Type must be 'video', 'pdf', or 'article'")
		return
	}

	// Check if resource type is article AND url is not provided
	if resource.Type == constants.ResourceTypeArticle && resource.URL == "" {
		errors.RespondWithError(c, errors.ErrMissingRequired, "URL must be provided for 'article' type")
		return
	}

	// Handle file uploads for video and pdf types if url is not provided
	if (resource.Type == constants.ResourceTypeVideo || resource.Type == constants.ResourceTypePDF) && resource.URL == "" {
		file, header, err := c.Request.FormFile(constants.FormFieldFile)
		if err != nil {
			errors.RespondWithErrorDetails(c, errors.ErrMissingRequired, fmt.Sprintf("File is required for %s resources", resource.Type), err.Error())
			return
		}
		// successfully opened the file
		defer file.Close()

		// Upload file to Cloud Storage
		url, err := utils.UploadFile(ctx, file, header, product, resource.Type)
		if err != nil {
			// Check if this is a file validation error
			if strings.Contains(err.Error(), constants.ErrFileValidationFailed) {
				logger.Warnf("File validation failed: %v", err)
				errors.RespondWithErrorDetails(c, errors.ErrInvalidFileType, "Invalid file type", fileTypeErrorDetail(resource.Type))
				return
			}
			errors.RespondWithErrorDetails(c, errors.ErrUploadFailed, "Failed to upload file", err.Error())
			return
		}
		resource.URL = url.PublicURL
	}

	// Handle thumbnail upload if thumbnail url not provided
	if resource.ThumbnailURL == "" {
		// Handle thumbnail upload (optional)
		thumbnailFile, thumbnailHeader, err := c.Request.FormFile(constants.FormFieldThumbnail)
		if err == nil {
			defer thumbnailFile.Close()

			thumbnailURL, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, product, constants.ResourceTypeImage)

			if err != nil {
				// Check if this is a file validation error for thumbnail
				if strings.Contains(err.Error(), constants.ErrFileValidationFailed) {
					logger.Warnf("Thumbnail validation failed: %v", err)
					errors.RespondWithErrorDetails(c, errors.ErrInvalidFileType, "Invalid thumbnail file type", fileTypeErrorDetail(constants.ResourceTypeImage))
					return
				}
				logger.Infof("Failed to upload thumbnail: %v", err)
			} else {
				resource.ThumbnailURL = thumbnailURL.PublicURL
			}
		}
	}

	// Create database services
	database := db.New()
	resourceService := db.NewResourceService(database)

	// Save to Firestore in product-specific collection
	docRef, err := resourceService.Create(ctx, product, resource)
	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrMutationFailed, "Failed to save resource", err.Error())
		return
	}

	// Update tag usage counts
	utils.UpdateTagUsage(ctx, product, resource.Tags, 1)

	resource.ID = docRef.ID

	// Convert URLs to signed URLs before returning
	signedURL, signedThumbnailURL, err := utils.ConvertResourceURLsToSigned(
		ctx,
		resource.URL,
		resource.ThumbnailURL,
		constants.DefaultSignedURLExpiration,
	)
	if err != nil {
		logger.Infof("Error generating signed URLs for new resource: %v", err)
	} else {
		resource.URL = signedURL
		resource.ThumbnailURL = signedThumbnailURL
	}

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
	product, exists := middleware.GetProductFromContext(c)
	if !exists {
		errors.RespondWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
		return
	}

	// Create database services
	database := db.New()
	resourceService := db.NewResourceService(database)

	// Get existing resource from product-specific collection
	doc, err := resourceService.GetByID(ctx, product, id)
	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrResourceNotFound, "Resource not found", err.Error())
		return
	}

	var existingResource models.Resource
	if err := doc.DataTo(&existingResource); err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrDataConversionFailed, "Failed to process existing resource data", err.Error())
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

	if title, titleExists := c.GetPostForm(constants.FormFieldTitle); titleExists {
		updatedResource.Title = title
	}
	if description, descriptionExists := c.GetPostForm(constants.FormFieldDescription); descriptionExists {
		updatedResource.Description = description
	}
	if resourceType, typeExists := c.GetPostForm(constants.FormFieldType); typeExists {
		updatedResource.Type = resourceType
		// Validate resource type
		if !utils.IsValidResourceType(updatedResource.Type) {
			errors.RespondWithError(c, errors.ErrUnsupportedType, "Type must be 'video', 'pdf', or 'article'")
			return
		}

		// Check if trying to change resource type
		if existingResource.Type != updatedResource.Type {
			errors.RespondWithError(c, errors.ErrInvalidParam, "Resource type cannot be changed")
			return
		}
	}
	if tagsStr, tagsStrExists := c.GetPostForm(constants.FormFieldTags); tagsStrExists {
		oldTags = existingResource.Tags
		newTags = utils.NormalizeTags(strings.Split(tagsStr, ","))

		updatedResource.Tags = newTags
	}
	updatedResource.UpdatedAt = time.Now()

	// Handle URL and file updates
	urlFromForm, urlFromFormExists := c.GetPostForm(constants.FormFieldURL)
	_, fileExists := c.Request.MultipartForm.File[constants.FormFieldFile]

	// If both URL and file are provided, return error
	if urlFromFormExists && fileExists {
		errors.RespondWithError(c, errors.ErrInvalidParam, "Either provide url or file")
		return
	}

	if urlFromFormExists {
		if existingResource.URL != "" {
			if err = utils.DeleteFileFromURL(ctx, existingResource.URL); err != nil {
				logger.Infof("Failed to delete old file: %v", err)
			}
		}

		updatedResource.URL = urlFromForm
	}

	if fileExists && (existingResource.Type == constants.ResourceTypeVideo || existingResource.Type == constants.ResourceTypePDF) {
		// User provided a new file to upload
		if file, header, err := c.Request.FormFile(constants.FormFieldFile); err == nil {
			defer file.Close()

			// Delete old file if it was stored in our storage
			if existingResource.URL != "" {
				if err = utils.DeleteFileFromURL(ctx, existingResource.URL); err != nil {
					logger.Infof("Failed to delete old file: %v", err)
				}
			}

			// Upload new file
			uploadResult, err := utils.UploadFile(ctx, file, header, product, existingResource.Type)
			if err != nil {
				// Check if this is a file validation error
				if strings.Contains(err.Error(), constants.ErrFileValidationFailed) {
					logger.Warnf("File validation failed: %v", err)
					errors.RespondWithErrorDetails(c, errors.ErrInvalidFileType, "Invalid file type", fileTypeErrorDetail(existingResource.Type))
					return
				}
				errors.RespondWithErrorDetails(c, errors.ErrUploadFailed, "Failed to upload new file", err.Error())
				return
			}
			updatedResource.URL = uploadResult.PublicURL
		}
	}

	// Handle thumbnail URL and file updates
	thumbnailURLFromForm, thumbnailURLFromFormExists := c.GetPostForm(constants.FormFieldThumbnailURL)
	_, thumbnailFileExists := c.Request.MultipartForm.File[constants.FormFieldThumbnail]

	// If both thumbnail URL and thumbnail file are provided, return error
	if thumbnailURLFromFormExists && thumbnailFileExists {
		errors.RespondWithError(c, errors.ErrInvalidParam, "Either provide url or thumbnail")
		return
	}

	if thumbnailURLFromFormExists {
		// User provided a new thumbnail URL
		// Delete old thumbnail if it was stored in our storage
		if existingResource.ThumbnailURL != "" {
			if err = utils.DeleteFileFromURL(ctx, existingResource.ThumbnailURL); err != nil {
				logger.Infof("Failed to delete old thumbnail: %v", err)
			}
		}
		updatedResource.ThumbnailURL = thumbnailURLFromForm
	}

	if thumbnailFileExists {
		// User provided a new thumbnail file to upload
		if thumbnailFile, thumbnailHeader, err := c.Request.FormFile(constants.FormFieldThumbnail); err == nil {
			defer thumbnailFile.Close()

			// Delete old thumbnail if it was stored in our storage
			if existingResource.ThumbnailURL != "" {
				if err = utils.DeleteFileFromURL(ctx, existingResource.ThumbnailURL); err != nil {
					logger.Infof("Failed to delete old thumbnail: %v", err)
				}
			}

			// Upload new thumbnail
			thumbnailResult, err := utils.UploadFile(ctx, thumbnailFile, thumbnailHeader, product, constants.ResourceTypeImage)
			if err != nil {
				// Check if this is a file validation error for thumbnail
				if strings.Contains(err.Error(), constants.ErrFileValidationFailed) {
					logger.Warnf("Thumbnail validation failed: %v", err)
					errors.RespondWithErrorDetails(c, errors.ErrInvalidFileType, "Invalid thumbnail file type", fileTypeErrorDetail(constants.ResourceTypeImage))
					return
				}
				logger.Infof("Failed to upload thumbnail: %v", err)
			} else {
				updatedResource.ThumbnailURL = thumbnailResult.PublicURL
			}
		}
	}

	// Save updated resource to product-specific collection
	err = resourceService.Update(ctx, product, id, updatedResource)
	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrMutationFailed, "Failed to update resource", err.Error())
		return
	}

	// Update tag usage counts
	if len(oldTags) > 0 || len(newTags) > 0 {
		utils.UpdateTagUsage(ctx, product, oldTags, -1)
		utils.UpdateTagUsage(ctx, product, newTags, 1)
	}

	updatedResource.ID = id

	// Convert URLs to signed URLs before returning
	signedURL, signedThumbnailURL, err := utils.ConvertResourceURLsToSigned(
		ctx,
		updatedResource.URL,
		updatedResource.ThumbnailURL,
		constants.DefaultSignedURLExpiration,
	)
	if err != nil {
		logger.Infof("Error generating signed URLs for updated resource: %v", err)
	} else {
		updatedResource.URL = signedURL
		updatedResource.ThumbnailURL = signedThumbnailURL
	}

	c.JSON(http.StatusOK, updatedResource)
}

// DeleteResource handles DELETE /resource/:id
//   - Deletes a resource and associated files from storage.
func DeleteResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get product from context (validated by middleware)
	product, exists := middleware.GetProductFromContext(c)
	if !exists {
		errors.RespondWithError(c, errors.ErrInvalidProduct, "Invalid product parameter")
		return
	}

	// Create database services
	database := db.New()
	resourceService := db.NewResourceService(database)

	// Get existing resource to clean up files from product-specific collection
	doc, err := resourceService.GetByID(ctx, product, id)
	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrResourceNotFound, "Resource not found", err.Error())
		return
	}

	var resource models.Resource
	if err := doc.DataTo(&resource); err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrDataConversionFailed, "Failed to process resource data", err.Error())
		return
	}

	// Delete files from Cloud Storage
	if resource.URL != "" {
		if err := utils.DeleteFileFromURL(ctx, resource.URL); err != nil {
			logger.Infof("Failed to delete file: %v", err)
		}
	}
	if resource.ThumbnailURL != "" {
		if err := utils.DeleteFileFromURL(ctx, resource.ThumbnailURL); err != nil {
			logger.Infof("Failed to delete thumbnail: %v", err)
		}
	}

	// Update tag usage counts
	utils.UpdateTagUsage(ctx, product, resource.Tags, -1)

	// Delete from Firestore product-specific collection
	err = resourceService.Delete(ctx, product, id)
	if err != nil {
		errors.RespondWithErrorDetails(c, errors.ErrMutationFailed, "Failed to delete resource", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// handleMultipartFormError handles errors from ParseMultipartForm
//   - returns appropriate error response
func handleMultipartFormError(c *gin.Context, err error) {
	logger.Infof("ParseMultipartForm error: %v", err)

	if strings.Contains(err.Error(), "too large") {
		errors.RespondWithErrorDetails(c, errors.ErrFileTooLarge, fmt.Sprintf("File too large. Maximum size is %d MB", constants.MaxFileSize/(1<<20)), err.Error())
	} else if strings.Contains(err.Error(), "no multipart boundary") {
		errors.RespondWithErrorDetails(c, errors.ErrInvalidContentType, "Invalid multipart form data - no boundary found", err.Error())
	} else {
		errors.RespondWithErrorDetails(c, errors.ErrInvalidPayload, fmt.Sprintf("Failed to parse form data: %v", err), err.Error())
	}
}

// fileTypeErrorDetail returns a user-friendly error message for file validation failures
// based on the resource type.
func fileTypeErrorDetail(resourceType string) string {
	switch resourceType {
	case constants.ResourceTypeVideo:
		return "Only MP4 and WebM video formats are supported"
	case constants.ResourceTypePDF:
		return "Only PDF files are supported"
	case constants.ResourceTypeImage:
		return "The uploaded file is not a supported image format"
	default:
		return "The uploaded file type is not supported"
	}
}
