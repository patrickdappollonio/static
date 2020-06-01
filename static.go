package static

import (
	"net/http"
	"os"
	"strings"
)

var defaults = []string{
	"favicon.ico",
	"robots.txt",
}

// Static is a Go-standard middleware that receives a
// folder path (the "path" parameter) and a list of optional
// assets that are usually delivered from outside the assets
// folder, such as favicons, robots.txt and similar. If an
// asset list is passed, then it's read from the project's root.
func Static(path string, assets []string) func(http.Handler) http.Handler {
	return st(path, false, assets)
}

// StaticWildcard works the same way as Static, except that
// whatever value specified after the path and before the asset
// destination will be responded, so assuming there's a folder with
// file "examples/01/01-how-to-middleware.txt" and "path" is set to
// "examples", then any value after examples and before
// "/01/01-how-to-middleware.txt" will be accepted, like:
// examples/abc123/01/01-how-to-middleware.txt
// examples/111222333/01/01-how-to-middleware.txt
func StaticWildcard(path string, assets []string) func(http.Handler) http.Handler {
	return st(path, true, assets)
}

func st(path string, allowWild bool, assets []string) func(http.Handler) http.Handler {
	extras := append(defaults, assets...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// We don't serve static files over a method that's not GET
			if r.Method != http.MethodGet {
				next.ServeHTTP(w, r)
				return
			}

			// Grab the path without the initial slash
			location := r.URL.Path[1:]

			// Check if it's a known file
			for _, v := range extras {
				if location == v {
					if serve(v, w, r, true) {
						return
					}
				}
			}

			// Check if it's the static folder
			if strings.HasPrefix(location, path) {
				if serve(location, w, r, true) {
					return
				}

				// If we allow wildcard, remove the extra item from the path
				if allowWild {
					nl := strings.TrimPrefix(location, path+"/")
					if x := strings.Index(nl, "/"); x != -1 {
						nl = path + nl[x:]

						if serve(nl, w, r, false) {
							return
						}
					}
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
func serve(path string, w http.ResponseWriter, r *http.Request, canonical bool) bool {
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
		defer f.Close()

		if !canonical {
			w.Header().Set("X-Robots-Tag", "noindex, follow")
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeContent(w, r, fi.Name(), fi.ModTime(), f)
		return true
	}

	return false
}
