package main

import (
	"cat"
	"log"
	"net/http"
	"time"
)

func onlyForV2() cat.HandlerFunc {
	return func(c *cat.Context) {
		//Start timer
		t := time.Now()
		//if server error occurred
		c.Fail(500, "Internal Server Error")

		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))

	}
}

func main() {
	r := cat.New()
	r.Use(cat.Logger()) //global middlerware
	r.GET("/", func(c *cat.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Ngo</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) //v2 group middlerware
	{
		v2.GET("/hello/:name", func(c *cat.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":8000")
}
