package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ebenoist/enlace/db"
	"github.com/ebenoist/enlace/env"
	"github.com/gin-gonic/gin"

	"github.com/yuin/goldmark"

	_ "github.com/mattn/go-sqlite3"
)

type LinksRequest struct {
	URL      string `json:"url"`
	Category string `json:"category"`
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.POST("/links", func(c *gin.Context) {
		var req LinksRequest
		err := c.BindJSON(&req)
		user, _, _ := c.Request.BasicAuth()

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("request was malformed"))
			return
		}

		u, err := db.NewURL(req.URL)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("url was malformed"))
			return
		}

		link, err := db.CreateLink(&db.Link{
			UserID:   user,
			URL:      u,
			Category: req.Category,
		})
		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				fmt.Errorf("FATAL ERROR: %s", err),
			)
			return
		}

		// Needs to be extracted to an after save routine of some sort.
		// Can this be something pluggable? Some event pattern?
		go func() {
			og, err := ParseOG(req.URL)
			if err != nil {
				log.Printf("could not parse OG for %s - %s", req.URL, err)
			}

			md, err := ParseMD(req.URL)
			if err != nil {
				log.Printf("could not parse MD for %s - %s", req.URL, err)
			}

			link.Markdown = md
			link.Description = og.Description
			link.Title = og.Title

			updated, err := db.UpdateLink(link)
			if err != nil {
				log.Printf("could not update OG for %s - %s", req.URL, err)
				return
			}

			log.Printf("updated OG for %+v", updated)
		}()

		c.String(http.StatusOK, "success")
	})

	r.GET("/~:userID/rss.xml", func(c *gin.Context) {
		userID := c.Param("userID")
		log.Printf("returning rss feed for %s", userID)
		links, err := db.GetLinks(userID)

		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("FATAL ERROR: %s", err),
			)
			return
		}

		presented, err := presentRSS(userID, links)
		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("FATAL ERROR: %s", err),
			)
			return
		}

		c.Data(http.StatusOK, "application/rss+xml", presented)
	})

	// TODO: consider a route that isn't susceptible to iteration
	// attacks

	// just return markdown for `md`
	r.GET("/links/:id", func(c *gin.Context) {
		raw := c.Param("id")
		ext := filepath.Ext(raw)

		id := strings.TrimSuffix(filepath.Base(raw), ext)
		link, err := db.GetLink(id)

		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				err,
			)

			return
		}

		if ext == ".html" {
			var buf bytes.Buffer
			goldmark.Convert([]byte(link.Markdown), &buf)

			c.Data(http.StatusOK, "text/html; charset=utf-8", buf.Bytes())
			return
		}

		c.String(http.StatusOK, link.Markdown)
	})

	r.Run(env.Get("LISTEN", "127.0.0.1:9292"))
}
