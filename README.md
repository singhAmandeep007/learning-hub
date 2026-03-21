# Learning Hub

A comprehensive learning platform that provides resources in the form of videos, PDF and helpdesk articles. The platform features a public-facing landing page with searchable and filterable resource cards.

## Tech Stack

### Frontend
- React (TypeScript)
- Vite
- React Router
- React Query
- SCSS Modules

### Backend
- Go with Gin framework
- Firebase (Firestore & Storage) with emulator support

## Project Structure

```
learninghub/
├── frontend/          # React frontend application
├── backend/           # Go backend application
├── jenkinsfiles/      # Jenkins pipeline definitions
├── Makefile           # Makefile for common tasks
└── README.md
```


## Development Commands

### Frontend (React + Vite)
```bash
cd frontend
npm run dev        # Start development server (localhost:3000)
npm run build      # Build for production
npm run lint       # Run ESLint
npm run format     # Format code with Prettier
npm run preview    # Preview production build
```

### Backend (Go + Gin)
```bash
cd backend
go run main.go     # Start server directly
air -c .air.toml   # Hot reload development (requires Air)
go mod tidy        # Clean up dependencies
go test ./...      # Run tests
```


### Development Environment
```bash

# WITHOUT DOCKER
# Full development setup with Firebase emulators
make dev-local     # Starts Firebase emulators, backend (Air), and frontend

# Stop services
make stop-services # Stop all running services

# WITH DOCKER
# Docker development
make docker-dev    # Build and run with Docker Compose

# Stop services
make docker-dev-stop # Stop all dev docker services
```

## Architecture Overview

### Multi-Product Learning Hub
The application is designed to manage learning resources across different products (currently "ecomm"). All data is product-scoped in separate Firestore collections.

### Backend Architecture (Go)
- **Framework**: Gin web framework
- **Database**: Firestore (with emulator support for development)
- **Storage**: Firebase Storage for file uploads
- **Structure**:
  - `main.go`: Entry point, router setup
  - `config/`: Environment configuration management
  - `handlers/`: HTTP request handlers (resources, tags)
  - `models/`: Data models (Resource, Tag, Response types)
  - `middleware/`: CORS, rate limiting, product validation
  - `utils/`: File upload/deletion, tag management utilities
  - `firebase/`: Firebase initialization and client management

### Frontend Architecture (React)
- **Framework**: React 19 with TypeScript
- **Build Tool**: Vite
- **Routing**: React Router v7
- **State Management**: TanStack Query for server state
- **Styling**: SCSS modules
- **Key Components**:
  - Product-scoped routing (`/:product/resources`)
  - Resource management (CRUD operations)
  - File upload handling (videos, PDFs, thumbnails)
  - Rich text editing with TiptapEditor
  - Tag-based filtering and search

### Key Features
- **Product-specific collections**: Each product has its own resource collection
- **File upload management**: Handles video, PDF, and thumbnail uploads
- **Rich text editing**: TiptapEditor integration for article content
- **Tag-based organization**: Dynamic tagging with usage counts
- **Search and filtering**: By type, tags, and text search
- **Pagination**: Cursor-based pagination for large datasets

## Environment Configuration

### Required Environment Variables

#### Backend (Go)
```bash
ENV_MODE=dev                    # "dev" or "prod"
PORT=8000                       # Server port
```

**Authentication Methods:**
- **Development**: Firebase emulators (no authentication required)
- **Production**: GCP-native authentication

#### Frontend (React)
```bash
VITE_API_BASE_URL=/api/v1       # API base URL
VITE_PORT=3000                  # Dev server port
VITE_PROXY_API_HOST=http://localhost:8000  # Backend proxy for dev
```

## Database Structure

### Firestore Collections
- `ecomm`: Resources for ECOMM product
- `ecomm`: Tag usage counts for ECOMM product


## API Endpoints

All endpoints are product-scoped: `/api/v1/:product/`

### Resources
- `GET /:product/resources` - List resources with filtering and pagination
- `GET /:product/resources/:id` - Get single resource
- `POST /:product/resources` - Create resource (multipart/form-data)
- `PATCH /:product/resources/:id` - Update resource (multipart/form-data)
- `DELETE /:product/resources/:id` - Delete resource

### Tags
- `GET /:product/tags` - Get all tags with usage counts

## Testing

### Backend Tests
```bash
cd backend
go test ./handlers -v      # Test handlers
go test ./utils -v         # Test utilities
go test ./... -v           # Run all tests
```

### Frontend Development
- Mock Service Worker (MSW) is configured for API mocking
- React Query DevTools available in development mode

## Development Tips

1. **Product Validation**: All routes require valid product parameter ("ecomm")
2. **File Uploads**: Use multipart/form-data for creating/updating resources with files
3. **Firebase Emulators**: Use `make dev-local` to run with Firebase emulators for offline development
4. **Hot Reload**: Backend uses Air for hot reload, frontend uses Vite HMR
5. **CORS**: Configured for local development across different ports
6. **Rate Limiting**: 100 requests per minute per IP in backend
7. **Upload Size**: Max 500MB per file