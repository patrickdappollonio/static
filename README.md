# Static middleware for Go standard projects using `net/http`

This is just a simple, Go-standard middleware that allows to serve static assets from a
path in a router, and a folder in the filesystem.


### Usage

Simply call `static.Static()`, the function accepts two parameters:

* `path` (`string`): the path to the folder where the assets are served. This will also be the path in the URL.
* `assets` (`[]string`): a list of static files usually served from outside the assets folder, like `robots.txt` or `favicon.ico`

Example:

```go
// Example with go-chi

r := chi.NewRouter()
r.Use(static.Static("assets", []string{
  "robots.txt",
  "favicon.ico",
}))
```

`static` should work with any muxer / router that allows middlewares with the format `func (http.Handler) http.Handler`.
