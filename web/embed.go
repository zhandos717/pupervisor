package web

import (
	"embed"
	"html/template"
	"io/fs"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed js/*.js
var staticFS embed.FS

func GetTemplatesFS() fs.FS {
	sub, _ := fs.Sub(templatesFS, "templates")
	return sub
}

func GetStaticFS() fs.FS {
	return staticFS
}

func ParseTemplates() (*template.Template, error) {
	return template.ParseFS(templatesFS, "templates/*.html")
}
