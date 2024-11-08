package db

import (
	"net/url"
	"testing"
)

func newURL(href string) *URL {
	u, err := url.Parse(href)
	if err != nil {
		panic(err)
	}

	return &URL{u}
}

func Test_CreateLink(t *testing.T) {
	_, err := CreateLink(&Link{
		UserID: "52",
		URL:    newURL("https://google.com"),
	})
	if err != nil {
		t.Errorf("could not create link %s", err)
	}
}

func Test_GetLinks(t *testing.T) {
	purge()

	_, err := CreateLink(&Link{
		UserID: "123",
		URL:    newURL("https://foo.com"),
	})
	_, err = CreateLink(&Link{
		UserID: "42",
		URL:    newURL("https://bar.com"),
	})
	if err != nil {
		t.Errorf("could not create %s", err)
	}

	links, err := GetLinks("42")
	if err != nil {
		t.Errorf("could not query %s", err)
	}

	if len(links) != 1 {
		t.Errorf("got %d links, but expected 1", len(links))
	}
}
