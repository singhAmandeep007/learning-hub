# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

.PHONY: dev-local prod-local stop-services install-tools install-deps clean help

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
	@cd backend && ENV_MODE=dev air -c .air.toml &
	$(call wait_for_port, 8080, Backend Server)
	@echo "$(YELLOW)Starting frontend development server...$(NC)"
	@cd frontend && npm run dev &
	@echo "$(GREEN)All services started!$(NC) Run $(RED)make stop-services$(NC) in a new terminal to stop all services."
	@echo "$(GREEN)Backend: http://localhost:8080$(NC)"
	@echo "$(GREEN)Frontend: http://localhost:3003$(NC)"
	@echo "$(GREEN)Firebase UI: http://localhost:4000$(NC)"
	@wait

# Production-like local environment
prod-local:
	@echo "$(GREEN)Starting production firebase in local environment...$(NC)"
	@echo "$(YELLOW)Starting backend server with Air...$(NC)"
	@cd backend && ENV_MODE=prod air -c .air.toml &
	$(call wait_for_port, 8080, Backend Server)
	@echo "$(YELLOW)Building and serving frontend...$(NC)"
	@cd frontend && npm run dev &
	@echo "$(GREEN)All services started! Press Ctrl+C to stop all services$(NC)"
	@echo "$(GREEN)Backend: http://localhost:8080$(NC)"
	@echo "$(GREEN)Frontend: http://localhost:3003$(NC)"
	@wait

# Stop all running services
# firebase emulator not exiting properly - https://github.com/firebase/firebase-tools/issues/3578#issuecomment-1024876668
stop-services:
	@echo "$(YELLOW)Stopping all services...$(NC)"
	@pkill -f "firebase emulators:start" || true
	@pkill -f "air -c" || true
	@pkill -f "npm run dev" || true
	@echo "$(GREEN)All services stopped$(NC)"

# Install global tools
install-tools:
	@echo "$(GREEN)Installing global tools...$(NC)"
	@echo "$(YELLOW)Installing Firebase CLI...$(NC)"
	@npm install -g firebase-tools
	@echo "$(YELLOW)Installing Air (Go hot reload)...$(NC)"
	@go install github.com/air-verse/air@latest
	@echo "$(GREEN)Global tools installed successfully!$(NC)"

# Install dependencies
install-deps:
	@echo "$(GREEN)Installing frontend dependencies...$(NC)"
	@cd frontend && npm install
	@echo "$(GREEN)Installing backend dependencies...$(NC)"
	@cd backend && go mod tidy

# Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning build artifacts...$(NC)"
	@cd frontend && rm -rf dist node_modules/.vite
	@cd backend && rm -rf tmp
	@echo "$(GREEN)Clean complete$(NC)"

# Help command
help:
	@echo "$(GREEN)Available commands:$(NC)"
	@echo "  $(YELLOW)dev-local$(NC)     - Start development environment (hot reload)"
	@echo "  $(YELLOW)prod-local$(NC)    - Start production-like local environment"
	@echo "  $(YELLOW)stop-services$(NC) - Stop all running services"
	@echo "  $(YELLOW)install-tools$(NC) - Install Firebase CLI and Air globally"
	@echo "  $(YELLOW)install-deps$(NC)  - Install all dependencies"
	@echo "  $(YELLOW)clean$(NC)         - Clean build artifacts"
	@echo "  $(YELLOW)help$(NC)          - Show this help message"