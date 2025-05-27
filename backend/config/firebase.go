package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	firebaseStorage "firebase.google.com/go/v4/storage"
	"google.golang.org/api/option"
)

var (
	FirestoreClient *firestore.Client
	StorageClient   *firebaseStorage.Client

	StorageBucket string
	AdminSecret string

	ctx           = context.Background()
)

// InitializeFirebase initializes Firebase services
func InitializeFirebase() error {
	var opts []option.ClientOption

	isUseFirebaseEmulator, err := strconv.ParseBool(os.Getenv("IS_FIREBASE_EMULATOR"))
	if err != nil {
		// Log a warning or fatal error if the env var is set but invalid
		log.Printf("Warning: Invalid value for IS_FIREBASE_EMULATOR ('%v'). Defaulting to false (not using emulator). Error: %v\n", isUseFirebaseEmulator, err)

		isUseFirebaseEmulator = false // Default to not using emulator if parsing fails
	}

	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		// Fallback or error if project ID is crucial and not set
		log.Fatalf("FIREBASE_PROJECT_ID environment variable is not set.")
	}

	var effectiveStorageBucket string

	if isUseFirebaseEmulator {
		opts = append(opts, option.WithoutAuthentication())
		// Although not neccessary - https://firebase.google.com/docs/emulator-suite/connect_firestore#admin_sdks
		if err := os.Setenv("FIRESTORE_EMULATOR_HOST", os.Getenv("FIRESTORE_EMULATOR_HOST")); err != nil {
			log.Fatalf("Init emulator db err: %v", err)
		}
		// Although not neccessary - https://firebase.google.com/docs/emulator-suite/connect_storage#web
		if err := os.Setenv("FIREBASE_STORAGE_EMULATOR_HOST", os.Getenv("FIREBASE_STORAGE_EMULATOR_HOST")); err != nil {
			log.Fatalf("Init emulator db err: %v", err)
		}

		effectiveStorageBucket = projectID + ".appspot.com"
	}else {
		opts = append(opts, option.WithCredentialsFile(os.Getenv("FIREBASE_CREDENTIALS_FILE")))

		effectiveStorageBucket = os.Getenv("FIREBASE_STORAGE_BUCKET")
		if effectiveStorageBucket == "" {
			// Default to <project_id>.firebasestorage.com for live environment if not specified
			effectiveStorageBucket = projectID + ".firebasestorage.com"
			log.Printf("FIREBASE_STORAGE_BUCKET not set, defaulting to: %s for live environment", effectiveStorageBucket)
		}
	}

	conf := &firebase.Config{
		ProjectID: projectID,
				StorageBucket: effectiveStorageBucket, // This is mainly for firebaseApp.Storage(ctx).DefaultBucket()

	}
	
	firebaseApp, err := firebase.NewApp(ctx, conf, opts...)
	if err != nil {
		return fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Initialize Firestore
	FirestoreClient, err = firebaseApp.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("error initializing firestore: %v", err)
	}

	// Initialize Cloud Storage
	StorageClient, err = firebaseApp.Storage(ctx) // storage.NewClient(ctx,opts...)
	if err != nil {
		return fmt.Errorf("error initializing storage: %v", err)
	}

	AdminSecret = os.Getenv("ADMIN_SECRET")
	if AdminSecret == "" {
		AdminSecret = "your-admin-secret-key" // Should be set via environment variable
	}

	StorageBucket = effectiveStorageBucket

	return nil
}

// cleanup closes all clients
func CloseFirebase() {
	if FirestoreClient != nil {
		FirestoreClient.Close()
	}
	if StorageClient != nil {
		// StorageClient.Close()
	}
}