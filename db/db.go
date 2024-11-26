package db

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"

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
var idgen, _ = sqids.New()

var ErrNotFound = errors.New("entity not found")
var ErrBadRequest = errors.New("bad id")

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
	fmt.Printf("got to scan %+v", src)
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
	Content     string     `db:"content"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	ID          string     `db:"id"`
	IP          netip.Addr `db:"ip"`
}

func (l *Link) GUID() string {
	if l.Category == "hoarder" {
		return fmt.Sprintf(
			"%s/links/%s?src=hoarder",
			env.Get("HOST"),
			l.ID,
		)
	}

	return fmt.Sprintf(
		"%s/links/%s",
		env.Get("HOST"),
		l.UID(),
	)
}

func (l *Link) UID() string {
	n, _ := strconv.Atoi(l.ID)
	publicID, _ := idgen.Encode([]uint64{uint64(n), 42})
	return publicID
}

func init() {
	url := env.Get("DATABASE_URL", "./enlace.db")
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
			updated_at DATETIME NULL,
			url VARCHAR(2048),
			user_id VARCHAR(128),
			category VARCHAR(128) NULL,
			description VARCHAR(2048) NULL,
			content TEXT NULL,
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
		`SELECT
			id,
			title,
			url,
			description,
			created_at,
			updated_at
		FROM links WHERE user_id = ?
		ORDER BY created_at DESC`,
		userID,
	)

	return links, err
}

func GetLink(id string) (*Link, error) {
	rid := idgen.Decode(id)
	if len(rid) == 0 {
		return nil, ErrBadRequest
	}

	log.Println(rid)

	var link Link
	err := conn.Get(&link, `SELECT * FROM links WHERE id = ?`, rid[0])
	return &link, err
}

func UpdateLink(link *Link) (*Link, error) {
	now := time.Now()
	_, err := conn.Exec(`
		UPDATE links SET
			title = ?,
			description = ?,
			content = ?,
			updated_at = ?
		WHERE id = ?
	`,
		link.Title,
		link.Description,
		link.Content,
		now.Format(time.RFC3339),
		link.ID,
	)

	link.UpdatedAt = &now
	return link, err
}

func CreateLink(link *Link) (*Link, error) {
	now := time.Now()
	res, err := conn.Exec(`
		INSERT INTO links (
			created_at,
			url,
			user_id,
			category,
			description,
			title,
			content
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		now.Format(time.RFC3339),
		link.URL,
		link.UserID,
		link.Category,
		"",
		"",
		"",
	)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	link.ID = strconv.Itoa(int(id))
	link.CreatedAt = &now

	return link, err
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
