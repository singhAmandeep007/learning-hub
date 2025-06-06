package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"

	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"

	"learning-hub/config"
	"learning-hub/constants"
	"learning-hub/firebase"
	"learning-hub/models"
)

// NormalizeTags processes a slice of tags by:
// 1. Converting all tags to lowercase
// 2. Trimming whitespace
// 3. Removing empty tags
// 4. Removing duplicates
//
// This ensures consistent tag formatting across the application.
//
// Example:
//
//	tags := []string{"  GoLang  ", "golang", "Backend", "  ", "backend"}
//	normalized := NormalizeTags(tags)
//	// Result: []string{"golang", "backend"}
//
// Parameters:
//   - tags: []string - A slice of tags to normalize
//
// Returns:
//   - []string - A new slice containing unique, normalized tags
func NormalizeTags(tags []string) []string {
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

// Max returns the maximum of two integers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// UpdateTagUsage updates the usage count for tags
func UpdateTagUsage(ctx context.Context, tags []string, delta int) {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		tagRef := firebase.FirestoreClient.Collection(constants.CollectionTags).Doc(tag)

		// Use a transaction to ensure atomicity
		err := firebase.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
			doc, err := tx.Get(tagRef)
			if err != nil {
				// Tag doesn't exist, create it
				return tx.Set(tagRef, models.Tag{
					Name:       tag,
					UsageCount: max(0, delta),
				})
			}

			var existingTag models.Tag
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

// FileUploadResult contains the result of a file upload operation
type FileUploadResult struct {
	PublicURL string
	Filename  string
	Size      int64
}

func UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, fileType string) (*FileUploadResult, error) {
	// Generate clean, unique filename
	filename, err := generateUniqueFilename(header.Filename, fileType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filename: %w", err)
	}

	log.Printf("Uploading file: %s to bucket: %s", filename, firebase.StorageBucket)

	bucketHandler := firebase.StorageClient.Bucket(firebase.StorageBucket)

	writer := bucketHandler.Object(filename).NewWriter(ctx)

	// Set content type
	if contentType := header.Header.Get("Content-Type"); contentType != "" {
		writer.ContentType = contentType
		// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Disposition
		writer.ContentDisposition = "inline"
	}

	bytesWritten, err := io.Copy(writer, file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize upload: %w", err)
	}

	// Make the file public
	// Note: For granular access control, Firebase Security Rules are preferred.
	// This makes the object publicly readable.
	isEmulator := config.AppConfig.IS_FIREBASE_EMULATOR
	if !isEmulator {
		acl := bucketHandler.Object(filename).ACL()
		if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			log.Printf("Warning: Failed to set public ACL: %v (File uploaded but may not be publicly accessible)", err)
		}
	}

	// Generate public URL
	publicURL, err := generatePublicURL(filename, firebase.StorageBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public URL: %w", err)
	}

	log.Printf("File uploaded successfully: %s (%d bytes)", publicURL, bytesWritten)

	return &FileUploadResult{
		PublicURL: publicURL,
		Filename:  filename,
		Size:      bytesWritten,
	}, nil
}

// generateUniqueFilename creates a unique filename with proper sanitization
func generateUniqueFilename(originalFilename, fileType string) (string, error) {
	if originalFilename == "" {
		return "", fmt.Errorf("original filename cannot be empty")
	}

	ext := filepath.Ext(originalFilename)
	baseName := strings.TrimSuffix(filepath.Base(originalFilename), ext)

	// Sanitize filename
	baseName = strings.ReplaceAll(baseName, " ", "_")
	baseName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, baseName)

	// Generate unique filename with timestamp
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s/%d_%s%s", fileType, timestamp, baseName, ext)

	return filename, nil
}

// generatePublicURL creates the appropriate public URL based on environment
func generatePublicURL(objectName, bucketName string) (string, error) {
	isEmulator := config.AppConfig.IS_FIREBASE_EMULATOR

	if isEmulator {
		emulatorHost := config.AppConfig.FIREBASE_STORAGE_EMULATOR_HOST
		if emulatorHost == "" {
			return "", fmt.Errorf("FIREBASE_STORAGE_EMULATOR_HOST not set for emulator mode")
		}

		encodedObjectName := url.PathEscape(objectName)
		// Eg. http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748580692_image1.png?alt=media
		publicURL := fmt.Sprintf("http://%s/v0/b/%s/o/%s?alt=media", emulatorHost, bucketName, encodedObjectName)

		log.Printf("Generated emulator URL: %s", publicURL)
		return publicURL, nil
	}

	// Production URL
	publicURL := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)
	log.Printf("Generated production URL: %s", publicURL)
	return publicURL, nil
}

// DeleteFileFromURL deletes a file from Cloud Storage given its public URL
func DeleteFileFromURL(ctx context.Context, fileUrl string) error {
	// Delete file if it is stored in our bucket
	if strings.Contains(fileUrl, firebase.StorageBucket) {
		bucketName, objectName, err := parseStorageURL(fileUrl)
		log.Printf("bucketName: %s objectName: %s", bucketName, objectName)

		if err != nil {
			return fmt.Errorf("failed to parse storage URL: %w", err)
		}

		// Get the bucket handle
		bucketHandler := firebase.StorageClient.Bucket(bucketName)

		// Get the object handle
		objHandler := bucketHandler.Object(objectName)

		// Delete the object
		if err := objHandler.Delete(ctx); err != nil {
			return fmt.Errorf("failed to delete object %s from bucket %s: %w", objectName, bucketName, err)
		}

		return nil
	}

	return nil
}

func parseStorageURL(fileUrl string) (bucketName, objectName string, err error) {
	parsedURL, err := url.Parse(fileUrl)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL format: %w", err)
	}

	isEmulator := config.AppConfig.IS_FIREBASE_EMULATOR

	if isEmulator {
		// http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748580692_image1.png?alt=media
		pathRegex := regexp.MustCompile(`^/v0/b/([^/]+)/o/(.+)$`)
		matches := pathRegex.FindStringSubmatch(parsedURL.Path)

		if len(matches) != 3 {
			return "", "", fmt.Errorf("invalid Firebase Storage URL path format")
		}

		bucketName = matches[1]
		encodedObjectName := matches[2]

		// URL decode the object name
		objectName, err = url.QueryUnescape(encodedObjectName)
		if err != nil {
			return "", "", fmt.Errorf("failed to decode object name: %w", err)
		}

		return bucketName, objectName, nil
	} else {
		pathParts := strings.SplitN(strings.TrimPrefix(parsedURL.Path, "/"), "/", 2)
		if len(pathParts) < 2 {
			return "", "", fmt.Errorf("invalid Google Cloud Storage URL path format")
		}

		bucketName = pathParts[0]
		objectName = pathParts[1]

		return bucketName, objectName, nil
	}
}
