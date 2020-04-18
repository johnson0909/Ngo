package main

import (
	"net/http"
	"cat"
)

func main() {
	r := cat.New()
	r.GET("/", func(c *cat.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Ngo<h1>")
	})

	r.GET("/hello", func(c *cat.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *cat.Context) {
		c.JSON(http.StatusOK, cat.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8000")
}