package datastore

import (
	"context"
	"database/sql"
)

type Datastore interface {
	// URLs
	GetURLBySlug(ctx context.Context, slug string) (string, error)
	SaveNewURL(ctx context.Context, slug string, url string) error
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
	var (
		id  int
		url string
	)
	row := ds.db.QueryRow(`SELECT id, url FROM url WHERE slug = ?`, slug)
	err := row.Scan(&id, &url)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return url, nil
}

func (ds datastore) SaveNewURL(ctx context.Context, slug string, url string) error {
	_, err := ds.db.Exec(`INSERT INTO url (slug, url) VALUES (?, ?)`, slug, url)
	return err
}
