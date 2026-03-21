package firebase

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"

	"google.golang.org/api/option"

	"learninghub/config"
	"learninghub/constants"
	"learninghub/pkg/logger"
)

var (
	FirestoreClient *firestore.Client
	StorageClient   *storage.Client

	StorageBucket string

	ctx    = context.Background()
	cancel context.CancelFunc
)

// createFirestoreClientWithDatabase creates a Firestore client for a specific database
func createFirestoreClientWithDatabase(ctx context.Context, projectID, databaseID string, opts []option.ClientOption) (*firestore.Client, error) {
	return firestore.NewClientWithDatabase(ctx, projectID, databaseID, opts...)
}

// InitializeFirebase initializes Firebase services
func InitializeFirebase() error {
	// Create a context with cancellation for proper cleanup
	ctx, cancel = context.WithCancel(ctx)

	opts, err := buildFirebaseConfig()
	if err != nil {
		cancel()
		return fmt.Errorf("error: failed to build Firebase config: %w", err)
	}

	StorageBucket = config.AppConfig.FIREBASE_STORAGE_BUCKET

	// Create Firestore client with database
	FirestoreClient, err = createFirestoreClientWithDatabase(ctx, config.AppConfig.FIREBASE_PROJECT_ID, config.AppConfig.FIRESTORE_DB_ID, opts)
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

		logger.Infof("Using Firebase emulator mode")
	}

	return opts, nil
}

// setEmulatorHosts checks the emulator host environment variables if not already set
func setEmulatorHosts() {
	os.Setenv("FIRESTORE_EMULATOR_HOST", config.AppConfig.FIRESTORE_EMULATOR_HOST)

	firebaseStorageEmulatorHost := config.AppConfig.FIREBASE_STORAGE_EMULATOR_HOST
	os.Setenv("FIREBASE_STORAGE_EMULATOR_HOST", firebaseStorageEmulatorHost)

	// need to set because we are using "cloud.google.com/go/storage" to create new storage client - https://github.com/firebase/firebase-admin-go/blob/570427a0f270b9adb061f54187a2b033548c3c9e/storage/storage.go#L38
	os.Setenv("STORAGE_EMULATOR_HOST", firebaseStorageEmulatorHost)
}

func CloseFirebase() {
	// Close Firestore client
	if FirestoreClient != nil {
		logger.Infof("Closing Firestore client")
		if err := FirestoreClient.Close(); err != nil {
			logger.Infof("error closing Firestore client: %v", err)
		} else {
			logger.Infof("Firestore client closed successfully")
		}
		FirestoreClient = nil
	}

	if StorageClient != nil {
		logger.Infof("Closing Storage client")
		if err := StorageClient.Close(); err != nil {
			logger.Infof("error closing Storage client: %v", err)
		} else {
			logger.Infof("Storage client closed successfully")
		}
		StorageClient = nil
	}

	// Cancel the context to clean up any ongoing operations
	if cancel != nil {
		cancel()
		cancel = nil
	}

	logger.Infof("Firebase clients shutdown process completed.")
}
