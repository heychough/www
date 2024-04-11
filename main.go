package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.Static("/assets", "assets")

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	e.GET("/", blogs)
	e.GET("*", notFound)
	e.GET("/blog/:path", blogs)

    e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}

func getMarkdownUtilities() (*parser.Parser, *html.Renderer) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return p, renderer
}

func index(c echo.Context) error {
	file, err := os.ReadFile("writing/index.md")

	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	p, renderer := getMarkdownUtilities()
	html := markdown.ToHTML(file, p, renderer)

	return c.Render(http.StatusOK, "index", map[string]interface{}{
		"content": string(html),
		"isIndex": true,
	})
}

func blogs(c echo.Context) error {
	path := c.Param("path")
	if path == "" {
		path = "index"
	}

	file, err := os.ReadFile("writing/" + path + ".md")

	if err != nil {
		fmt.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	p, renderer := getMarkdownUtilities()
	html := markdown.ToHTML(file, p, renderer)

	return c.Render(http.StatusOK, "index", map[string]interface{}{
		"content": string(html),
		"isIndex": path == "index",
	})
}

func notFound(c echo.Context) error {
	return c.String(http.StatusNotFound, "not found")
}
