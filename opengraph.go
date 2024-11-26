package main

import (
	"net/http"

	"github.com/ebenoist/enlace/htmlinfo"
)

func ParseOG(url string) (*htmlinfo.HTMLInfo, error) {
	info := htmlinfo.NewHTMLInfo()

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	err = info.Parse(resp.Body, &url, nil)
	return info, err
}
