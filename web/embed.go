// Package web holds the page the generator serves.
package web

import (
	"embed"
	"io/fs"
)

//go:embed static
var files embed.FS

// Static is the page files, rooted at static/.
func Static() fs.FS {
	sub, err := fs.Sub(files, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
