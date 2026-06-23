package web

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

// mountStatic registers a file server for the embedded static/ directory
// under contextPath + "/static/". The StripPrefix step removes the URL prefix
// so that lookups inside the file server are relative to the static/ root —
// e.g. GET /hotel-search/static/app.css → opens "app.css" within the sub FS.
func mountStatic(mux *http.ServeMux, contextPath string) {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		// Unreachable when the embed pattern is valid; treated as a programmer
		// error rather than a runtime concern.
		panic("static embed root missing: " + err.Error())
	}
	prefix := contextPath + "/static/"
	mux.Handle("GET "+prefix, http.StripPrefix(prefix, http.FileServer(http.FS(sub))))
}
