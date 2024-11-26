package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/ebenoist/enlace/db"
	"github.com/ebenoist/enlace/env"
	"github.com/ebenoist/enlace/scrapper"
	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

type LinksRequest struct {
	URL      string `json:"url"`
	Category string `json:"category"`
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.LoadHTMLGlob("templates/*.tmpl")

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

			html, err := scrapper.Scrape(req.URL)
			if err != nil {
				log.Printf("could not parse HTML for %s - %s", req.URL, err)
			}

			link.Content = string(html)
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
		bookmarks, err := db.GetHoarderLinks()

		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("FATAL ERROR: %s", err),
			)
			return
		}

		links = append(links, bookmarks...)

		presented, err := presentRSS(userID, links)
		if err != nil {
			c.String(
				http.StatusInternalServerError,
				fmt.Sprintf("FATAL ERROR: %s", err),
			)
			return
		}

		c.Header("X-Content-Type-Options", "nosniff")
		c.Data(http.StatusOK, "text/xml; charset=utf-8", presented)
	})

	// TODO: consider a route that isn't susceptible to iteration
	// attacks
	r.GET("/links/:id", func(c *gin.Context) {
		raw := c.Param("id")
		source := c.Query("src")
		ext := filepath.Ext(raw)

		id := strings.TrimSuffix(filepath.Base(raw), ext)

		var link *db.Link
		var err error

		if source == "hoarder" {
			link, err = db.GetHoarderLink(id)
		} else {
			link, err = db.GetLink(id)
		}

		if err != nil {
			c.AbortWithError(
				http.StatusInternalServerError,
				err,
			)

			return
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"Title":     link.Title,
			"CreatedAt": link.CreatedAt,
			"Content":   link.Content,
			"UserID":    link.UserID,
		})
	})

	r.GET("/~:userID/rss.xsl", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/xsl; charset=utf-8", XSL)
	})

	r.Run(env.Get("LISTEN", "127.0.0.1:9292"))
}
