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
)

var (
	FirestoreClient *firestore.Client
	StorageClient   *storage.Client

	StorageBucket string

	ctx           = context.Background()
	cancel          context.CancelFunc
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

	storageBucket := *config.AppConfig.FIREBASE_PROJECT_ID + ".firebasestorage.app"

	conf := &firebase.Config{
		ProjectID:     *config.AppConfig.FIREBASE_PROJECT_ID,
		// StorageBucket: storageBucket,
	}

	firebaseApp, err := firebase.NewApp(ctx, conf, opts...)
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

	// // Initialize Firebase Storage
	// StorageClient, err = firebaseApp.Storage(ctx)
	// if err != nil {
	// 	cancel()
	// 	// Clean up Firestore if Storage initialization fails
	// 	if FirestoreClient != nil {
	// 		CloseFirebase()
	// 	}
	// 	return fmt.Errorf("error initializing storage: %w", err)
	// }

	StorageBucket = storageBucket

	log.Printf("Firebase initialized successfully with bucket: %s", StorageBucket)
	return nil
}

// buildFirebaseConfig builds the Firebase configuration based on environment
func  buildFirebaseConfig() (firebaseOptions []option.ClientOption, error error) {
	var opts []option.ClientOption

	isEmulator := config.AppConfig.IS_FIREBASE_EMULATOR

	if isEmulator {
		opts = append(opts, option.WithoutAuthentication())
		
		// check emulator hosts if not already set
		checkEmulatorHosts()
		
		log.Printf("Using Firebase emulator mode")
	} else {
		credentialsFile := *config.AppConfig.FIREBASE_CREDENTIALS_FILE
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

// checkEmulatorHosts checks the emulator host environment variables if not already set
func checkEmulatorHosts() {
	firebaseStorageEmulatorHost := *config.AppConfig.FIREBASE_STORAGE_EMULATOR_HOST

	if *config.AppConfig.FIRESTORE_EMULATOR_HOST == "" {
		log.Printf("Warning: FIRESTORE_EMULATOR_HOST not set, emulator connection may fail")
	}
	if firebaseStorageEmulatorHost == "" {
		log.Printf("Warning: FIREBASE_STORAGE_EMULATOR_HOST not set, emulator connection may fail")
	}

	// need to set because we are using "cloud.google.com/go/storage" to create new storage client
	// https://github.com/firebase/firebase-admin-go/blob/570427a0f270b9adb061f54187a2b033548c3c9e/storage/storage.go#L38
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