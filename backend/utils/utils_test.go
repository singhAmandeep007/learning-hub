package utils

import (
	"learning-hub/config"
	"learning-hub/constants"
	"strconv"
	"strings"
	"testing"
	"time"

	"reflect"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize config for tests
	config.AppConfig = &config.EnvConfig{
		ENV_MODE: constants.EnvModeProd, // Default value
	}
}

func TestGenerateUniqueFilename(t *testing.T) {
	tests := []struct {
		name           string
		originalFile   string
		fileType       string
		wantErr        bool
		validateOutput func(t *testing.T, output string)
	}{
		{
			name:         "empty filename",
			originalFile: "",
			fileType:     "pdf",
			wantErr:      true,
		},
		{
			name:         "valid filename with spaces",
			originalFile: "my document.pdf",
			fileType:     "pdf",
			wantErr:      false,
			validateOutput: func(t *testing.T, output string) {
				if !strings.HasPrefix(output, "pdf/") {
					t.Errorf("expected prefix 'pdf/', got %s", output)
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
			name:         "filename with special characters",
			originalFile: "file@#$%^&*().txt",
			fileType:     "text",
			wantErr:      false,
			validateOutput: func(t *testing.T, output string) {
				if !strings.HasPrefix(output, "text/") {
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
			name:         "filename with multiple dots",
			originalFile: "my.file.name.pdf",
			fileType:     "pdf",
			wantErr:      false,
			validateOutput: func(t *testing.T, output string) {
				if !strings.HasSuffix(output, ".pdf") {
					t.Errorf("expected suffix '.pdf', got %s", output)
				}
				if !strings.Contains(output, "my_file_name") {
					t.Errorf("expected to contain 'my_file_name', got %s", output)
				}
			},
		},
		{
			name:         "filename with unicode characters",
			originalFile: "résumé.pdf",
			fileType:     "pdf",
			wantErr:      false,
			validateOutput: func(t *testing.T, output string) {
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
			got, err := generateUniqueFilename(tt.originalFile, tt.fileType)

			// Check error condition
			if (err != nil) != tt.wantErr {
				t.Errorf("generateUniqueFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Skip validation if we expected an error
			if tt.wantErr {
				return
			}

			// Validate timestamp
			parts := strings.Split(got, "/")
			if len(parts) != 2 {
				t.Errorf("expected format 'type/timestamp_name.ext', got %s", got)
				return
			}

			// Check if timestamp is recent (within last 5 seconds)
			timestampStr := strings.Split(parts[1], "_")[0]
			timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
			if err != nil {
				t.Errorf("invalid timestamp format: %v", err)
				return
			}
			if time.Since(time.Unix(timestamp, 0)) > 5*time.Second {
				t.Errorf("timestamp is not recent: %v", time.Unix(timestamp, 0))
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
		envMode        string
		expectedBucket string
		expectedObject string
		expectError    bool
	}{
		{
			name:           "Valid emulator URL",
			url:            "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/image%2F1748580692_image1.png?alt=media",
			envMode:        constants.EnvModeDev,
			expectedBucket: "learning-hub-81cc6.firebasestorage.app",
			expectedObject: "image/1748580692_image1.png",
			expectError:    false,
		},
		{
			name:           "Valid production URL",
			url:            "https://storage.googleapis.com/my-bucket/path/to/file.jpg",
			envMode:        constants.EnvModeProd,
			expectedBucket: "my-bucket",
			expectedObject: "path/to/file.jpg",
			expectError:    false,
		},
		{
			name:           "Invalid emulator URL format",
			url:            "http://127.0.0.1:8082/invalid/path",
			envMode:        constants.EnvModeDev,
			expectedBucket: "",
			expectedObject: "",
			expectError:    true,
		},
		{
			name:           "Invalid production URL format",
			url:            "https://storage.googleapis.com/invalid",
			envMode:        constants.EnvModeProd,
			expectedBucket: "",
			expectedObject: "",
			expectError:    true,
		},
		{
			name:           "Invalid URL",
			url:            "not-a-url",
			envMode:        constants.EnvModeProd,
			expectedBucket: "",
			expectedObject: "",
			expectError:    true,
		},
		{
			name:           "Emulator URL with special characters",
			url:            "http://127.0.0.1:8082/v0/b/learning-hub-81cc6.firebasestorage.app/o/path%2Fwith%20spaces%20and%20%26%20symbols.jpg?alt=media",
			envMode:        constants.EnvModeDev,
			expectedBucket: "learning-hub-81cc6.firebasestorage.app",
			expectedObject: "path/with spaces and & symbols.jpg",
			expectError:    false,
		},
		{
			name:           "Production URL with special characters",
			url:            "https://storage.googleapis.com/my-bucket/path/with spaces and & symbols.jpg",
			envMode:        constants.EnvModeProd,
			expectedBucket: "my-bucket",
			expectedObject: "path/with spaces and & symbols.jpg",
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the emulator flag
			config.AppConfig.ENV_MODE = tt.envMode

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
