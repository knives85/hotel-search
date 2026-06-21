// Package usecase holds the application logic, orchestrating the domain
// ports. Use cases are deliberately thin: they depend only on interfaces
// from the domain package, never on concrete adapters.
package usecase

import (
	"context"

	"github.com/knives85/hotel-search/internal/domain"
)

// SearchHotels is the use case behind GET /hotels/results.
type SearchHotels struct {
	Search domain.SearchPort
}

// Execute runs the hotel search for the given query.
func (uc SearchHotels) Execute(ctx context.Context, q domain.HotelSearchQuery) (domain.SearchResult, error) {
	return uc.Search.Search(ctx, q)
}

// GetSidebarFilterCounts is the use case behind GET /hotels/filter-counts.
type GetSidebarFilterCounts struct {
	Search domain.SearchPort
}

// Execute computes the per-option badge counts for the given query.
func (uc GetSidebarFilterCounts) Execute(ctx context.Context, q domain.HotelSearchQuery) (domain.SidebarFilterCounts, error) {
	return uc.Search.SidebarFilterCounts(ctx, q)
}

// GetHotelStats is the use case behind GET /hotels/stats.
type GetHotelStats struct {
	Stats domain.HotelStatsPort
}

// Execute computes the aggregate stats for the given query.
func (uc GetHotelStats) Execute(ctx context.Context, q domain.HotelSearchQuery) (domain.HotelStats, error) {
	return uc.Stats.Stats(ctx, q)
}
