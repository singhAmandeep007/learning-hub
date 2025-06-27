package config

import (
	"learning-hub/constants"
	"log"
	"os"
)

var AppConfig *EnvConfig

type EnvConfig struct {
	ENV_MODE     string `env:"ENV_MODE"`
	PORT         string `env:"PORT"`
	ADMIN_SECRET string `env:"ADMIN_SECRET"`

	CORS_ORIGINS string `env:"CORS_ORIGINS"` // Comma-separated

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

	config.PORT = getEnvOrDefault("PORT", "8000")

	config.CORS_ORIGINS = getEnvOrDefault("CORS_ORIGINS", "*")

	config.FIREBASE_CREDENTIALS_FILE = getEnvOrDefault("FIREBASE_CREDENTIALS_FILE", "firebase_credentials.json")

	config.FIREBASE_PROJECT_ID = getEnvOrDefault("FIREBASE_PROJECT_ID", "learning-hub-81cc6")

	config.FIRESTORE_EMULATOR_HOST = getEnvOrDefault("FIRESTORE_EMULATOR_HOST", "127.0.0.1:8080")
	config.FIREBASE_STORAGE_EMULATOR_HOST = getEnvOrDefault("FIREBASE_STORAGE_EMULATOR_HOST", "127.0.0.1:9199")

	AppConfig = config

	log.Printf("Loaded configuration: %+v", AppConfig)

	return nil
}
