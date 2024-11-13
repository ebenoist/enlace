package main

import (
	"log"
	"net/http"

	"github.com/ebenoist/enlace/htmlinfo"
)

func ParseOG(url string) (*htmlinfo.HTMLInfo, error) {
	info := htmlinfo.NewHTMLInfo()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	log.Printf("og parse results - %+v", resp)

	err = info.Parse(resp.Body, &url, nil)
	log.Printf("og parse results - %s", info)
	return info, err
}
