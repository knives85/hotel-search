package web

import (
	"errors"
	"net/url"
	"reflect"
	"testing"
)

func TestParseHotelSearchParams_Empty(t *testing.T) {
	p, err := parseHotelSearchParams(url.Values{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(p, HotelSearchParams{}) {
		t.Errorf("params = %+v, want zero value", p)
	}
}

func TestParseHotelSearchParams_UniqueID(t *testing.T) {
	v := url.Values{"uniqueId": {"42"}}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.UniqueID == nil || *p.UniqueID != 42 {
		t.Errorf("UniqueID = %v, want 42", p.UniqueID)
	}
}

func TestParseHotelSearchParams_UniqueID_BlankIgnored(t *testing.T) {
	v := url.Values{"uniqueId": {""}}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.UniqueID != nil {
		t.Errorf("UniqueID = %v, want nil for blank", p.UniqueID)
	}
}

func TestParseHotelSearchParams_UniqueID_InvalidReturnsError(t *testing.T) {
	v := url.Values{"uniqueId": {"not-a-number"}}
	_, err := parseHotelSearchParams(v)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidParam) {
		t.Errorf("error %v is not ErrInvalidParam", err)
	}
}

func TestParseHotelSearchParams_HotelName(t *testing.T) {
	v := url.Values{"hotelName": {"Hilton"}}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.HotelName != "Hilton" {
		t.Errorf("HotelName = %q, want %q", p.HotelName, "Hilton")
	}
}

func TestParseHotelSearchParams_SellStatus(t *testing.T) {
	cases := map[string]bool{"true": true, "false": false}
	for in, want := range cases {
		v := url.Values{"sellStatus": {in}}
		p, err := parseHotelSearchParams(v)
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", in, err)
		}
		if p.SellStatus == nil || *p.SellStatus != want {
			t.Errorf("input %q: SellStatus = %v, want %v", in, p.SellStatus, want)
		}
	}
}

func TestParseHotelSearchParams_SellStatus_InvalidReturnsError(t *testing.T) {
	v := url.Values{"sellStatus": {"yes"}}
	_, err := parseHotelSearchParams(v)
	if !errors.Is(err, ErrInvalidParam) {
		t.Errorf("err = %v, want ErrInvalidParam", err)
	}
}

func TestParseHotelSearchParams_StringLists(t *testing.T) {
	// Multi-valued params come as multiple values OR a comma-separated single
	// value — both shapes must be accepted (HTMX/Thymeleaf forms emit both).
	t.Run("multi value", func(t *testing.T) {
		v := url.Values{"countryCodes": {"IT", "FR", "DE"}}
		p, _ := parseHotelSearchParams(v)
		want := []string{"IT", "FR", "DE"}
		if !reflect.DeepEqual(p.CountryCodes, want) {
			t.Errorf("CountryCodes = %v, want %v", p.CountryCodes, want)
		}
	})
	t.Run("comma separated single value", func(t *testing.T) {
		v := url.Values{"countryCodes": {"IT,FR,DE"}}
		p, _ := parseHotelSearchParams(v)
		want := []string{"IT", "FR", "DE"}
		if !reflect.DeepEqual(p.CountryCodes, want) {
			t.Errorf("CountryCodes = %v, want %v", p.CountryCodes, want)
		}
	})
	t.Run("mixed empty entries are dropped", func(t *testing.T) {
		v := url.Values{"countryCodes": {"IT,,FR", ""}}
		p, _ := parseHotelSearchParams(v)
		want := []string{"IT", "FR"}
		if !reflect.DeepEqual(p.CountryCodes, want) {
			t.Errorf("CountryCodes = %v, want %v", p.CountryCodes, want)
		}
	})
}

func TestParseHotelSearchParams_AllCodeLists(t *testing.T) {
	v := url.Values{
		"starRatings":          {"ONE", "TWO"},
		"types":                {"Hotels"},
		"cityCodes":             {"ROM"},
		"regionCodes":           {"LAZ"},
		"touristicRegionCodes":  {"TUS"},
		"nonAdminCityCodes":     {"NAC1"},
		"poiCodes":              {"POI1"},
		"neighbourhoodCodes":    {"NB1"},
		"chainCodes":            {"HIL"},
		"facilityCodes":         {"WIFI"},
		"badgeCodes":            {"NEW"},
	}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checks := map[string]struct{ got, want []string }{
		"StarRatings":          {p.StarRatings, []string{"ONE", "TWO"}},
		"Types":                {p.Types, []string{"Hotels"}},
		"CityCodes":            {p.CityCodes, []string{"ROM"}},
		"RegionCodes":          {p.RegionCodes, []string{"LAZ"}},
		"TouristicRegionCodes": {p.TouristicRegionCodes, []string{"TUS"}},
		"NonAdminCityCodes":    {p.NonAdminCityCodes, []string{"NAC1"}},
		"PoiCodes":             {p.PoiCodes, []string{"POI1"}},
		"NeighbourhoodCodes":   {p.NeighbourhoodCodes, []string{"NB1"}},
		"ChainCodes":           {p.ChainCodes, []string{"HIL"}},
		"FacilityCodes":        {p.FacilityCodes, []string{"WIFI"}},
		"BadgeCodes":           {p.BadgeCodes, []string{"NEW"}},
	}
	for name, c := range checks {
		if !reflect.DeepEqual(c.got, c.want) {
			t.Errorf("%s = %v, want %v", name, c.got, c.want)
		}
	}
}

func TestParseHotelSearchParams_IntPointers(t *testing.T) {
	v := url.Values{
		"reviewScoreMin":     {"60"},
		"reviewScoreMax":     {"90"},
		"numberOfReviewsMin": {"10"},
		"numberOfReviewsMax": {"500"},
	}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ReviewScoreMin == nil || *p.ReviewScoreMin != 60 {
		t.Errorf("ReviewScoreMin = %v", p.ReviewScoreMin)
	}
	if p.ReviewScoreMax == nil || *p.ReviewScoreMax != 90 {
		t.Errorf("ReviewScoreMax = %v", p.ReviewScoreMax)
	}
	if p.NumberOfReviewsMin == nil || *p.NumberOfReviewsMin != 10 {
		t.Errorf("NumberOfReviewsMin = %v", p.NumberOfReviewsMin)
	}
	if p.NumberOfReviewsMax == nil || *p.NumberOfReviewsMax != 500 {
		t.Errorf("NumberOfReviewsMax = %v", p.NumberOfReviewsMax)
	}
}

func TestParseHotelSearchParams_IntInvalidReturnsError(t *testing.T) {
	v := url.Values{"reviewScoreMin": {"abc"}}
	_, err := parseHotelSearchParams(v)
	if !errors.Is(err, ErrInvalidParam) {
		t.Errorf("err = %v, want ErrInvalidParam", err)
	}
}

func TestParseHotelSearchParams_CreationDateRaw(t *testing.T) {
	v := url.Values{
		"creationDateFrom": {"2024-01-01"},
		"creationDateTo":   {"2024-01-31"},
	}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.CreationDateFrom != "2024-01-01" {
		t.Errorf("CreationDateFrom = %q", p.CreationDateFrom)
	}
	if p.CreationDateTo != "2024-01-31" {
		t.Errorf("CreationDateTo = %q", p.CreationDateTo)
	}
}

func TestParseHotelSearchParams_Page_Default(t *testing.T) {
	p, _ := parseHotelSearchParams(url.Values{})
	if p.Page != 0 {
		t.Errorf("Page = %d, want 0", p.Page)
	}
}

func TestParseHotelSearchParams_Page_Explicit(t *testing.T) {
	v := url.Values{"page": {"3"}}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Page != 3 {
		t.Errorf("Page = %d, want 3", p.Page)
	}
}

func TestParseHotelSearchParams_FullRoundTripThroughBuildQuery(t *testing.T) {
	// Smoke test: parse → buildQuery should succeed and reflect the inputs.
	v := url.Values{
		"hotelName":      {"Marriott"},
		"starRatings":    {"FOUR", "FIVE"},
		"countryCodes":   {"IT"},
		"sellStatus":     {"true"},
		"reviewScoreMin": {"70"},
		"page":           {"2"},
	}
	p, err := parseHotelSearchParams(v)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	q, err := buildQuery(p)
	if err != nil {
		t.Fatalf("buildQuery error: %v", err)
	}
	if q.HotelName == nil || *q.HotelName != "Marriott" {
		t.Errorf("HotelName = %v", q.HotelName)
	}
	if !reflect.DeepEqual(q.StarRatings, []string{"FOUR", "FIVE"}) {
		t.Errorf("StarRatings = %v", q.StarRatings)
	}
	if !reflect.DeepEqual(q.CountryCodes, []string{"IT"}) {
		t.Errorf("CountryCodes = %v", q.CountryCodes)
	}
	if q.SellStatus == nil || *q.SellStatus != true {
		t.Errorf("SellStatus = %v", q.SellStatus)
	}
	if q.ReviewScoreRange == nil || q.ReviewScoreRange.Min != 70 || q.ReviewScoreRange.Max != 100 {
		t.Errorf("ReviewScoreRange = %+v", q.ReviewScoreRange)
	}
	if q.Page != 2 {
		t.Errorf("Page = %d", q.Page)
	}
}
