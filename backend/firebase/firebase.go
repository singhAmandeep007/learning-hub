package firebase

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"

	firebase "firebase.google.com/go/v4"

	"google.golang.org/api/option"

	"learning-hub/config"
	"learning-hub/constants"
)

var (
	FirestoreClient *firestore.Client
	StorageClient   *storage.Client

	StorageBucket string

	ctx    = context.Background()
	cancel context.CancelFunc
)

// InitializeFirebase initializes Firebase services
func InitializeFirebase() error {
	// Create a context with cancellation for proper cleanup
	ctx, cancel = context.WithCancel(ctx)

	opts, err := buildFirebaseConfig()
	if err != nil {
		cancel()
		return fmt.Errorf("error: failed to build Firebase config: %w", err)
	}

	StorageBucket = config.AppConfig.FIREBASE_PROJECT_ID + ".firebasestorage.app"

	config := &firebase.Config{
		ProjectID: config.AppConfig.FIREBASE_PROJECT_ID,
	}

	firebaseApp, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		cancel()
		return fmt.Errorf("error initializing firebase app: %w", err)
	}

	// Initialize Firestore
	FirestoreClient, err = firebaseApp.Firestore(ctx)
	if err != nil {
		cancel()
		return fmt.Errorf("error initializing firestore: %w", err)
	}

	// opts = append([]option.ClientOption{option.WithScopes()}, opts...)

	// Initialize Cloud Storage
	StorageClient, err = storage.NewClient(ctx, opts...)
	if err != nil {
		cancel()
		return fmt.Errorf("error initializing storage: %v", err)
	}

	log.Printf("Firebase initialized successfully with bucket: %s", StorageBucket)
	return nil
}

// buildFirebaseConfig builds the Firebase configuration based on environment
func buildFirebaseConfig() (firebaseOptions []option.ClientOption, error error) {
	var opts []option.ClientOption

	isDev := config.AppConfig.ENV_MODE == constants.EnvModeDev

	if isDev {
		opts = append(opts, option.WithoutAuthentication())

		// set emulator hosts
		setEmulatorHosts()

		log.Printf("Using Firebase emulator mode")
	} else {
		credentialsFile := config.AppConfig.FIREBASE_CREDENTIALS_FILE
		if credentialsFile == "" {
			return nil, fmt.Errorf("FIREBASE_CREDENTIALS_FILE is required for production mode")
		}

		if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
			return nil, fmt.Errorf("credentials file not found: %s", credentialsFile)
		}

		opts = append(opts, option.WithCredentialsFile(credentialsFile))

		log.Printf("Using Firebase production mode")
	}

	return opts, nil
}

// setEmulatorHosts checks the emulator host environment variables if not already set
func setEmulatorHosts() {
	firebaseStorageEmulatorHost := config.AppConfig.FIREBASE_STORAGE_EMULATOR_HOST

	os.Setenv("FIRESTORE_EMULATOR_HOST", config.AppConfig.FIRESTORE_EMULATOR_HOST)
	os.Setenv("FIREBASE_STORAGE_EMULATOR_HOST", firebaseStorageEmulatorHost)

	// need to set because we are using "cloud.google.com/go/storage" to create new storage client - https://github.com/firebase/firebase-admin-go/blob/570427a0f270b9adb061f54187a2b033548c3c9e/storage/storage.go#L38
	os.Setenv("STORAGE_EMULATOR_HOST", firebaseStorageEmulatorHost)
}

func CloseFirebase() error {
	var errors []error

	// Close Firestore client
	if FirestoreClient != nil {
		if err := FirestoreClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("error closing Firestore client: %w", err))
		}
		FirestoreClient = nil
	}

	if StorageClient != nil {
		if err := StorageClient.Close(); err != nil {
			errors = append(errors, fmt.Errorf("error closing Storage client: %w", err))
		}
		StorageClient = nil
	}

	// Cancel the context to clean up any ongoing operations
	if cancel != nil {
		cancel()
		cancel = nil
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during cleanup: %v", errors)
	}

	log.Printf("Firebase clients closed successfully")
	return nil
}
