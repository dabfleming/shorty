package datastore

import (
	"context"
	"database/sql"
)

type Datastore interface {
	// URLs
	GetURLBySlug(ctx context.Context, slug string) (string, error)
}

type datastore struct {
	db *sql.DB
}

func New(db *sql.DB) (Datastore, error) {
	ds := datastore{
		db: db,
	}
	return ds, nil
}

func (ds datastore) GetURLBySlug(ctx context.Context, slug string) (string, error) {
	return "", nil
}
