package web

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/hotels/*.html
var templatesFS embed.FS

// templates is parsed once at startup; html/template is safe for concurrent use.
var templates = template.Must(template.ParseFS(templatesFS, "templates/hotels/*.html"))

// renderHTML executes a named template into a buffer first so that template
// errors surface as 500 with a clean response, not a half-written body.
func renderHTML(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, name, data); err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(buf.Bytes())
}
