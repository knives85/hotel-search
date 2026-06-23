package web

import "net/http"

// Handlers for the /hotels routes.

func (s *Server) handleHotelsIndex(w http.ResponseWriter, r *http.Request) {
	params, err := parseHotelSearchParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query, err := buildQuery(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := s.deps.Search.Search(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderHTML(w, "index", indexView{Result: result, Query: query})
}

func (s *Server) handleHotelsResults(w http.ResponseWriter, r *http.Request) {
	params, err := parseHotelSearchParams(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query, err := buildQuery(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := s.deps.Search.Search(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderHTML(w, "results-table", result)
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
