// Package domain holds the core types and the port interfaces of the app
// module. It has no dependency on any framework, database driver or HTTP
// library: adapters depend on the domain, never the other way around.
package domain

import "errors"

// ErrNotImplemented is returned by adapter and use-case stubs that have not
// been ported yet. Replace these as the implementation grows.
var ErrNotImplemented = errors.New("not implemented")

// HotelSearchQuery captures the filters coming from the search UI.
//
// TODO: port the real fields from the Kotlin HotelSearchQuery — country,
// type, star rating, facilities, chain, geo location, date range and
// pagination.
type HotelSearchQuery struct {
	Term   string
	Limit  int
	Offset int
}

// Hotel is the projection of a hotel document returned to the UI.
//
// TODO: port the real fields from the OpenSearch hotel document.
type Hotel struct {
	ID   string
	Name string
}

// SearchResult is a page of hotels plus the total match count.
type SearchResult struct {
	Hotels []Hotel
	Total  int64
}

// SidebarFilterCounts holds the per-option badge counts (FILT-005).
//
// TODO: port the real aggregation buckets per filter dimension.
type SidebarFilterCounts struct{}

// HotelStats holds the aggregate stats shown above the results table.
type HotelStats struct{}

// Suggestion is a single autocomplete entry (code + human description).
type Suggestion struct {
	Code        string
	Description string
}

// Job represents a background job (e.g. a CSV export) shown on /jobs.
type Job struct {
	ID     string
	Status string
}

// InventoryList is a saved set of filters shown on /inventory-lists.
type InventoryList struct {
	ID   string
	Name string
}
