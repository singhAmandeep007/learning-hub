# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: help dev-local stop-services install-tools install-deps docker-dev docker-dev-no-cache docker-dev-stop e2e-docker e2e-docker-vrt e2e-docker-vrt-update e2e-docker-stop e2e-local e2e-local-vrt e2e-local-vrt-update clean docker-clean status

# Function to wait for a port to be open
# Usage: $(call wait_for_port, <port_number>, <service_name>)
define wait_for_port
	@echo "$(YELLOW)Waiting for $(2) on port $(1)...$(NC)"
	@until nc -z localhost $(1) >/dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done
	@echo "$(GREEN)$(2) is available!$(NC)"
endef

# Development environment - runs all services concurrently
dev-local:
	@echo "$(GREEN)Starting development environment...$(NC)"
	@echo "$(YELLOW)Starting Firebase emulators...$(NC)"
	@cd backend && firebase emulators:start --only firestore,storage &
	$(call wait_for_port, 4000, Firebase Emulators)
	@echo "$(YELLOW)Starting backend server with Air...$(NC)"
	@cd backend && ENV_MODE=dev VALID_PRODUCTS=ecomm CORS_ORIGINS=http://localhost:3000,http://127.0.0.1:3000,http://localhost:5173,http://127.0.0.1:5173 FIREBASE_PROJECT_ID=learninghub-81cc6 FIRESTORE_EMULATOR_HOST=127.0.0.1:8080 FIREBASE_STORAGE_EMULATOR_HOST=127.0.0.1:9199 air -c .air.toml &
	$(call wait_for_port, 8000, Backend Server)
	@echo "$(YELLOW)Setting up Node.js version...$(NC)"
	@cd frontend && bash -c "source ~/.nvm/nvm.sh && nvm use"
	@echo "$(YELLOW)Starting frontend development server...$(NC)"
	@cd frontend && npm run dev &
	@echo "$(GREEN)All services started!$(NC) Run $(RED)make stop-services$(NC) in a new terminal to stop all services."
	@echo "$(GREEN)Backend: http://localhost:8000$(NC)"
	@echo "$(GREEN)Frontend: http://localhost:3000$(NC)"
	@echo "$(GREEN)Firebase UI: http://localhost:4000$(NC)"
	@wait

# Stop all running services
# firebase emulator not exiting properly - https://github.com/firebase/firebase-tools/issues/3578#issuecomment-1024876668
stop-services:
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@pkill -f "air -c" || true
	@pkill -f "npm run dev" || true
	@pkill -f "firebase emulators:start" || true
	@echo "$(GREEN)All services stopped$(NC)"

# Install global tools
install-tools:
	@echo "$(GREEN)Installing global tools...$(NC)"
	@echo "$(YELLOW)Installing Firebase CLI...$(NC)"
	@npm install -g firebase-tools@14.2.0
	@echo "$(YELLOW)Installing Air (Go hot reload)...$(NC)"
	@go install github.com/air-verse/air@v1.52.3
	@echo "$(GREEN)Global tools installed successfully!$(NC)"

# Install dependencies
install-deps:
	@echo "$(GREEN)Installing frontend dependencies...$(NC)"
	@cd frontend && npm install
	@echo "$(GREEN)Installing backend dependencies...$(NC)"
	@cd backend && go mod download && go mod verify

# Docker dev environment
docker-dev:
	@echo "$(GREEN)🚀 Starting development environment with docker...$(NC)"
	@echo "$(YELLOW)Stopping existing containers and removing volumes...$(NC)"
	docker compose -f docker-compose.dev.yml down -v 2>/dev/null || true
	@echo "$(YELLOW)Building containers...$(NC)"
	docker compose -f docker-compose.dev.yml build
	@echo "$(YELLOW)Starting containers...$(NC)"
	docker compose -f docker-compose.dev.yml up

docker-dev-no-cache:
	@echo "$(GREEN)🚀 Starting development environment with docker...$(NC)"
	@echo "$(YELLOW)Stopping existing containers and removing volumes...$(NC)"
	docker compose -f docker-compose.dev.yml down -v 2>/dev/null || true
	@echo "$(YELLOW)Building containers (no cache)...$(NC)"
	docker compose -f docker-compose.dev.yml build --no-cache
	@echo "$(YELLOW)Starting containers...$(NC)"
	docker compose -f docker-compose.dev.yml up

# Stop all services
docker-dev-stop:
	@echo "🛑 Stopping all dev docker services..."
	docker compose -f docker-compose.dev.yml down

# Run dedicated e2e docker stack (frontend + backend + firebase + playwright)
e2e-docker:
	@echo "$(GREEN)🚀 Starting E2E Docker stack...$(NC)"
	docker compose -f docker-compose.e2e.yml down -v 2>/dev/null || true
	docker compose -f docker-compose.e2e.yml up --build --abort-on-container-exit --exit-code-from e2e

# Run visual regression tests in Docker stack (real backend + frontend)
e2e-docker-vrt:
	@echo "$(GREEN)🖼️  Running visual regression tests in E2E Docker stack...$(NC)"
	docker compose -f docker-compose.e2e.yml down -v 2>/dev/null || true
	E2E_TEST_COMMAND="npm run test:visual" docker compose -f docker-compose.e2e.yml up --build --abort-on-container-exit --exit-code-from e2e

# Refresh visual snapshot baselines in Docker stack (real backend + frontend)
e2e-docker-vrt-update:
	@echo "$(GREEN)🔁 Updating visual regression snapshots in E2E Docker stack...$(NC)"
	docker compose -f docker-compose.e2e.yml down -v 2>/dev/null || true
	E2E_TEST_COMMAND="npm run test:visual:update" docker compose -f docker-compose.e2e.yml up --build --abort-on-container-exit --exit-code-from e2e

e2e-docker-stop:
	@echo "🛑 Stopping E2E docker services..."
	docker compose -f docker-compose.e2e.yml down -v --remove-orphans

# Run e2e tests against locally running app services
e2e-local:
	@echo "$(GREEN)🧪 Running E2E tests against local services...$(NC)"
	@cd e2e && npm ci && npm run install:browsers && E2E_BASE_URL=http://localhost:3000 E2E_API_BASE_URL=http://localhost:8000 E2E_PRODUCT=$${E2E_PRODUCT:-ecomm} npm test

# Run visual regression tests against locally running services
e2e-local-vrt:
	@echo "$(GREEN)🖼️  Running visual regression tests against local services...$(NC)"
	@cd e2e && npm ci && npm run install:browsers && E2E_BASE_URL=http://localhost:3000 E2E_API_BASE_URL=http://localhost:8000 E2E_PRODUCT=$${E2E_PRODUCT:-ecomm} npm run test:visual

# Update visual snapshot baselines against locally running services
e2e-local-vrt-update:
	@echo "$(GREEN)🔁 Updating visual regression snapshots against local services...$(NC)"
	@cd e2e && npm ci && npm run install:browsers && E2E_BASE_URL=http://localhost:3000 E2E_API_BASE_URL=http://localhost:8000 E2E_PRODUCT=$${E2E_PRODUCT:-ecomm} npm run test:visual:update

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@cd frontend && rm -rf dist node_modules/.vite
	@cd backend && rm -rf tmp
	@echo "$(GREEN)Clean complete$(NC)"

# Help command
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo "  $(YELLOW)dev-local$(NC)          - Start development environment (hot reload)"
	@echo "  $(YELLOW)stop-services$(NC)      - Stop all running services"
	@echo "  $(YELLOW)docker-dev$(NC)         - Start development environment with Docker"
	@echo "  $(YELLOW)docker-dev-stop$(NC)    - Stop all dev docker services"
	@echo "  $(YELLOW)e2e-docker$(NC)         - Run dedicated E2E Docker stack"
	@echo "  $(YELLOW)e2e-docker-vrt$(NC)     - Run visual regression tests in E2E Docker stack"
	@echo "  $(YELLOW)e2e-docker-vrt-update$(NC) - Update visual snapshots in E2E Docker stack"
	@echo "  $(YELLOW)e2e-docker-stop$(NC)    - Stop E2E docker services"
	@echo "  $(YELLOW)e2e-local$(NC)          - Run E2E tests against local services"
	@echo "  $(YELLOW)e2e-local-vrt$(NC)      - Run visual regression tests against local services"
	@echo "  $(YELLOW)e2e-local-vrt-update$(NC) - Update visual snapshots against local services"
	@echo "  $(YELLOW)install-tools$(NC)      - Install Firebase CLI and Air globally"
	@echo "  $(YELLOW)install-deps$(NC)       - Install all dependencies"
	@echo "  $(YELLOW)clean$(NC)              - Clean build artifacts"
	@echo "  $(YELLOW)help$(NC)               - Show this help message"

# Perform deep Docker cleanup
docker-clean:
	@echo "$(YELLOW)Performing deep Docker cleanup...$(NC)"
	@docker compose -f docker-compose.dev.yml down --rmi all --volumes --remove-orphans
	@docker system prune -af --volumes
	@echo "$(GREEN)Docker environment cleaned$(NC)"

status:
	@echo "$(YELLOW)Docker status:$(NC)"
	@docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
	@echo "\n Volume usage:"
	@docker volume ls
	@echo "\n Image usage:"
	@docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}"
