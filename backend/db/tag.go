package db

import (
	"context"

	"cloud.google.com/go/firestore"

	"learninghub/constants"
	"learninghub/models"
)

type TagService struct {
	db *DB
}

func NewTagService(db *DB) *TagService {
	return &TagService{db: db}
}

// List retrieves all tags ordered by usage count
func (ts *TagService) List(ctx context.Context, product string) ([]*firestore.DocumentSnapshot, error) {
	collectionName := constants.GetTagsCollectionName(product)
	return ts.db.client.Collection(collectionName).OrderBy("usageCount", firestore.Desc).Documents(ctx).GetAll()
}

// UpdateUsage updates the usage count for tags in a product-specific collection
func (ts *TagService) UpdateUsage(ctx context.Context, product string, tags []string, delta int) error {
	for _, tag := range tags {
		if tag == "" {
			continue
		}

		if err := ts.updateSingleTagUsage(ctx, product, tag, delta); err != nil {
			return err
		}
	}
	return nil
}

// updateSingleTagUsage updates the usage count for a single tag using a transaction
func (ts *TagService) updateSingleTagUsage(ctx context.Context, product, tag string, delta int) error {
	collectionName := constants.GetTagsCollectionName(product)
	tagRef := ts.db.client.Collection(collectionName).Doc(tag)

	// Use a transaction to ensure atomicity
	return ts.db.RunTransaction(ctx, func(_ context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(tagRef)
		if err != nil {
			// Tag doesn't exist, create it
			return tx.Set(tagRef, models.Tag{
				Name:       tag,
				UsageCount: max(0, delta),
			})
		}

		var existingTag models.Tag
		if err := doc.DataTo(&existingTag); err != nil {
			return err
		}

		newCount := max(0, existingTag.UsageCount+delta)
		if newCount == 0 {
			// Delete tag if usage count reaches 0
			return tx.Delete(tagRef)
		}

		return tx.Update(tagRef, []firestore.Update{
			{Path: "usageCount", Value: newCount},
		})
	})
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
