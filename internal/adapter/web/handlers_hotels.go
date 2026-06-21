package web

import "net/http"

// Handlers for the /hotels routes. Each one currently returns 501 and will be
// filled in as the OpenSearch/Postgres adapters and templates are ported.

func (s *Server) handleHotelsIndex(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /hotels")
}

func (s *Server) handleHotelsResults(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /hotels/results")
}

func (s *Server) handleHotelsStats(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /hotels/stats")
}

func (s *Server) handleHotelsFilterCounts(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /hotels/filter-counts")
}

func (s *Server) handleHotelsExportCSV(w http.ResponseWriter, _ *http.Request) {
	notImplemented(w, "GET /hotels/export.csv")
}

// handleSuggest backs every autocomplete endpoint (country, city, facility,
// chain, region, POI, ...). TODO: split per suggester when implementing.
func (s *Server) handleSuggest(w http.ResponseWriter, r *http.Request) {
	notImplemented(w, "GET "+r.URL.Path)
}
