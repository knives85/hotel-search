// Package domain holds the core types and the port interfaces of the app
// module. It has no dependency on any framework, database driver or HTTP
// library: adapters depend on the domain, never the other way around.
package domain

import "errors"

// ErrNotImplemented is returned by adapter and use-case stubs that have not
// been ported yet. Replace these as the implementation grows.
var ErrNotImplemented = errors.New("not implemented")

// UNMAPPED_GEO_CODE is the sentinel used by the OpenSearch documents when a
// geo dimension (city, region, ...) has no resolved code.
const UNMAPPED_GEO_CODE = "__unmapped__"

// PageSize is the default page size for the hotel search results table.
const PageSize = 200

// CSVPageSize is the page size used for the CSV export endpoint.
const CSVPageSize = 10_000

// StarRatingOptions is the closed set of accepted star-rating filter values.
// Anything outside this set is dropped by buildQuery.
var StarRatingOptions = []string{"ONE", "TWO", "THREE", "FOUR", "FIVE", "SIX", "UNRATED"}

// IntRange is the closed integer interval [Min, Max] used for review-score
// and number-of-reviews filters.
type IntRange struct {
	Min int
	Max int
}

// Int64Range is the closed interval [Min, Max] over int64. Used for the
// creation-date filter (epoch millis).
type Int64Range struct {
	Min int64
	Max int64
}

// HotelSearchQuery captures the filters coming from the search UI.
//
// Pointer fields are "nullable": a nil pointer means the user did not provide
// that filter. Slice fields use the same convention: nil == absent, non-nil
// (even if empty) is treated as "present" — but buildQuery normalises empty
// slices to nil before populating this struct.
type HotelSearchQuery struct {
	UniqueID             *int64
	UniqueIDs            []int64
	HotelName            *string
	SellStatus           *bool
	StarRatings          []string
	Types                []string
	CountryCodes         []string
	CityCodes            []string
	CityNamePrefix       *string
	RegionCodes          []string
	TouristicRegionCodes []string
	NonAdminCityCodes    []string
	PoiCodes             []string
	NeighbourhoodCodes   []string
	ChainCodes           []string
	FacilityCodes        []string
	BadgeCodes           []string
	ContentScoreRange    *IntRange
	ReviewScoreRange     *IntRange
	NumberOfReviewsRange *IntRange
	LocationScoreRange   *IntRange
	CreationDateRange    *Int64Range
	Page                 int
	PageSize             int
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
