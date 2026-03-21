package logger

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestLogLevelString tests the String method of LogLevel
func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, LogLevelDebug},
		{InfoLevel, LogLevelInfo},
		{WarnLevel, LogLevelWarn},
		{ErrorLevel, LogLevelError},
		{FatalLevel, LogLevelFatal},
		{LogLevel(999), LogLevelUnknown}, // Unknown level
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}

// TestLogLevelToZapLevel tests the ToZapLevel method of LogLevel
func TestLogLevelToZapLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{LogLevel(999), "info"}, // Default to info for unknown
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level_%d", tt.level), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.ToZapLevel().String())
		})
	}
}

// CustomTestDestination is a test destination that captures log entries
type CustomTestDestination struct {
	Entries []LogEntry
	mu      sync.Mutex
}

func NewTestDestination() *CustomTestDestination {
	return &CustomTestDestination{
		Entries: make([]LogEntry, 0),
	}
}

func (d *CustomTestDestination) Write(entry LogEntry) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Entries = append(d.Entries, entry)
	return nil
}

func (d *CustomTestDestination) Close() error {
	return nil
}

func (d *CustomTestDestination) GetEntries() []LogEntry {
	d.mu.Lock()
	defer d.mu.Unlock()
	// Return a copy to avoid race conditions
	entries := make([]LogEntry, len(d.Entries))
	copy(entries, d.Entries)
	return entries
}

func (d *CustomTestDestination) GetEntryCount() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.Entries)
}

// TestLoggerCreation tests the creation of a logger with options
func TestLoggerCreation(t *testing.T) {
	// Test default logger
	logger := New()
	logger.AddDestination("custom", NewTestDestination())
	assert.Equal(t, "", logger.serviceName)
	assert.Equal(t, InfoLevel, logger.minLevel)
	assert.False(t, logger.isProd)
	assert.Len(t, logger.destinations, 1)
	assert.Contains(t, logger.destinations, "custom")
	assert.Equal(t, []string{}, logger.defaultDests)

	// Test with options
	logger = New(
		WithServiceName("test-service"),
		WithMinLevel(DebugLevel),
		WithProduction(true),
		WithDefaultDestinations("test"),
	)
	assert.Equal(t, "test-service", logger.serviceName)
	assert.Equal(t, InfoLevel, logger.minLevel) // Should be InfoLevel because isProd=true
	assert.True(t, logger.isProd)
	assert.Equal(t, []string{"test"}, logger.defaultDests)

	// Clean up
	logger.Close()
}

// TestLoggerLevelFiltering tests that logs below minimum level are filtered
func TestLoggerLevelFiltering(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger with InfoLevel minimum
	logger := New(
		WithServiceName("test-service"),
		WithMinLevel(InfoLevel),
	)
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log messages at different levels
	logger.Debug("Debug message", nil)
	logger.Info("Info message", nil)
	logger.Warn("Warn message", nil)

	// Verify only Info and Warn were logged
	assert.Len(t, testDest.Entries, 2)
	assert.Equal(t, "Info message", testDest.Entries[0].Message)
	assert.Equal(t, "Warn message", testDest.Entries[1].Message)

	// Clean up
	logger.Close()
}

// TestLoggerProductionMode tests that debug logs are filtered in production mode
func TestLoggerProductionMode(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger in production mode with DebugLevel
	logger := New(
		WithServiceName("test-service"),
		WithMinLevel(DebugLevel),
		WithProduction(true), // This should override the min level for debug
	)
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log messages at different levels
	logger.Debug("Debug message", nil)
	logger.Info("Info message", nil)

	// Verify only Info was logged (Debug filtered in production)
	assert.Len(t, testDest.Entries, 1)
	assert.Equal(t, "Info message", testDest.Entries[0].Message)

	// Clean up
	logger.Close()
}

// TestLoggerFields tests that fields are properly included in log entries
func TestLoggerFields(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger
	logger := New(
		WithServiceName("test-service"),
	)
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log with fields
	fields := map[string]interface{}{
		"string": "value",
		"number": 42,
		"bool":   true,
	}
	logger.Info("Test message", fields)

	// Verify fields are included
	assert.Len(t, testDest.Entries, 1)
	assert.Equal(t, "Test message", testDest.Entries[0].Message)
	assert.Equal(t, fields, testDest.Entries[0].Fields)
	assert.Equal(t, "test-service", testDest.Entries[0].ServiceName)

	// Clean up
	logger.Close()
}

// TestLoggerMultipleDestinations tests logging to multiple destinations
func TestLoggerMultipleDestinations(t *testing.T) {
	testDest1 := NewTestDestination()
	testDest2 := NewTestDestination()

	// Create logger with multiple destinations
	logger := New(
		WithServiceName("test-service"),
	)
	logger.AddDestination("test1", testDest1)
	logger.AddDestination("test2", testDest2)
	logger.SetDefaultDestinations("test1", "test2")

	// Log a message
	logger.Info("Test message", nil)

	// Verify it went to both destinations
	assert.Len(t, testDest1.Entries, 1)
	assert.Len(t, testDest2.Entries, 1)
	assert.Equal(t, "Test message", testDest1.Entries[0].Message)
	assert.Equal(t, "Test message", testDest2.Entries[0].Message)

	// Test explicit destination
	logger.Error("Error message", nil, "test1")

	// Verify it only went to test1
	assert.Len(t, testDest1.Entries, 2)
	assert.Len(t, testDest2.Entries, 1)
	assert.Equal(t, "Error message", testDest1.Entries[1].Message)

	// Clean up
	logger.Close()
}

// TestLoggerFileDestination tests the file destination option
func TestLoggerFileDestination(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	// Create logger with file destination
	logger := New(
		WithServiceName("test-service"),
		WithFileDestination(logFile, 10, 1, 1, false),
		WithDefaultDestinations(FileLogger),
	)

	// Log some messages
	logger.Info("Info to file", map[string]interface{}{"key": "value"})
	logger.Errorf("Error to file")

	// Close to ensure everything is flushed
	logger.Close()

	// Verify log file was created and contains the messages
	f, err := os.Open(logFile)
	require.NoError(t, err)
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.NoError(t, scanner.Err())

	// Should have two log entries
	assert.Len(t, lines, 2)

	// Parse and verify JSON entries
	var infoEntry, errorEntry map[string]interface{}

	err = json.Unmarshal([]byte(lines[0]), &infoEntry)
	require.NoError(t, err)
	assert.Equal(t, "Info to file", infoEntry["msg"])
	assert.Equal(t, "test-service", infoEntry["service"])
	assert.Equal(t, "value", infoEntry["key"])

	err = json.Unmarshal([]byte(lines[1]), &errorEntry)
	require.NoError(t, err)
	assert.Equal(t, "Error to file", errorEntry["msg"])
	assert.Equal(t, "test-service", errorEntry["service"])
}

// TestLoggerAddRemoveDestination tests adding and removing destinations
func TestLoggerAddRemoveDestination(t *testing.T) {
	testDest := NewTestDestination()

	// Create logger
	logger := New()
	logger.AddDestination("test", testDest)
	logger.SetDefaultDestinations("test")

	// Log a message
	logger.Info("Test message", nil)
	assert.Len(t, testDest.Entries, 1)

	// Remove the destination
	logger.RemoveDestination("test")
	logger.Info("Another message", nil)

	// The message shouldn't have been logged to the removed destination
	assert.Len(t, testDest.Entries, 1)

	// Clean up
	logger.Close()
}

// TestConsoleDestination tests the console destination implementation
func TestConsoleDestination(t *testing.T) {
	// Redirect stdout temporarily
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create console destination
	consoleDest := NewConsoleDestination()

	// Log a test entry
	entry := LogEntry{
		ServiceName: "test-service",
		Level:       InfoLevel,
		Message:     "Console test message",
		Fields:      map[string]interface{}{"test": true},
	}
	err := consoleDest.Write(entry)
	require.NoError(t, err)

	// Close the writer to capture output
	w.Close()

	// Read the output before closing the destination
	var buf strings.Builder
	_, err = io.Copy(&buf, r)
	require.NoError(t, err)

	// Close the destination BEFORE restoring stdout
	// This ensures we sync while the redirect is still active
	_ = consoleDest.Close()

	// Now restore stdout
	os.Stdout = oldStdout

	// Verify output contains the message
	output := buf.String()
	assert.Contains(t, output, "Console test message")
	assert.Contains(t, output, "test-service")
	assert.Contains(t, output, `"test": true`)
}

// TestFatalExit tests that Fatal logs cause program exit
// This test is skipped because it would terminate the test process
func TestFatalExit(t *testing.T) {
	if os.Getenv("TEST_FATAL_EXIT") == "1" {
		logger := New()
		logger.Fatal("This should exit", nil)
		// Should not reach here
		t.Fail()
	} else {
		t.Skip("Skipping fatal exit test as it would terminate the process")
	}
}

// TestWithFileDestinationOption tests the WithFileDestination option
func TestWithFileDestinationOption(t *testing.T) {
	// Create a temporary directory for test logs
	tempDir, err := os.MkdirTemp("", "logger_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logFile := filepath.Join(tempDir, "test.log")

	// Create logger with file destination option
	logger := New(
		WithFileDestination(logFile, 10, 1, 1, false),
		WithDefaultDestinations(FileLogger),
	)

	// Verify the file destination was added
	assert.Contains(t, logger.destinations, "file")
	assert.Contains(t, logger.defaultDests, "file")

	// Log a test message
	logger.Info("Test with file destination option", nil)

	// Close to ensure everything is flushed
	logger.Close()

	// Verify log file was created
	_, err = os.Stat(logFile)
	assert.NoError(t, err)
}

// TestGlobalLogger tests the global logger functionality
func TestGlobalLogger(t *testing.T) {
	// Test InitGlobal and GetGlobal
	InitGlobal(
		WithServiceName("test-global"),
		WithMinLevel(DebugLevel),
		WithConsoleDestination(),
		WithDefaultDestinations(ConsoleLogger),
	)

	globalLogger := GetGlobal()
	assert.NotNil(t, globalLogger)
	assert.Equal(t, "test-global", globalLogger.serviceName)
	assert.Equal(t, DebugLevel, globalLogger.minLevel)

	// Clean up
	CloseGlobal()
}

// TestGlobalLoggerFallback tests that GetGlobal creates a fallback logger
func TestGlobalLoggerFallback(t *testing.T) {
	// Reset global logger to nil (simulate uninitialized state)
	globalLogger = nil

	// GetGlobal should create a fallback logger
	globalLogger := GetGlobal()
	assert.NotNil(t, globalLogger)
	assert.Equal(t, "", globalLogger.serviceName)                // Default empty service name
	assert.Contains(t, globalLogger.destinations, ConsoleLogger) // Should have console destination

	// Clean up
	CloseGlobal()
}

// TestGlobalConvenienceFunctions tests the global convenience functions
func TestGlobalConvenienceFunctions(t *testing.T) {
	testDest := NewTestDestination()

	// Initialize global logger with test destination
	InitGlobal(
		WithServiceName("test-global"),
		WithMinLevel(DebugLevel),
	)

	globalLogger := GetGlobal()
	globalLogger.AddDestination("test", testDest)
	globalLogger.SetDefaultDestinations("test")

	// Test global convenience functions
	Debugf("Debug message %d", 1)
	Infof("Info message %s", "test")
	Warnf("Warn message %v", true)
	Errorf("Error message %f", 3.14)

	Debug("Debug message", map[string]interface{}{"key": "value"})
	Info("Info message", map[string]interface{}{"key": "value"})

	// Verify messages were logged
	assert.Len(t, testDest.Entries, 6)
	assert.Equal(t, "Debug message 1", testDest.Entries[0].Message)
	assert.Equal(t, "Info message test", testDest.Entries[1].Message)
	assert.Equal(t, "Warn message true", testDest.Entries[2].Message)
	assert.Equal(t, "Error message 3.140000", testDest.Entries[3].Message)

	assert.Equal(t, "Debug message", testDest.Entries[4].Message)
	assert.Equal(t, "Info message", testDest.Entries[5].Message)

	// Verify all have correct service name
	for _, entry := range testDest.Entries {
		assert.Equal(t, "test-global", entry.ServiceName)
	}

	// Clean up
	CloseGlobal()
}

// TestGlobalLoggerProductionMode tests global logger in production mode
func TestGlobalLoggerProductionMode(t *testing.T) {
	testDest := NewTestDestination()

	// Initialize global logger in production mode
	InitGlobal(
		WithServiceName("test-global-prod"),
		WithMinLevel(DebugLevel),
		WithProduction(true), // This should override min level to InfoLevel
	)

	globalLogger := GetGlobal()
	globalLogger.AddDestination("test", testDest)
	globalLogger.SetDefaultDestinations("test")

	// Test that debug is filtered in production
	Debugf("Debug message")
	Infof("Info message")

	// Verify only info message was logged
	assert.Len(t, testDest.Entries, 1)
	assert.Equal(t, "Info message", testDest.Entries[0].Message)
	assert.Equal(t, InfoLevel, testDest.Entries[0].Level)

	// Clean up
	CloseGlobal()
}

// TestGlobalLoggerConcurrency tests that global logger is thread-safe
func TestGlobalLoggerConcurrency(t *testing.T) {
	testDest := NewTestDestination()

	InitGlobal(
		WithServiceName("test-concurrent"),
		WithMinLevel(InfoLevel),
	)

	globalLogger := GetGlobal()
	globalLogger.AddDestination("test", testDest)
	globalLogger.SetDefaultDestinations("test")

	// Run concurrent logging operations
	const numGoroutines = 10
	const messagesPerGoroutine = 5

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				Infof("Message from goroutine %d, iteration %d", id, j)
			}
		}(i)
	}

	wg.Wait()

	// Verify all messages were logged
	expectedMessages := numGoroutines * messagesPerGoroutine
	assert.Equal(t, expectedMessages, testDest.GetEntryCount())

	// Clean up
	CloseGlobal()
}

// TestGlobalLoggerReinitialization tests reinitializing the global logger
func TestGlobalLoggerReinitialization(t *testing.T) {
	// First initialization
	InitGlobal(WithServiceName("first-service"))
	firstLogger := GetGlobal()
	assert.Equal(t, "first-service", firstLogger.serviceName)

	// Second initialization should replace the first
	InitGlobal(WithServiceName("second-service"))
	secondLogger := GetGlobal()
	assert.Equal(t, "second-service", secondLogger.serviceName)

	// Should be a different instance
	assert.NotEqual(t, firstLogger, secondLogger)

	// Clean up
	CloseGlobal()
}
