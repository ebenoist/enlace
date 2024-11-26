package db

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ebenoist/enlace/env"
)

type BookmarkContent struct {
	Title       string `json:"title"`
	HtmlContent string `json:"htmlContent"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

type Bookmark struct {
	ID        string          `json:"id"`
	CreatedAt *time.Time      `json:"createdAt"`
	Title     string          `json:"title"`
	Content   BookmarkContent `json:"content"`
}

type HoarderResponse struct {
	Bookmarks []Bookmark `json:"bookmarks"`
}

var client = http.Client{}
var hoarderHost = env.Get("HOARDER_HOST", "")
var hoarderKey = env.Get("HOARDER_API_KEY", "")

func GetHoarderLink(id string) (*Link, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/bookmarks/%s", hoarderHost, id),
		nil,
	)
	req.Header = http.Header{
		"Accept": {"application/json"},
		"Authorization": {
			fmt.Sprintf("Bearer %s", hoarderKey),
		},
	}

	if err != nil {
		log.Fatalf("HOARDER_HOST is invalid - %s", hoarderHost)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to call hoarder %s", err)
	}

	var parsed Bookmark
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		log.Printf("failed to parse - %s", err)
	}

	u, err := NewURL(parsed.Content.URL)
	return &Link{
		ID:          parsed.ID,
		URL:         u,
		UserID:      "hoarder-user",
		Category:    "hoarder",
		Title:       parsed.Content.Title,
		Description: parsed.Content.Description,
		CreatedAt:   parsed.CreatedAt,
		Content:     parsed.Content.HtmlContent,
	}, nil
}

func GetHoarderLinks() ([]*Link, error) {
	links := make([]*Link, 0)

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/v1/bookmarks", hoarderHost),
		nil,
	)
	if err != nil {
		log.Fatalf("HOARDER_HOST is invalid - %s", hoarderHost)
	}

	req.Header = http.Header{
		"Accept": {"application/json"},
		"Authorization": {
			fmt.Sprintf("Bearer %s", hoarderKey),
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to call hoarder %s", err)
	}

	var parsed HoarderResponse
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		log.Printf("failed to parse - %s", err)
	}

	for _, b := range parsed.Bookmarks {
		u, err := NewURL(b.Content.URL)
		if err != nil {
			log.Printf("invalid link %s", b.Content.URL)
		}

		links = append(links, &Link{
			ID:          b.ID,
			URL:         u,
			UserID:      "hoarder-user",
			Category:    "hoarder",
			Title:       b.Content.Title,
			Description: b.Content.Description,
			CreatedAt:   b.CreatedAt,
		})
	}

	return links, err
}
