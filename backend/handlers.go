package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	firebaseStorage "firebase.google.com/go/v4/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

// App holds the application dependencies
type App struct {
	firestoreClient *firestore.Client
	storageClient   *firebaseStorage.Client
	
	storageBucket   string
	adminSecret     string
}

// Constants
const (
	CollectionResources = "resources"
	CollectionTags      = "tags"
	DefaultPageSize     = 20
	MaxPageSize         = 100
	MaxFileSize         = 100 << 20 // 100MB
)

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

// Tag represents a tag with usage statistics
type Tag struct {
	Name       string `json:"name" firestore:"name"`
	UsageCount int    `json:"usageCount" firestore:"usageCount"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data        []Resource `json:"data"`
	NextCursor  string     `json:"nextCursor,omitempty"`
	HasMore     bool       `json:"hasMore"`
	Total       int        `json:"total,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// InitializeFirebase initializes Firebase services
func initializeFirebase() (*App, error) {
	ctx := context.Background()

	isUseFirebaseEmulator, err := strconv.ParseBool(os.Getenv("IS_FIREBASE_EMULATOR"))
	if err != nil {
		// Log a warning or fatal error if the env var is set but invalid
		log.Printf("Warning: Invalid value for IS_FIREBASE_EMULATOR ('%v'). Defaulting to false (not using emulator). Error: %v\n", isUseFirebaseEmulator, err)
		isUseFirebaseEmulator = false // Default to not using emulator if parsing fails
	}

	var opts []option.ClientOption
	var effectiveStorageBucket string

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		// Fallback or error if project ID is crucial and not set
		log.Fatalf("FIREBASE_PROJECT_ID environment variable is not set.")
	}


	if isUseFirebaseEmulator {
		opts = append(opts, option.WithoutAuthentication())
		// Although not neccessary - https://firebase.google.com/docs/emulator-suite/connect_firestore#admin_sdks
		if err := os.Setenv("FIRESTORE_EMULATOR_HOST", os.Getenv("FIRESTORE_EMULATOR_HOST")); err != nil {
			log.Fatalf("Init emulator db err: %v", err)
		}
		// Although not neccessary - https://firebase.google.com/docs/emulator-suite/connect_storage#web
		if err := os.Setenv("FIREBASE_STORAGE_EMULATOR_HOST", os.Getenv("FIREBASE_STORAGE_EMULATOR_HOST")); err != nil {
			log.Fatalf("Init emulator db err: %v", err)
		}

		effectiveStorageBucket = projectID + ".appspot.com"
	}else {
		opts = append(opts, option.WithCredentialsFile("learning-hub-81cc6-firebase-adminsdk-fbsvc-f78df0ee32.json"))

		effectiveStorageBucket = os.Getenv("FIREBASE_STORAGE_BUCKET")
		if effectiveStorageBucket == "" {
			// Default to <project_id>.firebasestorage.com for live environment if not specified
			effectiveStorageBucket = projectID + ".firebasestorage.com"
			log.Printf("FIREBASE_STORAGE_BUCKET not set, defaulting to: %s for live environment", effectiveStorageBucket)
		}
	}

	conf := &firebase.Config{
		ProjectID: projectID,
				StorageBucket: effectiveStorageBucket, // This is mainly for firebaseApp.Storage(ctx).DefaultBucket()

	}
	
	firebaseApp, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Initialize Firestore
	firestoreClient, err := firebaseApp.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore: %v", err)
	}

	// Initialize Cloud Storage
	storageClient, err := firebaseApp.Storage(ctx) // storage.NewClient(ctx,opts...)
	if err != nil {
		return nil, fmt.Errorf("error initializing storage: %v", err)
	}

	adminSecret := os.Getenv("ADMIN_SECRET")
	if adminSecret == "" {
		adminSecret = "your-admin-secret-key" // Should be set via environment variable
	}

	return &App{
		firestoreClient: firestoreClient,
		storageClient:   storageClient,
		storageBucket:   effectiveStorageBucket,
		adminSecret:     adminSecret,
	}, nil
}

// cleanup closes all clients
func (app *App) cleanup() {
	if app.firestoreClient != nil {
		app.firestoreClient.Close()
	}
	if app.storageClient != nil {
		// app.storageClient.Close()
	}
}

// adminAuth middleware for admin-only routes
func (app *App) adminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			// Expected format: "Bearer SECRET_KEY"
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" && parts[1] == app.adminSecret {
				c.Next()
				return
			}
		}

		// Check query parameter as fallback
		secret := c.Query("admin_secret")
		if secret == app.adminSecret {
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "Admin authentication required",
		})
		c.Abort()
	}
}

// getResources handles GET /api/resources
func (app *App) getResources(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse query parameters
	typeFilter := c.DefaultQuery("type", "all") // "all" | "video" | "pdf" | "article"
	tagsParam := c.Query("tags") // "onboarding,tutorial" | "onboarding"
	search := c.Query("search") // "getting%20started" | "v1.2"
	cursor := c.Query("cursor") 
	limitStr := c.DefaultQuery("limit", "20")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > MaxPageSize {
		limit = DefaultPageSize
	}

	// Build Firestore query
	query := app.firestoreClient.Collection(CollectionResources).OrderBy("createdAt", firestore.Desc)

	// Apply type filter
	if typeFilter != "all" {
		query = query.Where("type", "==", typeFilter)
	}

	// Apply tags filter
	if tagsParam != "" {
		tags := normalizeTags(strings.Split(tagsParam, ","))
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
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "query_failed",
			Message: "Failed to fetch resources",
		})
		return
	}

	// Process results
	resources := make([]Resource, 0, len(docs))
	for i, doc := range docs {
		if i >= limit { // We fetched one extra to check if there are more
			break
		}

		var resource Resource
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
	response := PaginatedResponse{
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

// getResource handles GET /api/resources/:id
func (app *App) getResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	doc, err := app.firestoreClient.Collection(CollectionResources).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "resource_not_found",
			Message: "Resource not found",
		})
		return
	}

	var resource Resource
	if err := doc.DataTo(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "data_conversion_failed",
			Message: "Failed to process resource data",
		})
		return
	}
	resource.ID = doc.Ref.ID

	c.JSON(http.StatusOK, resource)
}

// createResource handles POST /api/resources
func (app *App) createResource(c *gin.Context) {
	ctx := c.Request.Context()

	// Log the content type for debugging
	log.Printf("Content-Type: %s", c.GetHeader("Content-Type"))
	log.Printf("Content-Length: %s", c.GetHeader("Content-Length"))


	// Set max multipart memory (this replaces ParseMultipartForm)
	// c.Request.ParseMultipartForm(MaxFileSize)
	

		// Check if it's actually a multipart form
		contentType := c.GetHeader("Content-Type")
		if contentType == "" || !strings.Contains(contentType, "multipart/form-data") {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "invalid_content_type",
				Message: "Request must be multipart/form-data",
			})
			return
		}

		if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
			log.Printf("ParseMultipartForm error: %v", err)
			
			// Provide more specific error messages
			var message string
			if strings.Contains(err.Error(), "too large") {
				message = fmt.Sprintf("File too large. Maximum size is %d MB", MaxFileSize/(1<<20))
			} else if strings.Contains(err.Error(), "no multipart boundary") {
				message = "Invalid multipart form data - no boundary found"
			} else {
				message = fmt.Sprintf("Failed to parse form data: %v", err)
			}
			
			c.JSON(http.StatusBadRequest, ErrorResponse{
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
	resource := Resource{
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		Type:        c.PostForm("type"),
		URL:         c.PostForm("url"),
		Tags:        normalizeTags(strings.Split(c.PostForm("tags"), ",")),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Validate required fields
	if resource.Title == "" || resource.Description == "" || resource.Type == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "missing_required_fields",
			Message: "Title, description, and type are required",
		})
		return
	}

	// Validate resource type
	if resource.Type != "video" && resource.Type != "pdf" && resource.Type != "article" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_type",
			Message: "Type must be 'video', 'pdf', or 'article'",
		})
		return
	}

	// Handle file uploads for video and pdf types
	if resource.Type == "video" || resource.Type == "pdf" {
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error:   "file_required",
				Message: fmt.Sprintf("File is required for %s resources", resource.Type),
			})
			return
		}
		defer file.Close()

		// Upload file to Cloud Storage
		url, err := app.uploadFile(ctx, file, header, resource.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
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
		thumbnailURL, err := app.uploadFile(ctx, thumbnailFile, thumbnailHeader, "image")
		if err != nil {
			log.Printf("Failed to upload thumbnail: %v", err)
		} else {
			resource.ThumbnailURL = thumbnailURL
		}
	}

	// Save to Firestore
	docRef, _, err := app.firestoreClient.Collection(CollectionResources).Add(ctx, resource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "save_failed",
			Message: "Failed to save resource",
		})
		return
	}

	// Update tag usage counts
	app.updateTagUsage(ctx, resource.Tags, 1)

	resource.ID = docRef.ID
	c.JSON(http.StatusCreated, resource)
}

// updateResource handles PUT /api/resources/:id
func (app *App) updateResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get existing resource
	doc, err := app.firestoreClient.Collection(CollectionResources).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "resource_not_found",
			Message: "Resource not found",
		})
		return
	}

	var existingResource Resource
	if err := doc.DataTo(&existingResource); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "data_conversion_failed",
			Message: "Failed to process existing resource data",
		})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(MaxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "form_parse_error",
			Message: "Failed to parse form data",
		})
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
	if tagsStr := c.PostForm("tags"); tagsStr != "" {
		oldTags := existingResource.Tags
		newTags := normalizeTags(strings.Split(tagsStr, ","))
		updatedResource.Tags = newTags
		
		// Update tag usage counts
		app.updateTagUsage(ctx, oldTags, -1)
		app.updateTagUsage(ctx, newTags, 1)
	}
	updatedResource.UpdatedAt = time.Now()

	// Handle file replacement
	if file, header, err := c.Request.FormFile("file"); err == nil {
		defer file.Close()

		// Delete old file if it exists and is stored in our bucket
		if existingResource.URL != "" && strings.Contains(existingResource.URL, app.storageBucket) {
			// app.deleteFileFromURL(ctx, existingResource.URL)
		}

		// Upload new file
		url, err := app.uploadFile(ctx, file, header, updatedResource.Type)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "upload_failed",
				Message: "Failed to upload new file",
			})
			return
		}
		updatedResource.URL = url
	}

	// Handle thumbnail replacement
	if thumbnailFile, thumbnailHeader, err := c.Request.FormFile("thumbnail"); err == nil {
		defer thumbnailFile.Close()

		// Delete old thumbnail
		if existingResource.ThumbnailURL != "" {
			// app.deleteFileFromURL(ctx, existingResource.ThumbnailURL)
		}

		// Upload new thumbnail
		thumbnailURL, err := app.uploadFile(ctx, thumbnailFile, thumbnailHeader, "image")
		if err != nil {
			log.Printf("Failed to upload thumbnail: %v", err)
		} else {
			updatedResource.ThumbnailURL = thumbnailURL
		}
	}

	// Save updated resource
	_, err = app.firestoreClient.Collection(CollectionResources).Doc(id).Set(ctx, updatedResource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "update_failed",
			Message: "Failed to update resource",
		})
		return
	}

	updatedResource.ID = id
	c.JSON(http.StatusOK, updatedResource)
}

// deleteResource handles DELETE /api/resources/:id
func (app *App) deleteResource(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	// Get existing resource to clean up files
	doc, err := app.firestoreClient.Collection(CollectionResources).Doc(id).Get(ctx)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "resource_not_found",
			Message: "Resource not found",
		})
		return
	}

	var resource Resource
	if err := doc.DataTo(&resource); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "data_conversion_failed",
			Message: "Failed to process resource data",
		})
		return
	}

	// Delete files from Cloud Storage
	if resource.URL != "" && strings.Contains(resource.URL, app.storageBucket) {
		// app.deleteFileFromURL(ctx, resource.URL)
	}
	if resource.ThumbnailURL != "" {
		// app.deleteFileFromURL(ctx, resource.ThumbnailURL)
	}

	// Update tag usage counts
	app.updateTagUsage(ctx, resource.Tags, -1)

	// Delete from Firestore
	_, err = app.firestoreClient.Collection(CollectionResources).Doc(id).Delete(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "delete_failed",
			Message: "Failed to delete resource",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Resource deleted successfully"})
}

// getTags handles GET /api/tags
func (app *App) getTags(c *gin.Context) {
	ctx := c.Request.Context()

	docs, err := app.firestoreClient.Collection(CollectionTags).OrderBy("usageCount", firestore.Desc).Documents(ctx).GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "query_failed",
			Message: "Failed to fetch tags",
		})
		return
	}

	tags := make([]Tag, 0, len(docs))
	for _, doc := range docs {
		var tag Tag
		if err := doc.DataTo(&tag); err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, tags)
}

// uploadFile uploads a file to Cloud Storage and returns the public URL
func (app *App) uploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, fileType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	baseFilename := strings.ReplaceAll(filepath.Base(header.Filename), " ", "_")
	baseFilename = strings.ReplaceAll(baseFilename, ext, "") // Remove original extension before adding our own clean one
	filename := fmt.Sprintf("%s/%d_%s%s", fileType, time.Now().Unix(), strings.ReplaceAll(header.Filename, " ", "_"), ext)

	log.Printf("bucket name: %v", app.storageBucket)

	// Create Cloud Storage object
	bucketHandle, err := app.storageClient.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("failed to get default bucket from Firebase Admin SDK: %w. Check firebase.Config.StorageBucket and emulator logs", err)
	}

	log.Printf("Obtained bucket handle for bucket: '%s'", bucketHandle.BucketName())

	obj := bucketHandle.Object(filename)

	log.Printf("obj: %v", obj)

	// Create writer
	writer := obj.NewWriter(ctx)
	writer.ContentType = header.Header.Get("Content-Type")
	// writer.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	// Copy file content
	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return "", err
	}

	if err := writer.Close(); err != nil {
		return "", err
	}

	var publicURL string
	isEmulator := os.Getenv("IS_FIREBASE_EMULATOR") == "true"

	if isEmulator {
		emulatorHost := os.Getenv("FIREBASE_STORAGE_EMULATOR_HOST") // e.g., 127.0.0.1:8082
		// For the emulator, the URL format is typically: http://{host}/v0/b/{bucket}/o/{object_path_encoded}?alt=media
		// The object name needs to be URL path encoded.
		encodedObjectName := url.PathEscape(obj.ObjectName())
		publicURL = fmt.Sprintf("http://%s/v0/b/%s/o/%s?alt=media", emulatorHost, obj.BucketName(), encodedObjectName)
		log.Printf("Generated emulator public URL: %s", publicURL)
	} else {
		// For live Firebase Storage, the typical public URL format.
		publicURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", obj.BucketName(), obj.ObjectName())
		log.Printf("Generated production public URL: %s", publicURL)
	}

	return publicURL, nil
}

// deleteFileFromURL deletes a file from Cloud Storage given its public URL
// func (app *App) deleteFileFromURL(ctx context.Context, url string) error {
// 	// Extract object name from URL
// 	parts := strings.Split(url, "/")
// 	if len(parts) < 2 {
// 		return fmt.Errorf("invalid URL format")
// 	}
	
// 	objectName := strings.Join(parts[len(parts)-2:], "/") // Get the last two parts (folder/filename)
	
// 	bucket := app.storageClient.Bucket(app.storageBucket)
// 	obj := bucket.Object(objectName)
	
// 	return obj.Delete(ctx)
// }

// updateTagUsage updates the usage count for tags
func (app *App) updateTagUsage(ctx context.Context, tags []string, delta int) {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		tagRef := app.firestoreClient.Collection(CollectionTags).Doc(tag)
		
		// Use a transaction to ensure atomicity
		err := app.firestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			doc, err := tx.Get(tagRef)
			if err != nil {
				// Tag doesn't exist, create it
				return tx.Set(tagRef, Tag{
					Name:       tag,
					UsageCount: max(0, delta),
				})
			}

			var existingTag Tag
			if err := doc.DataTo(&existingTag); err != nil {
				return err
			}

			newCount := max(0, existingTag.UsageCount+delta)
			if newCount == 0 {
				// Delete tag if usage count reaches 0
				return tx.Delete(tagRef)
			}

			return tx.Update(tagRef, []firestore.Update{
				{Path: "usageCount", Value: newCount},
			})
		})

		if err != nil {
			log.Printf("Failed to update tag usage for '%s': %v", tag, err)
		}
	}
}

// normalizeTags normalizes a slice of tags by trimming spaces and converting to lowercase
func normalizeTags(tags []string) []string {
	normalized := make([]string, 0, len(tags))
	seen := make(map[string]bool)

	for _, tag := range tags {
		tag = strings.TrimSpace(strings.ToLower(tag))
		if tag != "" && !seen[tag] {
			normalized = append(normalized, tag)
			seen[tag] = true
		}
	}

	return normalized
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}