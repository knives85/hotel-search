// Package opensearch adapts OpenSearch to the domain search and stats ports.
package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"

	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/knives85/hotel-search/internal/domain"
)

// Repository implements the OpenSearch-backed read ports.
type Repository struct {
	client *opensearchapi.Client
	index  string
}

// NewRepository builds a Repository that targets the given index via client.
func NewRepository(client *opensearchapi.Client, index string) *Repository {
	return &Repository{client: client, index: index}
}

// Search runs the multi-filter hotel query and maps the response to the
// domain projection.
func (r *Repository) Search(ctx context.Context, q domain.HotelSearchQuery) (domain.SearchResult, error) {
	body, err := buildSearchRequest(q)
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("build search request: %w", err)
	}

	resp, err := r.client.Search(ctx, &opensearchapi.SearchReq{
		Indices: []string{r.index},
		Body:    bytes.NewReader(body),
	})
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("opensearch search: %w", err)
	}

	hotels := make([]domain.Hotel, 0, len(resp.Hits.Hits))
	for _, hit := range resp.Hits.Hits {
		var doc hotelSearchDocument
		if err := json.Unmarshal(hit.Source, &doc); err != nil {
			return domain.SearchResult{}, fmt.Errorf("decode hit %q: %w", hit.ID, err)
		}
		hotels = append(hotels, toHotel(doc))
	}

	maxLastUpdate, maxNumberOfReviews, err := parseMaxAggregations(resp.Aggregations)
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("decode aggregations: %w", err)
	}

	return domain.SearchResult{
		Hotels:             hotels,
		Total:              int64(resp.Hits.Total.Value),
		Page:               q.Page,
		PageSize:           q.PageSize,
		LastUpdateDate:     maxLastUpdate,
		MaxNumberOfReviews: maxNumberOfReviews,
	}, nil
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

// aggregations is the slice of the OpenSearch response body that we care
// about for the results page: two top-level max aggregations.
type aggregationsResponse struct {
	MaxLastUpdate      maxAggValue `json:"agg_max_last_update"`
	MaxNumberOfReviews maxAggValue `json:"agg_max_number_of_reviews"`
}

type maxAggValue struct {
	Value *float64 `json:"value"`
}

// parseMaxAggregations decodes the two max aggregations and applies the same
// finite-and-positive guard as the original: zero / negative / NaN → nil.
func parseMaxAggregations(raw json.RawMessage) (*int64, *int, error) {
	if len(raw) == 0 {
		return nil, nil, nil
	}
	var aggs aggregationsResponse
	if err := json.Unmarshal(raw, &aggs); err != nil {
		return nil, nil, err
	}
	return finitePositiveInt64(aggs.MaxLastUpdate.Value),
		finitePositiveInt(aggs.MaxNumberOfReviews.Value),
		nil
}

func finitePositiveInt64(v *float64) *int64 {
	if v == nil || !isFinitePositive(*v) {
		return nil
	}
	n := int64(*v)
	return &n
}

func finitePositiveInt(v *float64) *int {
	if v == nil || !isFinitePositive(*v) {
		return nil
	}
	n := int(*v)
	return &n
}

func isFinitePositive(f float64) bool {
	return !math.IsNaN(f) && !math.IsInf(f, 0) && f > 0
}
