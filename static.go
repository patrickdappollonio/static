package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var statics = []string{
	"favicon.ico",
	"robots.txt",
}

func Static(path string, assets []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We don't serve static files over a method that's not GET
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Grab the path without the initial slash
			originalPath := r.URL.Path[1:]

			// Add the static assets if the local ones
			// are empty
			if len(assets) == 0 {
				assets = []string{}
			}
			assets = append(assets, statics...)

			// Check if it's a known file
			for _, v := range assets {
				if originalPath == v {
					if serve(v, w, r) {
						return
					}
				}
			}

			// Check if it's the static folder
			if strings.HasPrefix(originalPath, path) {
				if serve(originalPath, w, r) {
					return
				}
			}

			// If it didn't match either of the flows above
			// let's try the handlers or 404
			next.ServeHTTP(w, r)
			return
		})
	}
}

// serve tries serving the static file, it will give us a boolean
// response back to the caller, which means if it was indeed handled or not
// if handled, it won't call subsequent elements from the stack, but
// if handled, it'll just return
func serve(originalPath string, w http.ResponseWriter, r *http.Request) bool {
	if path, err := filepath.Abs(originalPath); err == nil {
		if fi, err := os.Stat(path); err == nil && path != "" {
			if fi.IsDir() {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return true
			}

			f, err := os.Open(path)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return true
			}

			http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
			f.Close()

			return true
		}
	}

	return false
}
