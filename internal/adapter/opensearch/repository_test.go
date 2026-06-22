package opensearch

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	"github.com/knives85/hotel-search/internal/domain"
)

// stubOSServer returns an httptest server that captures the inbound request
// body and replies with the given canned response.
func stubOSServer(t *testing.T, response string) (*httptest.Server, *capturedRequest) {
	t.Helper()
	cap := &capturedRequest{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cap.method = r.Method
		cap.path = r.URL.Path
		body, _ := io.ReadAll(r.Body)
		cap.body = body
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, response)
	}))
	t.Cleanup(srv.Close)
	return srv, cap
}

type capturedRequest struct {
	method string
	path   string
	body   []byte
}

// newTestRepo wires a Repository against the stub server.
func newTestRepo(t *testing.T, srv *httptest.Server) *Repository {
	t.Helper()
	client, err := opensearchapi.NewClient(opensearchapi.Config{
		Client: opensearch.Config{Addresses: []string{srv.URL}},
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return NewRepository(client, indexHotels)
}

const canonicalResponse = `{
	"took": 1,
	"timed_out": false,
	"_shards": {"total": 1, "successful": 1, "skipped": 0, "failed": 0},
	"hits": {
		"total": {"value": 2, "relation": "eq"},
		"max_score": 1.0,
		"hits": [
			{
				"_index": "hotels",
				"_id": "42",
				"_score": 1.0,
				"_source": {
					"unique_id": 42,
					"index_status": "COMPLETE",
					"hotel_name": "Hilton Rome",
					"country": {"code": "IT", "name": "Italy"},
					"review_score": 90,
					"number_of_reviews": 1500
				}
			},
			{
				"_index": "hotels",
				"_id": "7",
				"_score": 0.5,
				"_source": {
					"unique_id": 7,
					"index_status": "COMPLETE",
					"hotel_name": "Hilton Milan",
					"country": {"code": "IT", "name": "Italy"},
					"review_score": 85,
					"number_of_reviews": 900
				}
			}
		]
	},
	"aggregations": {
		"agg_max_last_update":      {"value": 1.7356032e12},
		"agg_max_number_of_reviews": {"value": 1500.0}
	}
}`

func TestRepository_Search_HappyPath(t *testing.T) {
	srv, cap := stubOSServer(t, canonicalResponse)
	repo := newTestRepo(t, srv)

	q := domain.HotelSearchQuery{
		HotelName: ptr("Hilton"),
		Page:      0,
		PageSize:  200,
	}
	got, err := repo.Search(context.Background(), q)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}

	// Captured request: POST /hotels/_search with our query body.
	if cap.method != http.MethodPost {
		t.Errorf("method = %s, want POST", cap.method)
	}
	if !strings.HasPrefix(cap.path, "/hotels/_search") {
		t.Errorf("path = %s, want /hotels/_search...", cap.path)
	}
	var sentBody map[string]any
	if err := json.Unmarshal(cap.body, &sentBody); err != nil {
		t.Fatalf("captured body is not JSON: %v", err)
	}
	if sentBody["size"].(float64) != 200 {
		t.Errorf("sent size = %v, want 200", sentBody["size"])
	}

	// Result mapping.
	if got.Total != 2 {
		t.Errorf("Total = %d, want 2", got.Total)
	}
	if got.Page != 0 {
		t.Errorf("Page = %d, want 0", got.Page)
	}
	if got.PageSize != 200 {
		t.Errorf("PageSize = %d, want 200", got.PageSize)
	}
	if len(got.Hotels) != 2 {
		t.Fatalf("len(Hotels) = %d, want 2", len(got.Hotels))
	}
	if got.Hotels[0].UniqueID != 42 {
		t.Errorf("Hotels[0].UniqueID = %d", got.Hotels[0].UniqueID)
	}
	if got.Hotels[0].Country == nil || got.Hotels[0].Country.Code != "IT" {
		t.Errorf("Hotels[0].Country = %+v", got.Hotels[0].Country)
	}

	// Aggregations.
	if got.LastUpdateDate == nil || *got.LastUpdateDate != 1735603200000 {
		t.Errorf("LastUpdateDate = %v, want 1735603200000", got.LastUpdateDate)
	}
	if got.MaxNumberOfReviews == nil || *got.MaxNumberOfReviews != 1500 {
		t.Errorf("MaxNumberOfReviews = %v, want 1500", got.MaxNumberOfReviews)
	}
}

func TestRepository_Search_EmptyResult(t *testing.T) {
	emptyResp := `{
		"took": 1, "timed_out": false,
		"_shards": {"total": 1, "successful": 1, "skipped": 0, "failed": 0},
		"hits": {"total": {"value": 0, "relation": "eq"}, "hits": []},
		"aggregations": {
			"agg_max_last_update": {"value": null},
			"agg_max_number_of_reviews": {"value": null}
		}
	}`
	srv, _ := stubOSServer(t, emptyResp)
	repo := newTestRepo(t, srv)
	got, err := repo.Search(context.Background(), domain.HotelSearchQuery{PageSize: 200})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if got.Total != 0 {
		t.Errorf("Total = %d, want 0", got.Total)
	}
	if len(got.Hotels) != 0 {
		t.Errorf("Hotels = %v, want empty", got.Hotels)
	}
	if got.LastUpdateDate != nil {
		t.Errorf("LastUpdateDate = %v, want nil", got.LastUpdateDate)
	}
	if got.MaxNumberOfReviews != nil {
		t.Errorf("MaxNumberOfReviews = %v, want nil", got.MaxNumberOfReviews)
	}
}

func TestRepository_Search_AggregationsNonFinite_AreDropped(t *testing.T) {
	resp := `{
		"took": 1, "timed_out": false,
		"_shards": {"total": 1, "successful": 1, "skipped": 0, "failed": 0},
		"hits": {"total": {"value": 0, "relation": "eq"}, "hits": []},
		"aggregations": {
			"agg_max_last_update": {"value": 0},
			"agg_max_number_of_reviews": {"value": -1}
		}
	}`
	srv, _ := stubOSServer(t, resp)
	repo := newTestRepo(t, srv)
	got, err := repo.Search(context.Background(), domain.HotelSearchQuery{PageSize: 200})
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	// 0 / negative ignored: nil.
	if got.LastUpdateDate != nil {
		t.Errorf("LastUpdateDate = %v, want nil for non-positive value", got.LastUpdateDate)
	}
	if got.MaxNumberOfReviews != nil {
		t.Errorf("MaxNumberOfReviews = %v, want nil for non-positive value", got.MaxNumberOfReviews)
	}
}

func TestRepository_Search_PreservesQueryPagination(t *testing.T) {
	srv, _ := stubOSServer(t, canonicalResponse)
	repo := newTestRepo(t, srv)
	got, _ := repo.Search(context.Background(), domain.HotelSearchQuery{Page: 3, PageSize: 50})
	if got.Page != 3 || got.PageSize != 50 {
		t.Errorf("Page/PageSize = %d/%d, want 3/50", got.Page, got.PageSize)
	}
}

// Stable order: results follow the OpenSearch hits order verbatim.
func TestRepository_Search_Order(t *testing.T) {
	srv, _ := stubOSServer(t, canonicalResponse)
	repo := newTestRepo(t, srv)
	got, _ := repo.Search(context.Background(), domain.HotelSearchQuery{PageSize: 200})
	gotIDs := []int64{got.Hotels[0].UniqueID, got.Hotels[1].UniqueID}
	want := []int64{42, 7}
	if !reflect.DeepEqual(gotIDs, want) {
		t.Errorf("order = %v, want %v", gotIDs, want)
	}
}
