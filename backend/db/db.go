package db

import (
	"context"

	"cloud.google.com/go/firestore"

	"learninghub/firebase"
)

type DB struct {
	client *firestore.Client
}

func New() *DB {
	return &DB{
		client: firebase.FirestoreClient,
	}
}

func (db *DB) Client() *firestore.Client {
	return db.client
}

func (db *DB) RunTransaction(ctx context.Context, fn func(context.Context, *firestore.Transaction) error) error {
	return db.client.RunTransaction(ctx, fn)
}
