package db

import (
	"context"
	"strconv"

	"learninghub/constants"
	"learninghub/models"

	"cloud.google.com/go/firestore"
)

// ResourceQuery represents query parameters for listing resources
type ResourceQuery struct {
	Product string
	Type    string
	Tags    []string
	Search  string
	Cursor  string
	Limit   int
}

// ResourceService handles resource database operations
type ResourceService struct {
	db *DB
}

// NewResourceService creates a new resource service
func NewResourceService(db *DB) *ResourceService {
	return &ResourceService{db: db}
}

// List retrieves resources with filtering and pagination
func (rs *ResourceService) List(ctx context.Context, query ResourceQuery) ([]*firestore.DocumentSnapshot, error) {
	collectionName := constants.GetResourcesCollectionName(query.Product)
	firestoreQuery := rs.db.client.Collection(collectionName).OrderBy("createdAt", firestore.Desc)

	// Apply type filter
	if query.Type != "" {
		firestoreQuery = firestoreQuery.Where("type", "==", query.Type)
	}

	// Apply tags filter
	if len(query.Tags) > 0 {
		firestoreQuery = firestoreQuery.Where("tags", "array-contains-any", query.Tags)
	}

	// Apply cursor for pagination
	if query.Cursor != "" {
		if offset, err := strconv.Atoi(query.Cursor); err == nil && offset >= 0 {
			firestoreQuery = firestoreQuery.Offset(offset)
		}
	}

	// Execute query with limit
	return firestoreQuery.Limit(query.Limit).Documents(ctx).GetAll()
}

// GetByID retrieves a single resource by ID
func (rs *ResourceService) GetByID(ctx context.Context, product, id string) (*firestore.DocumentSnapshot, error) {
	collectionName := constants.GetResourcesCollectionName(product)
	return rs.db.client.Collection(collectionName).Doc(id).Get(ctx)
}

// Create creates a new resource
func (rs *ResourceService) Create(ctx context.Context, product string, resource models.Resource) (*firestore.DocumentRef, error) {
	collectionName := constants.GetResourcesCollectionName(product)
	docRef, _, err := rs.db.client.Collection(collectionName).Add(ctx, resource)
	return docRef, err
}

// Update updates an existing resource
func (rs *ResourceService) Update(ctx context.Context, product, id string, resource models.Resource) error {
	collectionName := constants.GetResourcesCollectionName(product)
	_, err := rs.db.client.Collection(collectionName).Doc(id).Set(ctx, resource)
	return err
}

// Delete deletes a resource by ID
func (rs *ResourceService) Delete(ctx context.Context, product, id string) error {
	collectionName := constants.GetResourcesCollectionName(product)
	_, err := rs.db.client.Collection(collectionName).Doc(id).Delete(ctx)
	return err
}
