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

# Start the server in development mode
export ENV_MODE="dev" && air -c .air.toml

# Start the server in production mode (requires firebase credentials file)
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