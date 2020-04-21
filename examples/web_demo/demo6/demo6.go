package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"cat"
)

type student struct {
	Name string
	Age int8
}

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := cat.New()
	r.Use(cat.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate": formatAsDate,
	})
	r.LoadHTMLGlob("tempaltes/*")
	r.Static("/assets", "./static")

	stu1 := &student{Name: "saulliu", Age: 22}
	stu2 := &student{Name: "chanceyin", Age: 30}
	r.GET("/", func(c *cat.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	r.GET("/students", func(c *cat.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", cat.H{
			"title": "cat",
			"stuArr": [2]*student{stu1, stu2},
		})
	})
	
	r.GET("/date", func(c *cat.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", cat.H{
			"title": "cat",
			"now": time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC),
		})
	})

	r.Run(":8000")
}