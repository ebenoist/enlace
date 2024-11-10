package main

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/ebenoist/enlace/db"
	"github.com/ebenoist/enlace/env"
)

type rss struct {
	Version string    `xml:"version,attr"`
	Schema  string    `xml:"xmlns:atom,attr"`
	Channel []Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Item  []Item `xml:"item"`
}

type Item struct {
	Title       string     `xml:"title"`
	Description string     `xml:"description"`
	Link        string     `xml:"link"`
	PubDate     *time.Time `xml:"pubDate"`
	Category    string     `xml:"category"`
	Guid        string     `xml:"guid"`
}

func presentRSS(userID string, links []*db.Link) ([]byte, error) {
	rss := &rss{
		Version: "2.0",
		Schema:  "http://www.w3.org/2005/Atom",
		Channel: []Channel{
			{
				Title: userID,
				Link:  presentLink(userID),
				Item:  presentItems(links),
			},
		},
	}

	return xml.MarshalIndent(rss, "", "  ")
}

func presentItems(links []*db.Link) []Item {
	items := make([]Item, 0, len(links))

	for _, link := range links {
		items = append(items, Item{
			Title:       link.Title,
			Description: link.Description,
			Link:        link.URL.String(),
			PubDate:     link.CreatedAt,
			Category:    link.Category,
		})
	}

	return items
}

func presentLink(userID string) string {
	return fmt.Sprintf(
		"%s/~%s/rss.xml",
		env.Get("HOST", "http://127.0.0.1:9292"),
		userID,
	)
}
