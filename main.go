package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ebenoist/enlace/db"
	"github.com/ebenoist/enlace/env"
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

		c.String(http.StatusOK, string(presented))
	})

	r.Run(env.Get("LISTEN", "127.0.0.1:9292"))
}
