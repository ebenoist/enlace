package main

import (
	"log"
	"net/http"

	"github.com/dyatlov/go-opengraph/opengraph"
)

func ParseOG(url string) (*opengraph.OpenGraph, error) {
	og := opengraph.NewOpenGraph()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	log.Printf("og parse results - %+v", resp)

	err = og.ProcessHTML(resp.Body)
	log.Printf("og parse results - %s", og)
	return og, err
}
