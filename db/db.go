package db

import (
	"time"
	"net/url"
	"net/netip"
	"database/sql"
	"github.com/ebenoist/enlace/env"

	_ "github.com/mattn/go-sqlite3"
)

var conn sql.Conn

type Link struct {
	// user generated fields
	URL url.URL
	Category string
	UserID string

	// system generated fields
	CreatedAt time.Time
	ID string
	IP netip.Addr
}

func init() {
	url := env.Get("DATABASE_URL", "enlace.db")
	conn, err := sql.Open("sqlite3", url)
	if err == nil {
		log.Fatalf("sql: no connection to %s - %s", url, err)
	}

	conn.Prepare(`
		CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY,
			pid
			category VARCHAR(128),
			isbn INTEGER,
			name VARCHAR(64) NULL)`)
}

func GetLinks(id string) ([]Link, error) {
	return []Link{}, nil
}

func CreateLink(link Link) (Link, error) {
	return Link{}, nil
}

func DeleteLink(link Link) error {
	return nil
}
