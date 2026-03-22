package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"learninghub/config"
	"learninghub/constants"
	"learninghub/db"
	"learninghub/firebase"
	"learninghub/pkg/logger"

	"cloud.google.com/go/storage"
	"github.com/gabriel-vasile/mimetype"
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
//	Result: []string{"golang", "backend"}
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

// UpdateTagUsage updates the usage count for tags in a product-specific collection
func UpdateTagUsage(ctx context.Context, product string, tags []string, delta int) {
	database := db.New()
	tagService := db.NewTagService(database)

	if err := tagService.UpdateUsage(ctx, product, tags, delta); err != nil {
		logger.Infof("Failed to update tag usage: %v", err)
	}
}

// FileUploadResult contains the result of a file upload operation
type FileUploadResult struct {
	PublicURL string
	Filename  string
	Size      int64
}

// UploadFile uploads a file to Firebase Cloud Storage and returns the public URL
func UploadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader, product, fileType string) (*FileUploadResult, error) {
	// SECURITY: Validate file content using magic bytes detection.
	// This prevents attackers from uploading malicious files by spoofing the
	// Content-Type header. The actual file bytes are inspected, not the header.
	validationResult := ValidateFileContent(file, fileType)
	if !validationResult.IsValid {
		return nil, fmt.Errorf("%s: %s", constants.ErrFileValidationFailed, validationResult.Error)
	}

	// SECURITY: Use the detected extension from magic bytes analysis, NOT the
	// extension from the original filename. This prevents an attacker from uploading
	// a valid JPEG (passes magic bytes) but naming it "exploit.html" so the stored
	// object has a dangerous extension.
	filename, err := generateUniqueFilename(header.Filename, product, fileType, validationResult.Extension)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filename: %w", err)
	}

	bucketHandler := firebase.StorageClient.Bucket(firebase.StorageBucket)
	writer := bucketHandler.Object(filename).NewWriter(ctx)

	// SECURITY: Use the detected MIME type from file content, not the client-provided
	// header. This ensures the Content-Type stored in GCS (and served to browsers)
	// reflects the actual file content.
	writer.ContentType = validationResult.DetectedMIME

	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Content-Disposition
	writer.ContentDisposition = "inline"

	bytesWritten, err := io.Copy(writer, file)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Close the writer to finalize the upload
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to finalize upload: %w", err)
	}

	// Generate public URL
	publicURL, err := generatePublicURL(filename, firebase.StorageBucket)
	if err != nil {
		return nil, fmt.Errorf("failed to generate public URL: %w", err)
	}

	return &FileUploadResult{
		PublicURL: publicURL,
		Filename:  filename,
		Size:      bytesWritten,
	}, nil
}

// generateUniqueFilename creates a unique filename with proper sanitization.
//
// SECURITY: The detectedExtension parameter MUST come from magic bytes analysis
// (validationResult.Extension), never from the original filename. This prevents
// an attacker from storing an object with a dangerous extension (e.g. .html, .php)
// by embedding it in their uploaded filename while the actual bytes pass validation.
func generateUniqueFilename(originalFilename, product, fileType, detectedExtension string) (string, error) {
	if originalFilename == "" {
		return "", fmt.Errorf("original filename cannot be empty")
	}

	// Strip the original extension and use the sanitized base name only.
	// The final extension comes from magic bytes detection, not from here.
	baseName := strings.TrimSuffix(filepath.Base(originalFilename), filepath.Ext(originalFilename))

	// Sanitize: replace spaces and any non-alphanumeric/underscore/dash characters
	baseName = strings.ReplaceAll(baseName, " ", "_")
	baseName = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == '-' {
			return r
		}
		return '_'
	}, baseName)

	// Generate unique filename with timestamp, using the DETECTED extension
	timestamp := time.Now().UnixNano()
	filename := fmt.Sprintf("%s/%s/%d_%s%s", product, fileType, timestamp, baseName, detectedExtension)

	return filename, nil
}

// generatePublicURL creates the appropriate public URL based on environment
func generatePublicURL(objectName, bucketName string) (string, error) {
	isDev := config.AppConfig.ENV_MODE == constants.EnvModeDev

	if isDev {
		emulatorHost := config.AppConfig.FIREBASE_STORAGE_EMULATOR_HOST
		if emulatorHost == "" {
			return "", fmt.Errorf("FIREBASE_STORAGE_EMULATOR_HOST not set for emulator mode")
		}

		emulatorHostPort := strings.Split(emulatorHost, ":")[1]

		encodedObjectName := url.PathEscape(objectName)
		// Eg. http://127.0.0.1:8082/v0/b/learninghub-81cc6.firebasestorage.app/o/image%2F1748580692_image1.png?alt=media
		publicURL := fmt.Sprintf("http://127.0.0.1:%s/v0/b/%s/o/%s?alt=media", emulatorHostPort, bucketName, encodedObjectName)

		return publicURL, nil
	}

	// Production URL
	encodedObjectName := url.PathEscape(objectName)
	// Eg. https://firebasestorage.googleapis.com/v0/b/qa-us-firestore.firebasestorage.app/o/ecomm%2Fimage%2F1759303167962121049_final_step.png?alt=media
	publicURL := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media", bucketName, encodedObjectName)

	return publicURL, nil
}

// GenerateSignedURL generates a signed URL for a storage object
// This allows temporary authenticated access to private objects
func GenerateSignedURL(ctx context.Context, storageURL string, expirationMinutes int) (string, error) {
	isDev := config.AppConfig.ENV_MODE == constants.EnvModeDev

	// In dev mode, emulator URLs are already accessible, return as-is
	if isDev {
		return storageURL, nil
	}

	// Only generate signed URLs for our storage bucket URLs
	if !IsValidStorageURL(storageURL) {
		// External URLs (like article links) - return as-is
		return storageURL, nil
	}

	// Parse the storage URL to get bucket and object name
	bucketName, objectName, err := parseStorageURL(storageURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse storage URL: %w", err)
	}

	// Get the bucket and object handles
	bucket := firebase.StorageClient.Bucket(bucketName)

	// Generate signed URL with expiration
	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(time.Duration(expirationMinutes) * time.Minute),
	}

	signedURL, err := bucket.SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}

// ConvertResourceURLsToSigned converts a Resource's URLs to signed URLs if needed.
// This is a helper to transform URLs before sending to frontend.
func ConvertResourceURLsToSigned(ctx context.Context, url, thumbnailURL string, expirationMinutes int) (signedURL, signedThumbnailURL string, err error) {
	if url != "" {
		signedURL, err = GenerateSignedURL(ctx, url, expirationMinutes)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate signed URL for resource: %w", err)
		}
	}

	if thumbnailURL != "" {
		signedThumbnailURL, err = GenerateSignedURL(ctx, thumbnailURL, expirationMinutes)
		if err != nil {
			return "", "", fmt.Errorf("failed to generate signed URL for thumbnail: %w", err)
		}
	}

	return signedURL, signedThumbnailURL, nil
}

// DeleteFileFromURL deletes a file from Cloud Storage given its public URL
func DeleteFileFromURL(ctx context.Context, fileURL string) error {
	// Delete file if it is stored in our bucket
	if IsValidStorageURL(fileURL) {
		bucketName, objectName, err := parseStorageURL(fileURL)
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

func parseStorageURL(fileURL string) (bucketName, objectName string, err error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", "", fmt.Errorf("invalid URL format: %w", err)
	}

	// /v0/b/{bucket}/o/{object}
	// http://127.0.0.1:8082/v0/b/learninghub-81cc6.firebasestorage.app/o/product/image%2F1748580692_image1.png?alt=media
	// https://firebasestorage.googleapis.com/v0/b/qa-us-firestore.firebasestorage.app/o/ecomm%2Fimage%2F1759318704892973803_Dialog_modal.png?alt=media
	pathRegex := regexp.MustCompile(`^/v0/b/([^/]+)/o/(.+)$`)
	matches := pathRegex.FindStringSubmatch(parsedURL.Path)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("invalid Firebase Storage URL path format: %s", parsedURL.Path)
	}

	bucketName = matches[1]
	encodedObjectName := matches[2]

	objectName, err = url.QueryUnescape(encodedObjectName)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode object name: %w", err)
	}

	// Remove query parameters if they got included (e.g., ?alt=media)
	if idx := strings.Index(objectName, "?"); idx != -1 {
		objectName = objectName[:idx]
	}

	return bucketName, objectName, nil
}

// Validations

// IsValidResourceType check if resource type is valid
func IsValidResourceType(t string) bool {
	for _, v := range constants.ResourceTypes {
		if t == v {
			return true
		}
	}
	return false
}

// IsValidStorageURL checks if url points to resource stored in storage
func IsValidStorageURL(fileURL string) bool {
	return strings.Contains(fileURL, firebase.StorageBucket)
}

// IsValidProduct checks if product name is valid
func IsValidProduct(product string) bool {
	for _, v := range config.AppConfig.VALID_PRODUCTS {
		if product == v {
			return true
		}
	}
	return false
}

// File content validation

// FileValidationResult contains the result of file content validation
type FileValidationResult struct {
	IsValid      bool
	DetectedMIME string
	Extension    string
	Error        string
}

// validationError creates a failed FileValidationResult with the given error message
func validationError(detectedMIME, errMsg string) *FileValidationResult {
	return &FileValidationResult{
		IsValid:      false,
		DetectedMIME: detectedMIME,
		Error:        errMsg,
	}
}

// validationSuccess creates a successful FileValidationResult
func validationSuccess(detectedMIME, extension string) *FileValidationResult {
	return &FileValidationResult{
		IsValid:      true,
		DetectedMIME: detectedMIME,
		Extension:    extension,
	}
}

// blockedMIMETypes contains MIME types that are explicitly blocked for security reasons.
// These types can contain executable code or be used for XSS attacks.
var blockedMIMETypes = []string{
	"text/html",
	"application/xhtml+xml",
	"image/svg+xml", // SVG can contain embedded <script> tags
	"application/javascript",
	"text/javascript",
	"application/x-httpd-php",
	"application/x-sh",
	"application/x-shellscript",
	"application/x-msdownload", // Windows PE executables (.exe, .dll)
	"application/x-executable", // Linux ELF executables
}

// allowedVideoMIMEs is an explicit allowlist of video MIME types the application
// accepts. Using an allowlist (rather than strings.HasPrefix("video/")) ensures
// that unusual or dangerous "video/*" sub-types are rejected by default.
var allowedVideoMIMEs = map[string]bool{
	"video/mp4":        true,
	"video/webm":       true,
	"video/x-matroska": true, // WebM is a subset of Matroska; some detectors report this MIME
}

// pdfSuspiciousPatterns is a comprehensive list of PDF dictionary keys and
// JavaScript patterns associated with active content, auto-actions, and known
// PDF exploit techniques.
//
// Categories:
//   - JavaScript execution: /JavaScript, /JS
//   - Automatic trigger actions: /AA (Additional Actions), /OpenAction
//   - External interaction: /Launch, /SubmitForm, /ImportData, /URI, /GoToR, /GoToE
//   - Rich media / multimedia: /RichMedia, /Sound, /Movie, /Rendition, /3D, /GoTo3DView
//   - Embedded content: /EmbeddedFile, /FileAttachment
//   - XFA forms (can execute JS): /XFA
//   - Obfuscation / hiding: /Hide, /SetOCGState
//   - Navigation abuse: /Named
//   - Common JS API calls found in PDF exploits: app.alert, this.exportDataObject,
//     util.printf, app.launchURL, this.submitForm, getAnnots, getField,
//     app.setTimeOut (used for delayed execution), event.value
//   - Malformed stream markers used in heap-spray exploits: %eval, %execute
var pdfSuspiciousPatterns = [][]byte{
	// ---- JavaScript execution ----
	[]byte("/JavaScript"),
	[]byte("/JS"),

	// ---- Automatic trigger actions ----
	[]byte("/AA"),         // Additional Actions — runs JS on field/page events
	[]byte("/OpenAction"), // Runs an action when the PDF is opened

	// ---- External interaction ----
	[]byte("/Launch"),     // Launches an external application
	[]byte("/SubmitForm"), // Submits form data to a remote URL
	[]byte("/ImportData"), // Imports FDF data from an external source
	[]byte("/URI"),        // URI action — can silently load external URLs
	[]byte("/GoToR"),      // GoTo remote — opens another PDF or file
	[]byte("/GoToE"),      // GoTo embedded — navigates into an embedded file

	// ---- Rich media / multimedia ----
	[]byte("/RichMedia"),  // Embeds Flash/video — historic attack surface
	[]byte("/Sound"),      // Sound action
	[]byte("/Movie"),      // Movie action
	[]byte("/Rendition"),  // Media rendition action
	[]byte("/3D"),         // 3D annotation (U3D/PRC — has a JS API)
	[]byte("/GoTo3DView"), // Navigates to a named 3D view

	// ---- Embedded content ----
	[]byte("/EmbeddedFile"),   // Attached/embedded files
	[]byte("/FileAttachment"), // File attachment annotation

	// ---- XFA forms ----
	[]byte("/XFA"), // XML Forms Architecture — executes JS during render

	// ---- Obfuscation / hiding ----
	[]byte("/Hide"),        // Hides/shows annotations — used in social engineering
	[]byte("/SetOCGState"), // Controls Optional Content Group visibility

	// ---- Navigation abuse ----
	[]byte("/Named"), // Named action (e.g. NextPage) — can be chained maliciously

	// ---- Common JS API calls seen in PDF exploits ----
	[]byte("app.alert"),             // Classic PoC and phishing lure
	[]byte("app.launchURL"),         // Opens external URL silently
	[]byte("app.setTimeOut"),        // Delayed execution — evasion technique
	[]byte("this.exportDataObject"), // Data exfiltration via embedded files
	[]byte("this.submitForm"),       // Exfiltrates data to attacker server
	[]byte("util.printf"),           // Used in format-string exploits (CVE-2008-2992)
	[]byte("getAnnots"),             // Annotation enumeration — recon step
	[]byte("getField"),              // Form field access — used in JS-based attacks
	[]byte("event.value"),           // Field event handler — common JS entry point

	// ---- Heap-spray / exploit markers ----
	[]byte("%eval"),    // Obfuscated eval patterns in malformed PDFs
	[]byte("%execute"), // Obfuscated execute patterns in malformed PDFs
}

// scanPDFForSuspiciousContent reads the entire PDF into memory and checks for
// known dangerous patterns (embedded JavaScript, auto-actions, exploit markers).
//
// The file seek position is reset to the start after scanning so the caller can
// still read the file normally (e.g. for upload).
//
// Note: This is a heuristic scan — it catches the overwhelming majority of
// real-world malicious PDFs. A determined attacker can obfuscate content beyond
// what byte-scanning detects; for higher-assurance environments consider a
// sandboxed PDF renderer.
//
// Parameters:
//   - file: multipart.File — the uploaded file, already validated as PDF by magic bytes
//
// Returns:
//   - bool   — true if a suspicious pattern was found
//   - []byte — the matched pattern (for logging/debugging)
//   - error  — non-nil if reading or seeking failed
func scanPDFForSuspiciousContent(file multipart.File) (bool, []byte, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return false, nil, fmt.Errorf("failed to read PDF for scanning: %w", err)
	}

	// Reset so the caller (UploadFile) can still read the complete file
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return false, nil, fmt.Errorf("failed to reset file position after PDF scan: %w", err)
	}

	for _, pattern := range pdfSuspiciousPatterns {
		if bytes.Contains(content, pattern) {
			return true, pattern, nil
		}
	}

	return false, nil, nil
}

// ValidateFileContent validates file content using magic bytes detection.
// It uses the gabriel-vasile/mimetype package for accurate MIME type detection
// based on file signatures (magic numbers), then applies additional type-specific
// checks (e.g. embedded JS scanning for PDFs).
//
// The file seek position is always reset to the beginning before returning so
// the caller can subsequently read the full file for upload.
//
// Parameters:
//   - file:         multipart.File — the uploaded file to validate
//   - expectedType: string         — the expected resource type (video, pdf, image)
//
// Returns:
//   - *FileValidationResult — contains IsValid, DetectedMIME, Extension, and Error
func ValidateFileContent(file multipart.File, expectedType string) *FileValidationResult {
	// --- Step 1: Magic bytes detection ---
	// mimetype.DetectReader reads file signatures (magic numbers) to determine
	// the true MIME type, completely ignoring the client-supplied Content-Type.
	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return validationError("", fmt.Sprintf("failed to detect file type: %v", err))
	}

	// Reset file position to the beginning for all subsequent reads
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return validationError("", fmt.Sprintf("failed to reset file position: %v", err))
	}

	detectedMIME := mtype.String()

	// --- Step 2: Block explicitly dangerous MIME types ---
	// Even if a type would pass the expected-type check below, types in this list
	// are unconditionally rejected because they can carry executable content.
	for _, blocked := range blockedMIMETypes {
		if mtype.Is(blocked) {
			return validationError(detectedMIME, fmt.Sprintf("file type '%s' is not allowed for security reasons", detectedMIME))
		}
	}

	// --- Step 3: Type-specific validation ---
	var isValid bool
	switch expectedType {

	case constants.ResourceTypeVideo:
		// Use an explicit allowlist — "video/*" is too broad and would accept
		// obscure or potentially dangerous video sub-types.
		isValid = allowedVideoMIMEs[detectedMIME]
		if !isValid {
			return validationError(detectedMIME, fmt.Sprintf(
				"video type '%s' is not supported; accepted types: mp4, webm", detectedMIME,
			))
		}

	case constants.ResourceTypePDF:
		// First confirm magic bytes identify this as a real PDF
		isValid = mtype.Is("application/pdf")
		if !isValid {
			return validationError(detectedMIME, fmt.Sprintf(
				"file content type '%s' does not match expected type '%s'", detectedMIME, expectedType,
			))
		}

		// --- Step 4 (PDF only): Scan for embedded JavaScript and exploit patterns ---
		// Magic bytes only confirm the file IS a PDF. They say nothing about what
		// is inside it. PDFs can contain JavaScript (/JS, /JavaScript), automatic
		// open-actions (/OpenAction), XFA forms, and many other active-content
		// features that can be weaponised for RCE or data exfiltration.
		found, matchedPattern, err := scanPDFForSuspiciousContent(file)
		if err != nil {
			return validationError(detectedMIME, fmt.Sprintf("PDF content scan failed: %v", err))
		}
		if found {
			// Log the matched pattern (helpful for incident response) but do not
			// expose the raw pattern string to the end user.
			logger.Infof("PDF upload rejected: suspicious pattern detected: %q", matchedPattern)
			return validationError(detectedMIME, "PDF contains potentially malicious embedded content and cannot be uploaded")
		}

	case constants.ResourceTypeImage:
		// Accept any image/* type but exclude SVG (already in blockedMIMETypes;
		// this is a second, explicit guard).
		isValid = strings.HasPrefix(detectedMIME, "image/") && !mtype.Is("image/svg+xml")
		if !isValid {
			return validationError(detectedMIME, fmt.Sprintf(
				"file type '%s' is not a supported image format", detectedMIME,
			))
		}

	default:
		return validationError(detectedMIME, fmt.Sprintf("unknown resource type: %s", expectedType))
	}

	return validationSuccess(detectedMIME, mtype.Extension())
}
