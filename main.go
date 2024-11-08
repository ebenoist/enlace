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
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("request was malformed"))
			return
		}

		u, err := db.NewURL(req.URL)
		if err != nil {
			c.AbortWithError(http.StatusBadRequest, errors.New("url was malformed"))
			return
		}

		db.CreateLink(&db.Link{
			URL:      u,
			Category: req.Category,
		})
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
