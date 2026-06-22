package opensearch

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/knives85/hotel-search/internal/domain"
)

func TestToHotel_FullDocument(t *testing.T) {
	raw := `{
		"unique_id": 42,
		"index_status": "COMPLETE",
		"hotel_name": "Grand Hotel",
		"sell_status": true,
		"star_rating": "FOUR",
		"type": "Hotels",
		"country": {"code": "IT", "name": "Italy"},
		"city": {"code": "ROM", "name": "Rome"},
		"admin_region": {"code": "LAZ", "name": "Lazio"},
		"touristic_region": {"code": "TUS", "name": "Tuscany"},
		"nonadmin_city": {"code": "NAC1", "name": "Some Town"},
		"neighbourhood": {"code": "NB1", "name": "Trastevere"},
		"chain": {"code": "HIL", "name": "Hilton"},
		"facility_codes": ["WIFI", "POOL"],
		"nearby_pois": [{"code": "P1", "name": "Colosseum"}, {"code": "P2", "name": "Forum"}],
		"badges": ["NEW", "TRENDING"],
		"content_score": 80,
		"review_score": 90,
		"number_of_reviews": 1500,
		"location_score": 85,
		"creation_date": 1704067200000,
		"last_update_date": 1735603200000,
		"coordinates": {"latitude": 41.9028, "longitude": 12.4964}
	}`
	var doc hotelSearchDocument
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	got := toHotel(doc)

	if got.UniqueID != 42 {
		t.Errorf("UniqueID = %d, want 42", got.UniqueID)
	}
	if got.IndexStatus != domain.IndexStatusComplete {
		t.Errorf("IndexStatus = %q, want COMPLETE", got.IndexStatus)
	}
	if got.HotelName == nil || *got.HotelName != "Grand Hotel" {
		t.Errorf("HotelName = %v", got.HotelName)
	}
	if got.SellStatus == nil || *got.SellStatus != true {
		t.Errorf("SellStatus = %v", got.SellStatus)
	}
	if got.StarRating == nil || *got.StarRating != "FOUR" {
		t.Errorf("StarRating = %v", got.StarRating)
	}
	if got.Country == nil || got.Country.Code != "IT" {
		t.Errorf("Country = %+v", got.Country)
	}
	if got.Country.Name == nil || *got.Country.Name != "Italy" {
		t.Errorf("Country.Name = %v", got.Country.Name)
	}
	if got.Chain == nil || got.Chain.Code != "HIL" {
		t.Errorf("Chain = %+v", got.Chain)
	}
	if !reflect.DeepEqual(got.Facilities, []string{"WIFI", "POOL"}) {
		t.Errorf("Facilities = %v", got.Facilities)
	}
	if !reflect.DeepEqual(got.Badges, []string{"NEW", "TRENDING"}) {
		t.Errorf("Badges = %v", got.Badges)
	}
	if len(got.PointsOfInterest) != 2 ||
		got.PointsOfInterest[0].Code != "P1" ||
		got.PointsOfInterest[1].Code != "P2" {
		t.Errorf("PointsOfInterest = %+v", got.PointsOfInterest)
	}
	if got.Coordinates == nil || got.Coordinates.Latitude != 41.9028 || got.Coordinates.Longitude != 12.4964 {
		t.Errorf("Coordinates = %+v", got.Coordinates)
	}
	if got.ReviewScore == nil || *got.ReviewScore != 90 {
		t.Errorf("ReviewScore = %v", got.ReviewScore)
	}
}

func TestToHotel_MinimalDocument(t *testing.T) {
	raw := `{"unique_id": 7, "index_status": "PARTIAL"}`
	var doc hotelSearchDocument
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	got := toHotel(doc)

	if got.UniqueID != 7 {
		t.Errorf("UniqueID = %d", got.UniqueID)
	}
	if got.IndexStatus != domain.IndexStatusPartial {
		t.Errorf("IndexStatus = %q", got.IndexStatus)
	}
	// Every optional must be nil / empty.
	if got.HotelName != nil || got.SellStatus != nil || got.StarRating != nil ||
		got.Country != nil || got.City != nil || got.Chain != nil || got.Coordinates != nil {
		t.Errorf("expected all optionals nil, got: %+v", got)
	}
	if len(got.Facilities) != 0 || len(got.Badges) != 0 || len(got.PointsOfInterest) != 0 {
		t.Errorf("expected empty slices, got facilities=%v badges=%v poi=%v",
			got.Facilities, got.Badges, got.PointsOfInterest)
	}
}

func TestToHotel_UnknownIndexStatus_FallsBackToPartial(t *testing.T) {
	doc := hotelSearchDocument{UniqueID: 1, IndexStatus: "WHATEVER"}
	if toHotel(doc).IndexStatus != domain.IndexStatusPartial {
		t.Errorf("unknown index_status should fall back to PARTIAL")
	}
}

func TestToHotel_Coordinates_PartialAreDropped(t *testing.T) {
	// If one coordinate is missing, the document mapper drops the whole pair.
	raw := `{"unique_id": 1, "index_status": "COMPLETE", "coordinates": {"latitude": 41.9}}`
	var doc hotelSearchDocument
	_ = json.Unmarshal([]byte(raw), &doc)
	got := toHotel(doc)
	if got.Coordinates != nil {
		t.Errorf("Coordinates = %+v, want nil for partial pair", got.Coordinates)
	}
}
