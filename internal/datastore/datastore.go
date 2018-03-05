package datastore

import (
	"context"
	"database/sql"
	"time"

	"github.com/ua-parser/uap-go/uaparser"
)

// Datastore is the exported interface for our datastore
type Datastore interface {
	// URLs
	GetURLBySlug(ctx context.Context, slug string) (*URLMap, error)
	SaveNewURL(ctx context.Context, slug string, url string) error

	// Tracking
	TrackHit(ctx context.Context, urlID int, client *uaparser.Client, remoteAddress string) error

	// Stats
	GetVisitCounts(ctx context.Context) ([]VisitCount, error)
	GetVisits(ctx context.Context, slug string) (*URLMap, []Visit, error)
}

// URLMap models our basic short url to long url relationship, or the url table
type URLMap struct {
	ID   int
	Slug string
	URL  string
}

// Visit models a single visit record for a short url
type Visit struct {
	ID      int
	Device  string
	OS      string
	Browser string
	IP      string
	Time    time.Time
}

// VisitCount models aggregate visit data for a short url
type VisitCount struct {
	Slug  string
	URL   string
	Count int
}

type datastore struct {
	db *sql.DB
}

// New creates a new Datastore, given a database connection
func New(db *sql.DB) (Datastore, error) {
	ds := datastore{
		db: db,
	}
	return ds, nil
}

func (ds datastore) GetURLBySlug(ctx context.Context, slug string) (*URLMap, error) {
	var url URLMap

	row := ds.db.QueryRow(`SELECT id, slug, url FROM url WHERE slug = ?`, slug)
	err := row.Scan(&url.ID, &url.Slug, &url.URL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &url, nil
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

func (ds datastore) GetVisitCounts(ctx context.Context) ([]VisitCount, error) {
	vc := make([]VisitCount, 0)

	const query = `SELECT slug, url, COALESCE(cnt, 0) FROM url u LEFT JOIN ( SELECT url_id, COUNT(*) cnt FROM visit GROUP BY url_id ) v ON v.url_id = u.id ORDER BY id`
	rows, err := ds.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var v VisitCount
		err = rows.Scan(&v.Slug, &v.URL, &v.Count)
		if err != nil {
			return nil, err
		}

		vc = append(vc, v)
	}

	return vc, nil
}

func (ds datastore) GetVisits(ctx context.Context, slug string) (*URLMap, []Visit, error) {
	url, err := ds.GetURLBySlug(ctx, slug)
	if err != nil {
		return nil, nil, err
	}

	visits := make([]Visit, 0)
	const query = `SELECT id, device, os, browser, ip, created_at FROM visit WHERE url_id = ? ORDER BY id DESC`

	rows, err := ds.db.Query(query, url.ID)
	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var v Visit
		err = rows.Scan(&v.ID, &v.Device, &v.OS, &v.Browser, &v.IP, &v.Time)
		if err != nil {
			return nil, nil, err
		}

		visits = append(visits, v)
	}

	return url, visits, nil
}
