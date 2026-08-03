package web

import (
	"io/fs"
	"net/http"

	"dbmw/frontend"
)

// GetFileSystem returns the http.FileSystem wrapper around the embedded frontend dist.
func GetFileSystem() (http.FileSystem, error) {
	sub, err := fs.Sub(frontend.DistFS, "dist")
	if err != nil {
		return nil, err
	}
	return http.FS(sub), nil
}
