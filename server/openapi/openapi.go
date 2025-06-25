package openapi

import (
	"embed"
	"fmt"
	"io/fs"
	"strings"
)

//go:embed all:api/*
var FS embed.FS

func init() {
	// recursively print out the files in the embedded FS
	fs.WalkDir(FS, ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".swagger.json") {
			fmt.Println(path)
		}
		return nil
	})
}
