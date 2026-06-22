// Package web wires the HTTP routes of the app module. The route set and the
// URLs mirror the Kotlin controllers exactly, so the existing HTMX UI keeps
// working unchanged once the handlers are implemented.
package web

import (
	"net/http"

	"github.com/knives85/hotel-search/internal/domain"
)

// Deps holds the ports and use cases the HTTP layer depends on. Fields may
// be left nil in tests when the exercised handler does not call into them.
type Deps struct {
	Search domain.SearchPort
}

// Server holds the HTTP router for the app module.
type Server struct {
	mux         *http.ServeMux
	contextPath string
	deps        Deps
}

// NewServer builds the router and registers all routes under contextPath
// (e.g. "/hotel-search").
func NewServer(contextPath string, deps Deps) *Server {
	s := &Server{mux: http.NewServeMux(), contextPath: contextPath, deps: deps}
	s.routes()
	return s
}

// Handler returns the root http.Handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	// Liveness/readiness probe (kept outside the app context path).
	s.mux.HandleFunc("GET /healthz", s.handleHealth)

	p := s.contextPath

	// Hotel search — the heart of the module.
	s.mux.HandleFunc("GET "+p+"/hotels", s.handleHotelsIndex)
	s.mux.HandleFunc("GET "+p+"/hotels/results", s.handleHotelsResults)
	s.mux.HandleFunc("GET "+p+"/hotels/stats", s.handleHotelsStats)
	s.mux.HandleFunc("GET "+p+"/hotels/filter-counts", s.handleHotelsFilterCounts)
	s.mux.HandleFunc("GET "+p+"/hotels/export.csv", s.handleHotelsExportCSV)

	// Autocomplete suggesters.
	s.mux.HandleFunc("GET "+p+"/hotels/country-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/city-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/non-admin-city-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/admin-region-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/touristic-region-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/neighbourhood-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/poi-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/facility-suggest", s.handleSuggest)
	s.mux.HandleFunc("GET "+p+"/hotels/chain-suggest", s.handleSuggest)

	// Jobs.
	s.mux.HandleFunc("GET "+p+"/jobs", s.handleJobsIndex)
	s.mux.HandleFunc("GET "+p+"/jobs/{id}/row", s.handleJobRow)
	s.mux.HandleFunc("GET "+p+"/jobs/{id}/download", s.handleJobDownload)
}

// notImplemented is the placeholder response for routes not yet ported.
func notImplemented(w http.ResponseWriter, name string) {
	http.Error(w, "not implemented: "+name, http.StatusNotImplemented)
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
