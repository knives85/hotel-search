// Package opensearch adapts OpenSearch to the domain search and stats ports.
//
// TODO: wire the official client github.com/opensearch-project/opensearch-go/v4
// (use the AWS SigV4 signer when targeting Amazon OpenSearch).
package opensearch

import (
	"context"

	"github.com/knives85/hotel-search/internal/domain"
)

// Repository implements the OpenSearch-backed read ports.
type Repository struct {
	// TODO: hold the opensearch-go client and the index/alias name here.
}

// NewRepository builds an OpenSearch repository.
func NewRepository() *Repository {
	return &Repository{}
}

// Search runs the multi-filter hotel query.
func (r *Repository) Search(ctx context.Context, q domain.HotelSearchQuery) (domain.SearchResult, error) {
	return domain.SearchResult{}, domain.ErrNotImplemented
}

// SidebarFilterCounts runs the terms aggregations for the badge counts.
func (r *Repository) SidebarFilterCounts(ctx context.Context, q domain.HotelSearchQuery) (domain.SidebarFilterCounts, error) {
	return domain.SidebarFilterCounts{}, domain.ErrNotImplemented
}

// Stats runs the aggregate stats query.
func (r *Repository) Stats(ctx context.Context, q domain.HotelSearchQuery) (domain.HotelStats, error) {
	return domain.HotelStats{}, domain.ErrNotImplemented
}

// Compile-time checks that Repository satisfies the intended ports.
var (
	_ domain.SearchPort     = (*Repository)(nil)
	_ domain.HotelStatsPort = (*Repository)(nil)
)
