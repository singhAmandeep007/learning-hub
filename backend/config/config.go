package config

import (
	"fmt"
	"learning-hub/constants"
	"os"
	"strconv"
)

var AppConfig *EnvConfig

type EnvConfig struct {
	ENV_MODE     string `env:"ENV_MODE"`
	PORT         string `env:"PORT"`
	ADMIN_SECRET string `env:"ADMIN_SECRET"`

	IS_FIREBASE_EMULATOR bool `env:"IS_FIREBASE_EMULATOR"`

	FIREBASE_CREDENTIALS_FILE string `env:"FIREBASE_CREDENTIALS_FILE"`
	FIREBASE_PROJECT_ID       string `env:"FIREBASE_PROJECT_ID"`

	FIRESTORE_EMULATOR_HOST        string `env:"FIRESTORE_EMULATOR_HOST"`
	FIREBASE_STORAGE_EMULATOR_HOST string `env:"FIREBASE_STORAGE_EMULATOR_HOST"`
}

// getEnvOrDefault retrieves an environment variable or returns a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// LoadConfig loads environment variables into an EnvConfig struct.
func LoadConfig() error {
	config := &EnvConfig{}

	// Load string variables with defaults
	config.ENV_MODE = getEnvOrDefault("ENV_MODE", constants.EnvModeDev)
	config.ADMIN_SECRET = getEnvOrDefault("ADMIN_SECRET", "your-admin-secret-key")

	config.PORT = getEnvOrDefault("PORT", "8080")

	// Load boolean variables with defaults
	isEmulatorStr := getEnvOrDefault("IS_FIREBASE_EMULATOR", "false")
	isEmulator, err := strconv.ParseBool(isEmulatorStr)
	if err != nil {
		return fmt.Errorf("invalid IS_FIREBASE_EMULATOR environment variable: %w", err)
	}
	config.IS_FIREBASE_EMULATOR = isEmulator

	config.FIREBASE_CREDENTIALS_FILE = getEnvOrDefault("FIREBASE_CREDENTIALS_FILE", "")
	config.FIREBASE_PROJECT_ID = getEnvOrDefault("FIREBASE_PROJECT_ID", "")

	config.FIRESTORE_EMULATOR_HOST = "127.0.0.1:" + getEnvOrDefault("FIRESTORE_EMULATOR_HOST_PORT", "8081")
	config.FIREBASE_STORAGE_EMULATOR_HOST = "127.0.0.1:" + getEnvOrDefault("FIREBASE_STORAGE_EMULATOR_HOST_PORT", "8082")

	AppConfig = config

	return nil
}
