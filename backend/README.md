## Prerequisites

- Firebase Tools (for Firebase emulators)

## Starting firebase emulators

```bash
# Check if Firebase CLI is installed
firebase --version
# Start Firebase emulators
firebase emulators:start --only firestore,storage
```

## Starting the Backend

```bash
# Install Air - Live reload for Go apps
go install github.com/air-verse/air@latest

# Check if Air is installed
which air

# Install dependencies
go mod download
go mod tidy

# Start the server in dev mode
export ENV_MODE="dev" && air -c .air.toml

# Start the server in prod mode (requires firebase credentials file in backend directory named `firebase_credentials.json`)
export ENV_MODE="prod" && air -c .air.toml
```

## Running Tests

```bash
# Run all tests
go test ./...
# Run tests with coverage
go test -cover ./...
```

## Create build

```bash
# Build the application
go build -o ./tmp/main .
```

## Running shell script to seed, update and delete resources

```shell
cd httpClientTest

chmod +x resources.sh

./resources.sh
```


## Testing docker locally

### Build the Docker image
```bash
docker build -t learning-hub-backend-prod -f backend/Dockerfile.prod backend
```

### Run the Docker container
```bash
docker run --rm -it --name test-backend-prod -p 8000:8000 -e ENV_MODE=prod -e FIREBASE_PROJECT_ID=learning-hub-81cc6 -e FIREBASE_CREDENTIALS_FILE=/app/firebase_credentials.json -e ADMIN_SECRET=your-secure-admin-secret -e CORS_ORIGINS=http://localhost:3000 -v $(pwd)/backend/firebase_credentials.json:/app/firebase_credentials.json:ro learning-hub-backend-prod
```