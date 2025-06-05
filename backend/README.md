## Starting the Backend

```bash
# Install Air - Live reload for Go apps
go install github.com/air-verse/air@latest

# Check if Air is installed
which air

# Install dependencies
go mod tidy

# Start the server in development mode
export ENV_MODE="dev" && air -c .air.toml

# Start the server in production mode
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