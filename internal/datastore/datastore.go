package datastore

import (
	"context"
	"database/sql"

	"github.com/ua-parser/uap-go/uaparser"
)

type Datastore interface {
	// URLs
	GetURLBySlug(ctx context.Context, slug string) (int, string, error)
	SaveNewURL(ctx context.Context, slug string, url string) error

	// Tracking
	TrackHit(ctx context.Context, urlID int, client *uaparser.Client, remoteAddress string) error
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

func (ds datastore) GetURLBySlug(ctx context.Context, slug string) (int, string, error) {
	var (
		id  int
		url string
	)
	row := ds.db.QueryRow(`SELECT id, url FROM url WHERE slug = ?`, slug)
	err := row.Scan(&id, &url)
	if err == sql.ErrNoRows {
		return 0, "", nil
	}
	if err != nil {
		return 0, "", err
	}

	return id, url, nil
}

func (ds datastore) SaveNewURL(ctx context.Context, slug string, url string) error {
	_, err := ds.db.Exec(`INSERT INTO url (slug, url) VALUES (?, ?)`, slug, url)
	return err
}

func (ds datastore) TrackHit(ctx context.Context, urlID int, client *uaparser.Client, remoteAddress string) error {
	const query = `INSERT INTO visit (url_id, device, os, browser, ip) VALUES (?, ?, ?, ?, ?)`
	_, err := ds.db.Exec(query, urlID, client.Device.Family, client.Os.Family, client.UserAgent.Family, remoteAddress)
	return err
}
