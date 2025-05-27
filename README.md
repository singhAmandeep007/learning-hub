# Learning Hub

A comprehensive learning platform that provides resources in the form of videos and PDF helpdesk articles. The platform features a public-facing landing page with searchable and filterable resource cards, and a private admin section for managing these resources.

## Tech Stack

### Frontend
- React (TypeScript)
- Vite
- React Router
- React Hook Form
- React Query
- SCSS Modules

### Backend
- Go with Gin framework
- Firebase (Firestore & Storage)

## Project Structure

```
learning-hub/
├── frontend/           # React frontend application
├── backend/           # Go backend application
└── README.md
```

## Getting Started

### Prerequisites
- Node.js (v21 or higher)
- Go (v1.24 or higher)
- Firebase account and project

### Frontend Setup
1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the development server:
   ```bash
   npm run dev
   ```

### Backend Setup
1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Start the server:
   ```bash
   go run main.go
   ```

## License

MIT 