package web

import (
	"embed"
	"io/fs"
)

//go:embed templates/*.html
var templatesFS embed.FS

//go:embed css/*.css js/*.js
var staticFS embed.FS

func GetTemplatesFS() fs.FS {
	sub, _ := fs.Sub(templatesFS, "templates")
	return sub
}

func GetStaticFS() fs.FS {
	return staticFS
}
