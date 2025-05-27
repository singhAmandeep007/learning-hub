package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/firestore"

	"learning-hub/config"
	"learning-hub/constants"
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
//   tags := []string{"  GoLang  ", "golang", "Backend", "  ", "backend"}
//   normalized := NormalizeTags(tags)
//   // Result: []string{"golang", "backend"}
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

		tagRef := config.FirestoreClient.Collection(constants.CollectionTags).Doc(tag)
		
		// Use a transaction to ensure atomicity
		err := config.FirestoreClient.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
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

// UploadFile uploads a file to Cloud Storage and returns the public URL
func UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, fileType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	baseFilename := strings.ReplaceAll(filepath.Base(header.Filename), " ", "_")
	baseFilename = strings.ReplaceAll(baseFilename, ext, "") // Remove original extension before adding our own clean one
	filename := fmt.Sprintf("%s/%d_%s%s", fileType, time.Now().Unix(), strings.ReplaceAll(header.Filename, " ", "_"), ext)

	log.Printf("bucket name: %v", config.StorageBucket)

	// Create Cloud Storage object
	bucketHandle, err := config.StorageClient.DefaultBucket()
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