package web

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/knives85/hotel-search/internal/domain"
)

func ptr[T any](v T) *T { return &v }

func TestBuildQuery_EmptyParams(t *testing.T) {
	q, err := buildQuery(HotelSearchParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.Page != 0 {
		t.Errorf("Page = %d, want 0", q.Page)
	}
	if q.PageSize != domain.PageSize {
		t.Errorf("PageSize = %d, want %d", q.PageSize, domain.PageSize)
	}
	// Everything optional should be the zero value.
	want := domain.HotelSearchQuery{Page: 0, PageSize: domain.PageSize}
	if !reflect.DeepEqual(q, want) {
		t.Errorf("query = %+v, want %+v", q, want)
	}
}

func TestBuildQuery_PassesThroughScalarPointers(t *testing.T) {
	p := HotelSearchParams{
		UniqueID:       ptr(int64(42)),
		GiataID:        ptr(int64(7)),
		InternalCityID: ptr(int64(123)),
		SellStatus:     ptr(true),
		Page:           3,
	}
	q, err := buildQuery(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.UniqueID == nil || *q.UniqueID != 42 {
		t.Errorf("UniqueID = %v, want 42", q.UniqueID)
	}
	if q.GiataID == nil || *q.GiataID != 7 {
		t.Errorf("GiataID = %v, want 7", q.GiataID)
	}
	if q.InternalCityID == nil || *q.InternalCityID != 123 {
		t.Errorf("InternalCityID = %v, want 123", q.InternalCityID)
	}
	if q.SellStatus == nil || *q.SellStatus != true {
		t.Errorf("SellStatus = %v, want true", q.SellStatus)
	}
	if q.Page != 3 {
		t.Errorf("Page = %d, want 3", q.Page)
	}
}

func TestBuildQuery_HotelName_BlankBecomesNil(t *testing.T) {
	cases := []string{"", " ", "   ", "\t\n"}
	for _, in := range cases {
		q, err := buildQuery(HotelSearchParams{HotelName: in})
		if err != nil {
			t.Fatalf("input %q: unexpected error: %v", in, err)
		}
		if q.HotelName != nil {
			t.Errorf("input %q: HotelName = %q, want nil", in, *q.HotelName)
		}
	}
}

func TestBuildQuery_HotelName_NonBlankPassedThrough(t *testing.T) {
	q, err := buildQuery(HotelSearchParams{HotelName: "Hilton"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.HotelName == nil || *q.HotelName != "Hilton" {
		t.Errorf("HotelName = %v, want \"Hilton\"", q.HotelName)
	}
}

func TestBuildQuery_StarRatings_FilteredAndEmptyToNil(t *testing.T) {
	t.Run("invalid filtered out, valid kept", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{StarRatings: []string{"ONE", "INVALID", "TWO"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []string{"ONE", "TWO"}
		if !reflect.DeepEqual(q.StarRatings, want) {
			t.Errorf("StarRatings = %v, want %v", q.StarRatings, want)
		}
	})
	t.Run("all invalid -> nil", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{StarRatings: []string{"INVALID", "ALSO_BAD"}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q.StarRatings != nil {
			t.Errorf("StarRatings = %v, want nil", q.StarRatings)
		}
	})
	t.Run("empty input -> nil", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{StarRatings: []string{}})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q.StarRatings != nil {
			t.Errorf("StarRatings = %v, want nil", q.StarRatings)
		}
	})
}

func TestBuildQuery_ListFields_EmptyToNil(t *testing.T) {
	// All "list of codes" fields share the same takeIf-isNotEmpty rule.
	p := HotelSearchParams{
		Types:                []string{},
		CountryCodes:         []string{},
		CityCodes:            []string{},
		RegionCodes:          []string{},
		TouristicRegionCodes: []string{},
		NonAdminCityCodes:    []string{},
		PoiCodes:             []string{},
		NeighbourhoodCodes:   []string{},
		ChainCodes:           []string{},
		FacilityCodes:        []string{},
		BadgeCodes:           []string{},
		ProviderIDs:          []int{},
	}
	q, err := buildQuery(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	checkNil := func(name string, got any) {
		t.Helper()
		v := reflect.ValueOf(got)
		if !v.IsNil() {
			t.Errorf("%s = %v, want nil", name, got)
		}
	}
	checkNil("Types", q.Types)
	checkNil("CountryCodes", q.CountryCodes)
	checkNil("CityCodes", q.CityCodes)
	checkNil("RegionCodes", q.RegionCodes)
	checkNil("TouristicRegionCodes", q.TouristicRegionCodes)
	checkNil("NonAdminCityCodes", q.NonAdminCityCodes)
	checkNil("PoiCodes", q.PoiCodes)
	checkNil("NeighbourhoodCodes", q.NeighbourhoodCodes)
	checkNil("ChainCodes", q.ChainCodes)
	checkNil("FacilityCodes", q.FacilityCodes)
	checkNil("BadgeCodes", q.BadgeCodes)
	checkNil("ProviderIDs", q.ProviderIDs)
}

func TestBuildQuery_ListFields_PassThrough(t *testing.T) {
	p := HotelSearchParams{
		Types:        []string{"Hotels", "Apartments"},
		CountryCodes: []string{"IT", "FR"},
		ProviderIDs:  []int{1, 2, 3},
	}
	q, err := buildQuery(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(q.Types, []string{"Hotels", "Apartments"}) {
		t.Errorf("Types = %v", q.Types)
	}
	if !reflect.DeepEqual(q.CountryCodes, []string{"IT", "FR"}) {
		t.Errorf("CountryCodes = %v", q.CountryCodes)
	}
	if !reflect.DeepEqual(q.ProviderIDs, []int{1, 2, 3}) {
		t.Errorf("ProviderIDs = %v", q.ProviderIDs)
	}
}

func TestBuildQuery_ReviewScoreRange(t *testing.T) {
	t.Run("both nil -> nil range", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q.ReviewScoreRange != nil {
			t.Errorf("ReviewScoreRange = %+v, want nil", q.ReviewScoreRange)
		}
	})
	t.Run("only min -> [min, 100]", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{ReviewScoreMin: ptr(80)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := &domain.IntRange{Min: 80, Max: 100}
		if !reflect.DeepEqual(q.ReviewScoreRange, want) {
			t.Errorf("ReviewScoreRange = %+v, want %+v", q.ReviewScoreRange, want)
		}
	})
	t.Run("only max -> [0, max]", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{ReviewScoreMax: ptr(70)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := &domain.IntRange{Min: 0, Max: 70}
		if !reflect.DeepEqual(q.ReviewScoreRange, want) {
			t.Errorf("ReviewScoreRange = %+v, want %+v", q.ReviewScoreRange, want)
		}
	})
	t.Run("both -> [min, max]", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{
			ReviewScoreMin: ptr(60),
			ReviewScoreMax: ptr(90),
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := &domain.IntRange{Min: 60, Max: 90}
		if !reflect.DeepEqual(q.ReviewScoreRange, want) {
			t.Errorf("ReviewScoreRange = %+v, want %+v", q.ReviewScoreRange, want)
		}
	})
}

func TestBuildQuery_NumberOfReviewsRange_MaxDefault(t *testing.T) {
	t.Run("only min -> [min, MaxInt]", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{NumberOfReviewsMin: ptr(10)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if q.NumberOfReviewsRange == nil ||
			q.NumberOfReviewsRange.Min != 10 ||
			q.NumberOfReviewsRange.Max != math.MaxInt {
			t.Errorf("NumberOfReviewsRange = %+v, want {10, MaxInt}", q.NumberOfReviewsRange)
		}
	})
	t.Run("only max -> [0, max]", func(t *testing.T) {
		q, err := buildQuery(HotelSearchParams{NumberOfReviewsMax: ptr(500)})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := &domain.IntRange{Min: 0, Max: 500}
		if !reflect.DeepEqual(q.NumberOfReviewsRange, want) {
			t.Errorf("NumberOfReviewsRange = %+v, want %+v", q.NumberOfReviewsRange, want)
		}
	})
}

func TestBuildQuery_CreationDateRange_Parsed(t *testing.T) {
	q, err := buildQuery(HotelSearchParams{
		CreationDateFrom: "2024-01-01",
		CreationDateTo:   "2024-01-31",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.CreationDateRange == nil {
		t.Fatal("CreationDateRange = nil, want non-nil")
	}
	if q.CreationDateRange.Min != 1704067200000 {
		t.Errorf("Min = %d, want 1704067200000", q.CreationDateRange.Min)
	}
	if q.CreationDateRange.Max != 1706745599999 {
		t.Errorf("Max = %d, want 1706745599999", q.CreationDateRange.Max)
	}
}

func TestBuildQuery_CreationDateRange_InvalidWrapsError(t *testing.T) {
	_, err := buildQuery(HotelSearchParams{CreationDateFrom: "not-a-date"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, domain.ErrInvalidDateFormat) {
		t.Errorf("error %v is not ErrInvalidDateFormat", err)
	}
}

func TestBuildQuery_ProviderStatusPassedThrough(t *testing.T) {
	ps := domain.ProviderStatusActive
	q, err := buildQuery(HotelSearchParams{ProviderStatus: &ps})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q.ProviderStatus == nil || *q.ProviderStatus != domain.ProviderStatusActive {
		t.Errorf("ProviderStatus = %v, want ACTIVE", q.ProviderStatus)
	}
}
