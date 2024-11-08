package db

import (
	"database/sql/driver"

	"github.com/ebenoist/enlace/env"
	"github.com/jmoiron/sqlx"
	"github.com/sqids/sqids-go"

	"log"
	"net/netip"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var conn *sqlx.DB
var s, _ = sqids.New()

func NewURL(r string) (*URL, error) {
	u, err := url.Parse(r)
	if err != nil {
		return nil, err
	}

	return &URL{URL: u}, nil
}

type URL struct {
	*url.URL
}

func (u *URL) Value() (driver.Value, error) {
	return u.String(), nil
}

func (u *URL) Scan(src interface{}) error {
	p, err := url.Parse(src.(string))
	if err != nil {
		return err
	}

	u.URL = p

	return err
}

type Link struct {
	// user generated fields
	URL      *URL   `db:"url"`
	Category string `db:"category"`
	UserID   string `db:"user_id"`

	// system generated fields
	Title       string     `db:"title"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	ID          string     `db:"id"`
	IP          netip.Addr `db:"ip"`
}

func init() {
	url := env.Get("DATABASE_URL", "../enlace.db")
	var err error

	conn, err = sqlx.Open("sqlite3", url)
	if err != nil {
		log.Fatalf("sql: no connection to %s - %s", url, err)
	}

	_, err = conn.Exec(`
		BEGIN;
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY,
			created_at DATETIME,
			url VARCHAR(2048),
			user_id VARCHAR(128),
			category VARCHAR(128) NULL,
			description VARCHAR(2048) NULL,
			title VARCHAR(50) NULL
		);

		CREATE INDEX IF NOT EXISTS links_user_id_idx ON links (user_id);
		CREATE INDEX IF NOT EXISTS links_category_idx ON links (category);

		COMMIT;
	`)
	if err != nil {
		log.Fatalf("sql: could not create tables - %s", err)
	}
}

func GetLinks(userID string) ([]*Link, error) {
	links := make([]*Link, 0)
	err := conn.Select(
		&links,
		`SELECT * FROM links WHERE user_id = ?`,
		userID,
	)

	return links, err
}

func CreateLink(link *Link) (*Link, error) {
	_, err := conn.Exec(`
		INSERT INTO links (
			created_at,
			url,
			user_id,
			category,
			description,
			title
		) VALUES (?, ?, ?, ?, ?, ?)`,
		time.Now().Format(time.RFC3339),
		link.URL,
		link.UserID,
		link.Category,
		link.Description,
		link.Title,
	)
	if err != nil {
		return nil, err
	}

	return &Link{}, err
}

func DeleteLink(link *Link) error {
	return nil
}

func purge() {
	_, err := conn.Exec("DELETE FROM links")
	if err != nil {
		panic(err)
	}
}
