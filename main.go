package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}

func main() {
  r := gin.Default()
	r.SetTrustedProxies(nil)

  r.POST("/links", func(c *gin.Context) {
		c.String(http.StatusOK, "success")
  })

	r.GET("/~:id/rss.xml", func(c *gin.Context) {
		id := c.Param("id")

		c.String(http.StatusOK, fmt.Sprintf(`
			<?xml version="1.0" encoding="UTF-8"?>
			<rss xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
				<channel>
					<title>%s</title>
					<link>https://enlace.space/~erik/rss.xml</link>
					<item>
						<title>og:title</title>
						<description>og:description</description>
						<link>https://foo.com/bar</link>
						<pubDate>Sat, 07 Nov 2024 18:34:56 +0000</pubDate>
						<category>fizzle</category>
					</item>
				</channel>
			</rss>`, id));

	})

	r.Run(env("LISTEN", "127.0.0.1:9292"))
}
