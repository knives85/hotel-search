package web

import (
	"encoding/json"
	"net/http"
)

// Handlers for the /hotels routes. Each one currently returns 501 and will be
// filled in as the OpenSearch/Postgres adapters and templates are ported.

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

	counts, err := s.deps.Search.SidebarFilterCounts(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var hotelDto []map[string]any
	for _, hotel := range result.Hotels {
		hotelDto = append(hotelDto, map[string]any{
			"uniqueId": hotel.UniqueID,
			"name":     hotel.HotelName,
		})
	}

	dto := map[string]any{
		"total":  result.Total,
		"hotels": hotelDto,
		"counts": counts.ByStarRating,
	}

	marshalled, err := json.Marshal(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(marshalled)
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
