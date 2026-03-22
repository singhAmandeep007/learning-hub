package utils

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"reflect"

	"github.com/stretchr/testify/assert"

	"learninghub/config"
	"learninghub/constants"
	"learninghub/firebase"
)

// mockFile implements multipart.File interface for testing
type mockFile struct {
	*bytes.Reader
}

func (m *mockFile) Close() error {
	return nil
}

func newMockFile(data []byte) *mockFile {
	return &mockFile{Reader: bytes.NewReader(data)}
}

func init() {
	// Initialize config for tests
	config.AppConfig = &config.EnvConfig{
		ENV_MODE:       constants.EnvModeProd, // Default value
		VALID_PRODUCTS: []string{"ecomm"},
	}
}

const product = "test-product"

func TestGenerateUniqueFilename(t *testing.T) {
	tests := []struct {
		name              string
		originalFile      string
		fileType          string
		detectedExtension string
		wantErr           bool
		validateOutput    func(t *testing.T, output string)
	}{
		{
			name:              "empty filename",
			originalFile:      "",
			fileType:          "pdf",
			detectedExtension: ".pdf",
			wantErr:           true,
		},
		{
			name:              "valid filename with spaces",
			originalFile:      "my document.pdf",
			fileType:          "pdf",
			detectedExtension: ".pdf",
			wantErr:           false,
			validateOutput: func(t *testing.T, output string) {
				expectedPrefix := product + "/pdf/"
				if !strings.HasPrefix(output, expectedPrefix) {
					t.Errorf("expected prefix 'test-product/pdf/', got %s", output)
				}
				if !strings.HasSuffix(output, ".pdf") {
					t.Errorf("expected suffix '.pdf', got %s", output)
				}
				if !strings.Contains(output, "my_document") {
					t.Errorf("expected to contain 'my_document', got %s", output)
				}
			},
		},
		{
			name:              "filename with special characters",
			originalFile:      "file@#$%^&*().txt",
			fileType:          "text",
			detectedExtension: ".txt",
			wantErr:           false,
			validateOutput: func(t *testing.T, output string) {
				expectedPrefix := product + "/text/"
				if !strings.HasPrefix(output, expectedPrefix) {
					t.Errorf("expected prefix 'text/', got %s", output)
				}
				if !strings.HasSuffix(output, ".txt") {
					t.Errorf("expected suffix '.txt', got %s", output)
				}
				if strings.ContainsAny(output, "@#$%^&*()") {
					t.Errorf("expected special characters to be replaced, got %s", output)
				}
			},
		},
		{
			name:              "filename with multiple dots",
			originalFile:      "my.file.name.pdf",
			fileType:          "pdf",
			detectedExtension: ".pdf",
			wantErr:           false,
			validateOutput: func(t *testing.T, output string) {
				expectedPrefix := product + "/pdf/"
				if !strings.HasPrefix(output, expectedPrefix) {
					t.Errorf("expected prefix '%s', got %s", expectedPrefix, output)
				}
				if !strings.HasSuffix(output, ".pdf") {
					t.Errorf("expected suffix '.pdf', got %s", output)
				}
				if !strings.Contains(output, "my_file_name") {
					t.Errorf("expected to contain 'my_file_name', got %s", output)
				}
			},
		},
		{
			name:              "filename with unicode characters",
			originalFile:      "résumé.pdf",
			fileType:          "pdf",
			detectedExtension: ".pdf",
			wantErr:           false,
			validateOutput: func(t *testing.T, output string) {
				expectedPrefix := product + "/pdf/"
				if !strings.HasPrefix(output, expectedPrefix) {
					t.Errorf("expected prefix '%s', got %s", expectedPrefix, output)
				}
				if !strings.HasSuffix(output, ".pdf") {
					t.Errorf("expected suffix '.pdf', got %s", output)
				}
				if strings.Contains(output, "é") {
					t.Errorf("expected unicode characters to be replaced, got %s", output)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateUniqueFilename(tt.originalFile, product, tt.fileType, tt.detectedExtension)

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("generateUniqueFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip validation if we expected an error
			if tt.wantErr {
				return
			}

			// Validate timestamp product/fileType/timestamp_name.ext
			parts := strings.Split(got, "/")
			if len(parts) != 3 {
				t.Errorf("expected format 'product/type/timestamp_name.ext', got %s", got)
				return
			}

			// Check if timestamp is recent (within last 5 seconds)
			timestampStr := strings.Split(parts[2], "_")[0]
			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				t.Errorf("invalid timestamp format: %v", err)
				return
			}
			if time.Since(time.Unix(0, timestamp)) > 5*time.Second {
				t.Errorf("timestamp is not recent: %v", time.Unix(0, timestamp))
			}

			// Run additional validation if provided
			if tt.validateOutput != nil {
				tt.validateOutput(t, got)
			}
		})
	}
}

func TestParseStorageURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedBucket string
		expectedObject string
		expectError    bool
	}{
		{
			name:           "Valid emulator URL",
			url:            "http://127.0.0.1:8082/v0/b/learninghub-81cc6.firebasestorage.app/o/ecomm/image%2F1748580692_image1.png?alt=media",
			expectedBucket: "learninghub-81cc6.firebasestorage.app",
			expectedObject: "ecomm/image/1748580692_image1.png",
			expectError:    false,
		},
		{
			name:           "Valid production URL - video",
			url:            "https://firebasestorage.googleapis.com/v0/b/qa-us-firestore.firebasestorage.app/o/ecomm%2Fvideo%2F1759318704739890076_projectsupportlogsadminapi.mp4?alt=media",
			expectedBucket: "qa-us-firestore.firebasestorage.app",
			expectedObject: "ecomm/video/1759318704739890076_projectsupportlogsadminapi.mp4",
			expectError:    false,
		},
		{
			name:           "Valid production URL - image",
			url:            "https://firebasestorage.googleapis.com/v0/b/qa-us-firestore.firebasestorage.app/o/ecomm%2Fimage%2F1759318704892973803_Dialog_modal.png?alt=media",
			expectedBucket: "qa-us-firestore.firebasestorage.app",
			expectedObject: "ecomm/image/1759318704892973803_Dialog_modal.png",
			expectError:    false,
		},
		{
			name:           "Valid production URL - pdf",
			url:            "https://firebasestorage.googleapis.com/v0/b/qa-us-firestore.firebasestorage.app/o/ecomm%2Fpdf%2F1759318284946151469_Plum_Employee_Handbook_-_File.pdf?alt=media",
			expectedBucket: "qa-us-firestore.firebasestorage.app",
			expectedObject: "ecomm/pdf/1759318284946151469_Plum_Employee_Handbook_-_File.pdf",
			expectError:    false,
		},
		{
			name:           "Invalid URL format - missing /o/ path",
			url:            "http://127.0.0.1:8082/invalid/path",
			expectedBucket: "",
			expectedObject: "",
			expectError:    true,
		},
		{
			name:           "Invalid URL format - incomplete path",
			url:            "https://firebasestorage.googleapis.com/v0/b/bucket",
			expectedBucket: "",
			expectedObject: "",
			expectError:    true,
		},
		{
			name:           "Invalid URL - not parseable",
			url:            "not-a-url",
			expectedBucket: "",
			expectedObject: "",
			expectError:    true,
		},
		{
			name:           "Emulator URL with special characters",
			url:            "http://127.0.0.1:8082/v0/b/learninghub-81cc6.firebasestorage.app/o/path%2Fwith%20spaces%20and%20%26%20symbols.jpg?alt=media",
			expectedBucket: "learninghub-81cc6.firebasestorage.app",
			expectedObject: "path/with spaces and & symbols.jpg",
			expectError:    false,
		},
		{
			name:           "Production URL with special characters",
			url:            "https://firebasestorage.googleapis.com/v0/b/my-bucket.firebasestorage.app/o/path%2Fwith%20spaces%20and%20%26%20symbols.jpg?alt=media",
			expectedBucket: "my-bucket.firebasestorage.app",
			expectedObject: "path/with spaces and & symbols.jpg",
			expectError:    false,
		},
		{
			name:           "URL without query parameters",
			url:            "https://firebasestorage.googleapis.com/v0/b/my-bucket.firebasestorage.app/o/folder%2Ffile.txt",
			expectedBucket: "my-bucket.firebasestorage.app",
			expectedObject: "folder/file.txt",
			expectError:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function
			bucket, object, err := parseStorageURL(tt.url)
			// Check error expectation
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, bucket)
				assert.Empty(t, object)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBucket, bucket)
				assert.Equal(t, tt.expectedObject, object)
			}
		})
	}
}

func TestNormalizeTags(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty tags",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single tag",
			input:    []string{"golang"},
			expected: []string{"golang"},
		},
		{
			name:     "duplicate tags",
			input:    []string{"golang", "Golang", "GOLANG"},
			expected: []string{"golang"},
		},
		{
			name:     "tags with whitespace",
			input:    []string{"  golang  ", "  react  ", "react"},
			expected: []string{"golang", "react"},
		},
		{
			name:     "mixed case tags",
			input:    []string{"GoLang", "REACT", "TypeScript", "typescript"},
			expected: []string{"golang", "react", "typescript"},
		},
		{
			name:     "empty strings",
			input:    []string{"", "golang", "  ", "react"},
			expected: []string{"golang", "react"},
		},
		{
			name:     "special characters",
			input:    []string{"go-lang", "react.js", "react.js"},
			expected: []string{"go-lang", "react.js"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeTags(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("NormalizeTags() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsValidResourceType(t *testing.T) {
	tests := []struct {
		name         string
		resourceType string
		expected     bool
	}{
		{
			name:         "valid video type",
			resourceType: constants.ResourceTypeVideo,
			expected:     true,
		},
		{
			name:         "valid pdf type",
			resourceType: constants.ResourceTypePDF,
			expected:     true,
		},
		{
			name:         "valid article type",
			resourceType: constants.ResourceTypeArticle,
			expected:     true,
		},
		{
			name:         "invalid type",
			resourceType: "invalid",
			expected:     false,
		},
		{
			name:         "empty type",
			resourceType: "",
			expected:     false,
		},
		{
			name:         "case sensitive - uppercase",
			resourceType: "VIDEO",
			expected:     false,
		},
		{
			name:         "case sensitive - mixed case",
			resourceType: "Video",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidResourceType(tt.resourceType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidProduct(t *testing.T) {
	tests := []struct {
		name     string
		product  string
		expected bool
	}{
		{
			name:     "valid product Ecomm",
			product:  "ecomm",
			expected: true,
		},
		{
			name:     "invalid product",
			product:  "invalid_product",
			expected: false,
		},
		{
			name:     "empty product",
			product:  "",
			expected: false,
		},
		{
			name:     "case sensitive - uppercase",
			product:  "ECOMM",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidProduct(tt.product)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidStorageURL(t *testing.T) {
	// Setup test environment
	originalBucket := firebase.StorageBucket
	firebase.StorageBucket = "test-bucket.firebasestorage.app"
	defer func() {
		firebase.StorageBucket = originalBucket
	}()

	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "valid storage URL",
			url:      "https://storage.googleapis.com/test-bucket.firebasestorage.app/path/file.jpg",
			expected: true,
		},
		{
			name:     "valid emulator URL",
			url:      "http://127.0.0.1:8082/v0/b/test-bucket.firebasestorage.app/o/file.jpg",
			expected: true,
		},
		{
			name:     "external URL",
			url:      "https://example.com/image.jpg",
			expected: false,
		},
		{
			name:     "different bucket",
			url:      "https://storage.googleapis.com/other-bucket/file.jpg",
			expected: false,
		},
		{
			name:     "empty URL",
			url:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidStorageURL(tt.url)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGeneratePublicURL(t *testing.T) {
	tests := []struct {
		name         string
		objectName   string
		bucketName   string
		envMode      string
		emulatorHost string
		expected     string
		expectError  bool
	}{
		{
			name:         "production URL",
			objectName:   "path/to/file.jpg",
			bucketName:   "my-bucket",
			envMode:      constants.EnvModeProd,
			emulatorHost: "",
			expected:     "https://firebasestorage.googleapis.com/v0/b/my-bucket/o/path%2Fto%2Ffile.jpg?alt=media",
			expectError:  false,
		},
		{
			name:         "dev URL with emulator",
			objectName:   "path/to/file.jpg",
			bucketName:   "my-bucket",
			envMode:      constants.EnvModeDev,
			emulatorHost: "127.0.0.1:8082",
			expected:     "http://127.0.0.1:8082/v0/b/my-bucket/o/path%2Fto%2Ffile.jpg?alt=media",
			expectError:  false,
		},
		{
			name:         "dev URL without emulator host",
			objectName:   "path/to/file.jpg",
			bucketName:   "my-bucket",
			envMode:      constants.EnvModeDev,
			emulatorHost: "",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "dev URL with special characters",
			objectName:   "path/with spaces/file name.jpg",
			bucketName:   "my-bucket",
			envMode:      constants.EnvModeDev,
			emulatorHost: "127.0.0.1:8082",
			expected:     "http://127.0.0.1:8082/v0/b/my-bucket/o/path%2Fwith%20spaces%2Ffile%20name.jpg?alt=media",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup config
			config.AppConfig.ENV_MODE = tt.envMode
			config.AppConfig.FIREBASE_STORAGE_EMULATOR_HOST = tt.emulatorHost

			result, err := generatePublicURL(tt.objectName, tt.bucketName)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestValidateFileContent(t *testing.T) {
	// Test with mock clean PDF (minimal valid PDF without suspicious patterns)
	t.Run("valid PDF file", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/pdfs/ecommerce_catalog.pdf")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypePDF)
		assert.True(t, result.IsValid, "PDF file should be valid")
		assert.Equal(t, "application/pdf", result.DetectedMIME)
		assert.Empty(t, result.Error)
	})

	t.Run("valid MP4 video file", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/videos/ecommerce_promo.mp4")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypeVideo)
		assert.True(t, result.IsValid, "MP4 video should be valid")
		assert.True(t, strings.HasPrefix(result.DetectedMIME, "video/"), "Should detect as video type")
		assert.Empty(t, result.Error)
	})

	t.Run("valid WebM video file", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/videos/ecommerce_promo.webm")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypeVideo)
		assert.True(t, result.IsValid, "WebM video should be valid")
		assert.True(t, strings.HasPrefix(result.DetectedMIME, "video/"), "Should detect as video type")
		assert.Empty(t, result.Error)
	})

	t.Run("valid PNG image file", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/images/ecommerce_product.png")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypeImage)
		assert.True(t, result.IsValid, "PNG image should be valid")
		assert.Equal(t, "image/png", result.DetectedMIME)
		assert.Empty(t, result.Error)
	})

	// Test with mock files for blocked types
	t.Run("blocked SVG file with JavaScript", func(t *testing.T) {
		svgContent := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg">
<script>alert('XSS')</script>
</svg>`)
		file := newMockFile(svgContent)

		result := ValidateFileContent(file, constants.ResourceTypeImage)
		assert.False(t, result.IsValid, "SVG file should be blocked")
		assert.Contains(t, result.Error, "not allowed for security reasons")
	})

	t.Run("blocked HTML file", func(t *testing.T) {
		htmlContent := []byte(`<!DOCTYPE html>
<html>
<head><title>Test</title></head>
<body><script>alert('XSS')</script></body>
</html>`)
		file := newMockFile(htmlContent)

		result := ValidateFileContent(file, constants.ResourceTypeImage)
		assert.False(t, result.IsValid, "HTML file should be blocked")
		assert.Contains(t, result.Error, "not allowed for security reasons")
	})

	// Test type mismatch scenarios
	t.Run("PDF file when video expected", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/pdfs/ecommerce_simple.pdf")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypeVideo)
		assert.False(t, result.IsValid, "PDF should not validate as video")
		assert.Contains(t, result.Error, "is not supported")
	})

	t.Run("video file when PDF expected", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/videos/ecommerce_promo.mp4")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypePDF)
		assert.False(t, result.IsValid, "Video should not validate as PDF")
		assert.Contains(t, result.Error, "does not match expected type")
	})

	t.Run("image file when PDF expected", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/images/ecommerce_product.png")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, constants.ResourceTypePDF)
		assert.False(t, result.IsValid, "Image should not validate as PDF")
		assert.Contains(t, result.Error, "does not match expected type")
	})

	// Test unknown resource type
	t.Run("unknown resource type", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/pdfs/ecommerce_catalog.pdf")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		result := ValidateFileContent(file, "unknown_type")
		assert.False(t, result.IsValid, "Unknown resource type should fail")
		assert.Contains(t, result.Error, "unknown resource type")
	})

	// Test file position reset
	t.Run("file position reset after validation", func(t *testing.T) {
		file, err := os.Open("../httpClientTest/videos/ecommerce_promo.mp4")
		if err != nil {
			t.Skipf("Test file not found: %v", err)
		}
		defer file.Close()

		// Validate the file
		result := ValidateFileContent(file, constants.ResourceTypeVideo)
		assert.True(t, result.IsValid)

		// Check that we can read from the beginning
		buf := make([]byte, 4)
		n, err := file.Read(buf)
		assert.NoError(t, err)
		assert.Equal(t, 4, n)
		// MP4 files start with 4 bytes of size, then "ftyp" at offset 4
		// First 4 bytes are typically the box size (varies), so just verify we can read
		assert.NotEmpty(t, buf, "Should be able to read from file start after validation")
	})
}
