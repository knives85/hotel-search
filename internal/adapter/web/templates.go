package web

import (
	"bytes"
	"embed"
	"html/template"
	"net/http"
	"slices"

	"github.com/knives85/hotel-search/internal/domain"
)

//go:embed templates/hotels/*.html
var templatesFS embed.FS

// templates is parsed once at startup; html/template is safe for concurrent use.
var templates = template.Must(
	template.New("hotels").
		Funcs(templateFuncs).
		ParseFS(templatesFS, "templates/hotels/*.html"),
)

// templateFuncs are the helpers referenced from the .html files.
var templateFuncs = template.FuncMap{
	"seq":         seq,
	"contains":    slices.Contains[[]string, string],
	"starOptions": func() []starOption { return starOptionsList },
}

// indexView is the data passed to the "index" template: the search result for
// the table and the query so the filter form can mark its inputs as checked.
type indexView struct {
	Result domain.SearchResult
	Query  domain.HotelSearchQuery
}

// starOption is a (rating-key, star-count) pair driving the star-rating
// checkbox loop. Count is 0 for the UNRATED row.
type starOption struct {
	Value string
	Count int
}

var starOptionsList = []starOption{
	{Value: "SIX", Count: 6},
	{Value: "FIVE", Count: 5},
	{Value: "FOUR", Count: 4},
	{Value: "THREE", Count: 3},
	{Value: "TWO", Count: 2},
	{Value: "ONE", Count: 1},
	{Value: "UNRATED", Count: 0},
}

// seq returns [1, n] so templates can do `{{range seq n}}` — html/template
// has no numeric range built-in.
func seq(n int) []int {
	if n <= 0 {
		return nil
	}
	out := make([]int, n)
	for i := range out {
		out[i] = i + 1
	}
	return out
}

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
