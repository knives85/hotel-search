package opensearch

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/knives85/hotel-search/internal/domain"
)

// assertJSONEqual fails the test unless got and want decode to the same
// untyped tree (order of object keys ignored, slice order preserved).
func assertJSONEqual(t *testing.T, got []byte, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("got is not valid JSON: %v\n%s", err, got)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want is not valid JSON: %v", err)
	}
	if !reflect.DeepEqual(g, w) {
		gp, _ := json.MarshalIndent(g, "", "  ")
		wp, _ := json.MarshalIndent(w, "", "  ")
		t.Fatalf("JSON mismatch.\nGOT:\n%s\nWANT:\n%s", gp, wp)
	}
}

func ptr[T any](v T) *T { return &v }

func TestBuildSearchRequest_Skeleton_EmptyQuery(t *testing.T) {
	q := domain.HotelSearchQuery{Page: 0, PageSize: domain.PageSize}
	got, err := buildSearchRequest(q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `{
		"from": 0,
		"size": 200,
		"track_total_hits": true,
		"query": {
			"bool": {
				"filter": [
					{"term": {"index_status": {"value": "COMPLETE"}}}
				]
			}
		},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_Pagination(t *testing.T) {
	q := domain.HotelSearchQuery{Page: 3, PageSize: 50}
	got, err := buildSearchRequest(q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var parsed map[string]any
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if parsed["from"].(float64) != 150 {
		t.Errorf("from = %v, want 150", parsed["from"])
	}
	if parsed["size"].(float64) != 50 {
		t.Errorf("size = %v, want 50", parsed["size"])
	}
}

func TestBuildSearchRequest_UniqueID(t *testing.T) {
	q := domain.HotelSearchQuery{UniqueID: ptr(int64(42)), PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"term": {"unique_id":   {"value": 42}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_SellStatus_SuppressedWhenUniqueID(t *testing.T) {
	// uniqueId present → sell_status filter must NOT be added.
	q := domain.HotelSearchQuery{
		UniqueID:   ptr(int64(42)),
		SellStatus: ptr(true),
		PageSize:   200,
	}
	got, _ := buildSearchRequest(q)
	var parsed struct {
		Query struct {
			Bool struct {
				Filter []map[string]any `json:"filter"`
			} `json:"bool"`
		} `json:"query"`
	}
	if err := json.Unmarshal(got, &parsed); err != nil {
		t.Fatal(err)
	}
	for _, f := range parsed.Query.Bool.Filter {
		if term, ok := f["term"].(map[string]any); ok {
			if _, has := term["sell_status"]; has {
				t.Errorf("sell_status filter should be suppressed, got %v", f)
			}
		}
	}
}

func TestBuildSearchRequest_SellStatus_AppliedAlone(t *testing.T) {
	q := domain.HotelSearchQuery{SellStatus: ptr(true), PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"term": {"sell_status":  {"value": "true"}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_HotelName_MatchPhrasePrefix(t *testing.T) {
	q := domain.HotelSearchQuery{HotelName: ptr("Hilton"), PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"match_phrase_prefix": {"hotel_name": {"query": "Hilton", "max_expansions": 10000}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_StarRatings_Rated(t *testing.T) {
	q := domain.HotelSearchQuery{StarRatings: []string{"FOUR", "FIVE"}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term":  {"index_status": {"value": "COMPLETE"}}},
			{"terms": {"star_rating": ["FOUR", "FIVE"]}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_StarRatings_UnratedOnly(t *testing.T) {
	q := domain.HotelSearchQuery{StarRatings: []string{"UNRATED"}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"bool": {
				"should": [
					{"bool": {"must_not": [{"exists": {"field": "star_rating"}}]}},
					{"term": {"star_rating": {"value": ""}}}
				],
				"minimum_should_match": "1"
			}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_StarRatings_RatedAndUnrated(t *testing.T) {
	q := domain.HotelSearchQuery{StarRatings: []string{"ONE", "UNRATED"}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"bool": {
				"should": [
					{"terms": {"star_rating": ["ONE"]}},
					{"bool": {
						"should": [
							{"bool": {"must_not": [{"exists": {"field": "star_rating"}}]}},
							{"term": {"star_rating": {"value": ""}}}
						],
						"minimum_should_match": "1"
					}}
				],
				"minimum_should_match": "1"
			}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_GeoCodes_Plain(t *testing.T) {
	q := domain.HotelSearchQuery{CountryCodes: []string{"IT", "FR"}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term":  {"index_status": {"value": "COMPLETE"}}},
			{"terms": {"country.code": ["IT", "FR"]}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_GeoCodes_UnmappedOnly(t *testing.T) {
	q := domain.HotelSearchQuery{CityCodes: []string{domain.UNMAPPED_GEO_CODE}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"bool": {
				"should": [
					{"bool": {"must_not": [{"exists": {"field": "city.code"}}]}},
					{"term": {"city.code": {"value": ""}}}
				],
				"minimum_should_match": "1"
			}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_GeoCodes_Mixed(t *testing.T) {
	q := domain.HotelSearchQuery{
		RegionCodes: []string{"LAZ", domain.UNMAPPED_GEO_CODE},
		PageSize:    200,
	}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"bool": {
				"should": [
					{"terms": {"admin_region.code": ["LAZ"]}},
					{"bool": {
						"should": [
							{"bool": {"must_not": [{"exists": {"field": "admin_region.code"}}]}},
							{"term": {"admin_region.code": {"value": ""}}}
						],
						"minimum_should_match": "1"
					}}
				],
				"minimum_should_match": "1"
			}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_FacilityCodes_AND(t *testing.T) {
	// Each facility becomes its OWN term filter (AND semantics).
	q := domain.HotelSearchQuery{FacilityCodes: []string{"WIFI", "POOL"}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status":   {"value": "COMPLETE"}}},
			{"term": {"facility_codes": {"value": "WIFI"}}},
			{"term": {"facility_codes": {"value": "POOL"}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_PoiCodes_Nested(t *testing.T) {
	q := domain.HotelSearchQuery{PoiCodes: []string{"P1", "P2"}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term": {"index_status": {"value": "COMPLETE"}}},
			{"nested": {"path": "nearby_pois", "query": {"term": {"nearby_pois.code": {"value": "P1"}}}}},
			{"nested": {"path": "nearby_pois", "query": {"term": {"nearby_pois.code": {"value": "P2"}}}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_ReviewScoreRange(t *testing.T) {
	q := domain.HotelSearchQuery{
		ReviewScoreRange: &domain.IntRange{Min: 60, Max: 90},
		PageSize:         200,
	}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term":  {"index_status": {"value": "COMPLETE"}}},
			{"range": {"review_score": {"gte": 60, "lte": 90}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_CreationDateRange(t *testing.T) {
	q := domain.HotelSearchQuery{
		CreationDateRange: &domain.Int64Range{Min: 1704067200000, Max: 1706745599999},
		PageSize:          200,
	}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term":  {"index_status":  {"value": "COMPLETE"}}},
			{"range": {"creation_date": {"gte": 1704067200000, "lte": 1706745599999}}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}

func TestBuildSearchRequest_UniqueIDs(t *testing.T) {
	q := domain.HotelSearchQuery{UniqueIDs: []int64{1, 2, 3}, PageSize: 200}
	got, _ := buildSearchRequest(q)
	want := `{
		"from": 0, "size": 200, "track_total_hits": true,
		"query": {"bool": {"filter": [
			{"term":  {"index_status": {"value": "COMPLETE"}}},
			{"terms": {"unique_id": [1, 2, 3]}}
		]}},
		"aggs": {
			"agg_max_last_update":       {"max": {"field": "last_update_date"}},
			"agg_max_number_of_reviews": {"max": {"field": "number_of_reviews"}}
		}
	}`
	assertJSONEqual(t, got, want)
}
