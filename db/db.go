package db

import (
	"github.com/ebenoist/enlace/env"
	"github.com/jmoiron/sqlx"
	"github.com/sqids/sqids-go"

	// "github.com/jmoiron/sqlx"
	"log"
	"net/netip"
	"net/url"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var conn sqlx.DB
var s, _ = sqids.New()

type Link struct {
	// user generated fields
	URL      *url.URL `db:"url"`
	Category string   `db:"category"`
	UserID   string   `db:"user_id"`

	// system generated fields
	CreatedAt time.Time  `db:"created_at"`
	ID        string     `db:"id"`
	IP        netip.Addr `db:"ip"`
}

func init() {
	url := env.Get("DATABASE_URL", "enlace.db")
	conn, err := sqlx.Open("sqlite3", url)
	if err != nil {
		log.Fatalf("sql: no connection to %s - %s", url, err)
	}

	_, err = conn.Exec(`
		BEGIN;
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY,
			created_at VARCHAR(24),
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

func GetLinks(userID string) ([]Link, error) {
	links := []Link{}

	err := conn.Select(&links, ` // TODO: not working
	SELECT *
	FROM links
	WHERE links.user_id = $1
	ORDER BY created_at DESC")
	`, userID)

	return links, err
}

func CreateLink(link Link) (Link, error) {
	return Link{}, nil
}

func DeleteLink(link Link) error {
	return nil
}
