# Learning Hub API Contract

## Overview

This document outlines the API contract for the Learning Hub application. The API provides endpoints for retrieving resources, tags, and related data for the learning hub.

## Base URL

```
https://example.com/api
```

## Admin Authentication

Creating, updating and deleting resources API requests require admin authentication using secret:

```
AdminSecret: <secret>
```

## Error Handling

All endpoints follow a consistent error response format:

```json
{
  "error": "string",
  "message": "string", // Optional additional information
}
```

Common error codes:
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Endpoints

### Resources

#### Get Resources

Retrieves a list of resources with filtering, searching, and pagination.

```
GET /resources
```

**Query Parameters:**

| Parameter | Type   | Required | Description                       |
|-----------|--------|----------|-----------------------------------|
| search    | string | No       | Search term for resource title and description
| type      | string | No       | Filter by resource type: 'video' or 'pdf' or 'article' or 'all'
| tags      | string | No       | Comma separated tags
| cursor    | number | No       | No. of items skipped  
| limit     | number | No       | Number of items per page (default: 20, max: 100) 

**Response:**

```json
{
  "data": [
    {
      "id": "string",
      "title": "string",
      "description": "string",
      "type": "video" | "article" | "pdf",
      "url": "string",
      "thumbnailUrl": "string", // Optional
      "tags": ["string"],
      "createdAt": "string",
      "updatedAt": "string",
    }
  ],
  "hasMore": boolean
}
```

**Status Codes:**
- `200` - Success
- `500` - Internal Server Error

#### Get Resource by ID

Retrieves a specific resource by its ID.

```
GET /resources/{id}
```

**URL Parameters:**

| Parameter | Type   | Required | Description     |
|-----------|--------|----------|-----------------|
| id        | string | Yes      | Resource ID     |

**Response:**

```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "type": "video" | "article" | "pdf",
  "url": "string",
  "thumbnailUrl": "string", // Optional
  "tags": ["string"],
  "createdAt": "string",
  "updatedAt": "string",
}
```

**Status Codes:**
- `200` - Success
- `404` - Resource not found
- `500` - Internal Server Error

#### Create Resource

Creates a new resource.

```
POST /resources
```

**Request Body:**

```json
{
  "title": "string",
  "description": "string",
  "type": "video" | "article" | "pdf",
  "url": "string", // Optional
  "thumbnailUrl": "string", // Optional
  "tags": ["string"],
}
```

```multipart/form-data
	file: File,
	thumbnail: File
```

**Response:**

```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "type": "video" | "article" | "pdf",
  "url": "string",
  "thumbnailUrl": "string", // Optional
  "tags": ["string"],
  "createdAt": "string",
  "updatedAt": "string",
}
```

**Status Codes:**
- `201` - Created
- `400` - Invalid request data
- `401` - Unauthorized
- `500` - Internal Server Error

#### Update Resource

Updates an existing resource.

```
PATCH /resources/{id}
```

**URL Parameters:**

| Parameter | Type   | Required | Description     |
|-----------|--------|----------|-----------------|
| id        | string | Yes      | Resource ID     |

**Request Body:**

```json
{
  "title": "string",
  "description": "string",
  "type": "video" | "article" | "pdf",
  "url": "string",
  "thumbnailUrl": "string", // Optional
  "tags": ["string"]
}
```

```multipart/form-data
	file: File,
	thumbnail: File
```

**Response:**

```json
{
  "id": "string",
  "title": "string",
  "description": "string",
  "type": "video" | "article" | "pdf",
  "url": "string",
  "thumbnailUrl": "string", // Optional
  "tags": ["string"],
  "createdAt": "string",
  "updatedAt": "string",
}
```

**Status Codes:**
- `200` - Success
- `400` - Invalid request data
- `404` - Resource not found
- `401` - Unauthorized
- `500` - Internal Server Error

#### Delete Resource

Deletes a resource.

```
DELETE /resources/{id}
```

**URL Parameters:**

| Parameter | Type   | Required | Description     |
|-----------|--------|----------|-----------------|
| id        | string | Yes      | Resource ID     |

**Response:**

```json
{
  "message": "string"
}
```

**Status Codes:**
- `200` - Success
- `404` - Resource not found
- `401` - Unauthorized
- `500` - Internal Server Error

### Tags

#### Get All Tags

Retrieves all available Tags.

```
GET /tags
```

**Response:**

```json
[
	{
		"name": "string",
		"usageCount": 0 // Number of resources with this tag
	}
]
```

**Status Codes:**
- `200` - Success

## Data Models

### Resource
```typescript
{
  id: string;
  title: string;
  description: string;
  type: 'video' | 'article' | 'pdf';
  url: string;
  thumbnailUrl?: string;
  tags: string[];
  createdAt: string;
  updatedAt: string;
}
```

### Tag
```typescript
{
  name: string;
  usageCount: number; // Number of resources with this tag
}
```