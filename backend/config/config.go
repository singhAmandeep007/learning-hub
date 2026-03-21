package config

import (
	"os"

	"learninghub/constants"
	"learninghub/pkg/logger"
)

var AppConfig *EnvConfig

type EnvConfig struct {
	ENV_MODE string `env:"ENV_MODE"`
	PORT     string `env:"PORT"`

	CORS_ORIGINS string `env:"CORS_ORIGINS"` // Comma-separated

	FIREBASE_PROJECT_ID string `env:"FIREBASE_PROJECT_ID"`

	FIRESTORE_EMULATOR_HOST        string `env:"FIRESTORE_EMULATOR_HOST"`
	FIREBASE_STORAGE_EMULATOR_HOST string `env:"FIREBASE_STORAGE_EMULATOR_HOST"`

	FIRESTORE_DB_ID         string `env:"FIRESTORE_DB_ID"`
	FIREBASE_STORAGE_BUCKET string `env:"FIREBASE_STORAGE_BUCKET"`
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

	// Load env as prod by default
	config.ENV_MODE = getEnvOrDefault("ENV_MODE", constants.EnvModeProd)

	config.PORT = getEnvOrDefault("PORT", "8000")

	config.CORS_ORIGINS = getEnvOrDefault("CORS_ORIGINS", "http://")

	config.FIREBASE_PROJECT_ID = getEnvOrDefault("FIREBASE_PROJECT_ID", "learninghub-81cc6")

	config.FIRESTORE_EMULATOR_HOST = getEnvOrDefault("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	config.FIREBASE_STORAGE_EMULATOR_HOST = getEnvOrDefault("FIREBASE_STORAGE_EMULATOR_HOST", "127.0.0.1:9199")

	config.FIRESTORE_DB_ID = getEnvOrDefault("FIRESTORE_DB_ID", "learninghub")

	config.FIREBASE_STORAGE_BUCKET = getEnvOrDefault("FIREBASE_STORAGE_BUCKET", config.FIREBASE_PROJECT_ID+".firebasestorage.app")

	AppConfig = config

	logger.Infof("Loaded configuration: %+v", AppConfig)

	return nil
}
