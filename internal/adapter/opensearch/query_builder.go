package opensearch

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/knives85/hotel-search/internal/domain"
)

const (
	indexHotels = "hotels"

	starRatingUnrated   = "UNRATED"
	maxPrefixExpansions = 10000

	aggMaxLastUpdate      = "agg_max_last_update"
	aggMaxNumberOfReviews = "agg_max_number_of_reviews"
	aggStarRating         = "agg_star_rating"
)

// buildSearchRequest serialises a HotelSearchQuery into the JSON body of an
// OpenSearch /_search request: paginated hits + two max aggregations for the
// results header. No facet/sidebar aggregations here.
func buildSearchRequest(q domain.HotelSearchQuery) ([]byte, error) {
	body := map[string]any{
		"from":             q.Page * q.PageSize,
		"size":             q.PageSize,
		"track_total_hits": true,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": buildAllFilters(q),
			},
		},
		"aggs": map[string]any{
			aggMaxLastUpdate:      maxAgg("last_update_date"),
			aggMaxNumberOfReviews: maxAgg("number_of_reviews"),
		},
	}
	return json.Marshal(body)
}

func buildSidebarCountRequest(q domain.HotelSearchQuery) ([]byte, error) {
	selectedFilters := buildAllFilters(q)
	byStarAggQuery := q
	byStarAggQuery.StarRatings = nil
	body := map[string]any{
		"size": 0,
		"query": map[string]any{
			"bool": map[string]any{
				"filter": selectedFilters,
			},
		},
		"aggs": map[string]any{
			aggStarRating: filterAgg("star_rating", buildAllFilters(byStarAggQuery), map[string]any{
				"field":   "star_rating",
				"size":    10,
				"missing": "UNRATED",
			}),
		},
	}

	respAsJson, _ := json.MarshalIndent(body, "", "  ")
	fmt.Printf("Search Request:\n%s\n", string(respAsJson))
	return json.Marshal(body)
}

// buildAllFilters returns the bool/filter clause list. Order matches the
// original module's order so reviewing the JSON side-by-side is easy.
func buildAllFilters(q domain.HotelSearchQuery) []map[string]any {
	filters := []map[string]any{
		termQuery("index_status", "COMPLETE"),
	}

	if q.UniqueID != nil {
		filters = append(filters, termInt64Query("unique_id", *q.UniqueID))
	}
	if len(q.UniqueIDs) > 0 {
		filters = append(filters, termsInt64Query("unique_id", q.UniqueIDs))
	}

	// Skip sell_status when looking up by unique identifier: both active and
	// inactive hotels must be reachable regardless of the browse toggle.
	if q.UniqueID == nil && len(q.UniqueIDs) == 0 {
		if q.SellStatus != nil {
			filters = append(filters, termQuery("sell_status", strconv.FormatBool(*q.SellStatus)))
		}
	}

	if q.HotelName != nil && *q.HotelName != "" {
		filters = append(filters, matchPhrasePrefixQuery("hotel_name", *q.HotelName))
	}

	if f := starRatingFilter(q.StarRatings); f != nil {
		filters = append(filters, f)
	}

	if len(q.Types) > 0 {
		filters = append(filters, termsQuery("type", q.Types))
	}

	if len(q.CountryCodes) > 0 {
		filters = append(filters, geoFilterQuery("country.code", q.CountryCodes))
	}
	if len(q.CityCodes) > 0 {
		filters = append(filters, geoFilterQuery("city.code", q.CityCodes))
	}
	if len(q.RegionCodes) > 0 {
		filters = append(filters, geoFilterQuery("admin_region.code", q.RegionCodes))
	}
	if len(q.TouristicRegionCodes) > 0 {
		filters = append(filters, geoFilterQuery("touristic_region.code", q.TouristicRegionCodes))
	}
	if len(q.NonAdminCityCodes) > 0 {
		filters = append(filters, geoFilterQuery("nonadmin_city.code", q.NonAdminCityCodes))
	}
	if len(q.NeighbourhoodCodes) > 0 {
		filters = append(filters, geoFilterQuery("neighbourhood.code", q.NeighbourhoodCodes))
	}

	if len(q.ChainCodes) > 0 {
		filters = append(filters, termsQuery("chain.code", q.ChainCodes))
	}

	// AND semantics — each facility / badge becomes its own term filter.
	for _, code := range q.FacilityCodes {
		filters = append(filters, termQuery("facility_codes", code))
	}
	for _, code := range q.BadgeCodes {
		filters = append(filters, termQuery("badges", code))
	}

	// AND semantics for POIs too, but each clause is wrapped in a nested
	// query because nearby_pois is a nested mapping in the index.
	for _, code := range q.PoiCodes {
		filters = append(filters, nestedQuery("nearby_pois", termQuery("nearby_pois.code", code)))
	}

	if q.ContentScoreRange != nil {
		filters = append(filters, rangeIntQuery("content_score", q.ContentScoreRange.Min, q.ContentScoreRange.Max))
	}
	if q.ReviewScoreRange != nil {
		filters = append(filters, rangeIntQuery("review_score", q.ReviewScoreRange.Min, q.ReviewScoreRange.Max))
	}
	if q.NumberOfReviewsRange != nil {
		filters = append(filters, rangeIntQuery("number_of_reviews", q.NumberOfReviewsRange.Min, q.NumberOfReviewsRange.Max))
	}
	if q.LocationScoreRange != nil {
		filters = append(filters, rangeIntQuery("location_score", q.LocationScoreRange.Min, q.LocationScoreRange.Max))
	}
	if q.CreationDateRange != nil {
		filters = append(filters, rangeInt64Query("creation_date", q.CreationDateRange.Min, q.CreationDateRange.Max))
	}

	return filters
}

// ---- DSL helpers ----

func termQuery(field, value string) map[string]any {
	return map[string]any{"term": map[string]any{field: map[string]any{"value": value}}}
}

func termInt64Query(field string, value int64) map[string]any {
	return map[string]any{"term": map[string]any{field: map[string]any{"value": value}}}
}

func termsQuery(field string, values []string) map[string]any {
	return map[string]any{"terms": map[string]any{field: values}}
}

func termsInt64Query(field string, values []int64) map[string]any {
	return map[string]any{"terms": map[string]any{field: values}}
}

func matchPhrasePrefixQuery(field, value string) map[string]any {
	return map[string]any{
		"match_phrase_prefix": map[string]any{
			field: map[string]any{
				"query":          value,
				"max_expansions": maxPrefixExpansions,
			},
		},
	}
}

func rangeIntQuery(field string, gte, lte int) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{"gte": gte, "lte": lte}}}
}

func rangeInt64Query(field string, gte, lte int64) map[string]any {
	return map[string]any{"range": map[string]any{field: map[string]any{"gte": gte, "lte": lte}}}
}

func nestedQuery(path string, inner map[string]any) map[string]any {
	return map[string]any{"nested": map[string]any{"path": path, "query": inner}}
}

func maxAgg(field string) map[string]any {
	return map[string]any{"max": map[string]any{"field": field}}
}

func filterAgg(filterKey string, filters []map[string]any, terms map[string]any) map[string]any {
	return map[string]any{
		"global": map[string]any{},
		"aggs": map[string]any{
			"f": map[string]any{
				"filter": map[string]any{
					"bool": map[string]any{"filter": filters},
				},
				"aggs": map[string]any{
					"b": map[string]any{"terms": terms},
				},
			},
		},
	}
}

// missingFieldQuery → bool { must_not: [{ exists: field }] }.
func missingFieldQuery(field string) map[string]any {
	return map[string]any{
		"bool": map[string]any{
			"must_not": []map[string]any{
				{"exists": map[string]any{"field": field}},
			},
		},
	}
}

// missingOrEmptyQuery matches either a missing/null field or one indexed as
// the legacy empty string.
func missingOrEmptyQuery(field string) map[string]any {
	return shouldOne(
		missingFieldQuery(field),
		termQuery(field, ""),
	)
}

// shouldOne wraps the given clauses in bool { should: [...], minimum_should_match: "1" }.
func shouldOne(clauses ...map[string]any) map[string]any {
	return map[string]any{
		"bool": map[string]any{
			"should":               clauses,
			"minimum_should_match": "1",
		},
	}
}

// starRatingFilter implements the rated / unrated / mixed branches.
func starRatingFilter(ratings []string) map[string]any {
	if len(ratings) == 0 {
		return nil
	}
	rated := make([]string, 0, len(ratings))
	includeUnrated := false
	for _, r := range ratings {
		if r == starRatingUnrated {
			includeUnrated = true
		} else {
			rated = append(rated, r)
		}
	}
	switch {
	case !includeUnrated:
		return termsQuery("star_rating", rated)
	case len(rated) == 0:
		return missingOrEmptyQuery("star_rating")
	default:
		return shouldOne(
			termsQuery("star_rating", rated),
			missingOrEmptyQuery("star_rating"),
		)
	}
}

// geoFilterQuery handles the UNMAPPED_GEO_CODE sentinel: a request for the
// sentinel becomes a "missing or empty" clause on the geo field.
func geoFilterQuery(field string, codes []string) map[string]any {
	real := make([]string, 0, len(codes))
	hasUnmapped := false
	for _, c := range codes {
		if c == domain.UNMAPPED_GEO_CODE {
			hasUnmapped = true
		} else {
			real = append(real, c)
		}
	}
	switch {
	case !hasUnmapped:
		return termsQuery(field, real)
	case len(real) == 0:
		return missingOrEmptyQuery(field)
	default:
		return shouldOne(
			termsQuery(field, real),
			missingOrEmptyQuery(field),
		)
	}
}
