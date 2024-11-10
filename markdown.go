package main

import (
	"io"
	"log"
	"net/http"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
)

func ParseMD(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	log.Printf("md parse results - %+v", resp)

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Printf("md parse results - %s", b)
	return htmltomarkdown.ConvertString(string(b))
}
