package main

import (
	"fmt"
	"net/http"

	"github.com/ebenoist/enlace/db"
	"github.com/ebenoist/enlace/env"
	"github.com/gin-gonic/gin"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)

	r.POST("/links", func(c *gin.Context) {
		db.CreateLink(&db.Link{})
		c.String(http.StatusOK, "success")
	})

	r.GET("/~:userID/rss.xml", func(c *gin.Context) {
		userID := c.Param("userID")
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
