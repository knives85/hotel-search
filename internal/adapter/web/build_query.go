package web

import (
	"math"
	"strings"

	"github.com/knives85/hotel-search/internal/domain"
)

// buildQuery normalises the raw HTTP-bound HotelSearchParams into the domain
// HotelSearchQuery consumed by use cases:
//
//   - blank strings and empty slices become nil ("absent" in the domain),
//   - star ratings outside the closed option set are dropped,
//   - half-open numeric and date ranges are completed with sensible defaults.
//
// Returns domain.ErrInvalidDateFormat (wrapped) when CreationDateFrom/To do
// not match yyyy-MM-dd.
func buildQuery(p HotelSearchParams) (domain.HotelSearchQuery, error) {
	creationDateRange, err := domain.ParseDateRangeEpoch(p.CreationDateFrom, p.CreationDateTo)
	if err != nil {
		return domain.HotelSearchQuery{}, err
	}

	return domain.HotelSearchQuery{
		UniqueID:             p.UniqueID,
		HotelName:            nilIfBlank(p.HotelName),
		SellStatus:           p.SellStatus,
		StarRatings:          nilIfEmpty(filterStarRatings(p.StarRatings)),
		Types:                nilIfEmpty(p.Types),
		CountryCodes:         nilIfEmpty(p.CountryCodes),
		CityCodes:            nilIfEmpty(p.CityCodes),
		RegionCodes:          nilIfEmpty(p.RegionCodes),
		TouristicRegionCodes: nilIfEmpty(p.TouristicRegionCodes),
		NonAdminCityCodes:    nilIfEmpty(p.NonAdminCityCodes),
		ChainCodes:           nilIfEmpty(p.ChainCodes),
		FacilityCodes:        nilIfEmpty(p.FacilityCodes),
		BadgeCodes:           nilIfEmpty(p.BadgeCodes),
		PoiCodes:             nilIfEmpty(p.PoiCodes),
		NeighbourhoodCodes:   nilIfEmpty(p.NeighbourhoodCodes),
		ReviewScoreRange:     buildIntRange(p.ReviewScoreMin, p.ReviewScoreMax, 0, 100),
		NumberOfReviewsRange: buildIntRange(p.NumberOfReviewsMin, p.NumberOfReviewsMax, 0, math.MaxInt),
		CreationDateRange:    creationDateRange,
		Page:                 p.Page,
		PageSize:             domain.PageSize,
	}, nil
}

func nilIfBlank(s string) *string {
	t := strings.TrimSpace(s)
	if t == "" {
		return nil
	}
	return &s
}

func nilIfEmpty[T any](xs []T) []T {
	if len(xs) == 0 {
		return nil
	}
	return xs
}

func filterStarRatings(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	allowed := make(map[string]struct{}, len(domain.StarRatingOptions))
	for _, o := range domain.StarRatingOptions {
		allowed[o] = struct{}{}
	}
	out := make([]string, 0, len(in))
	for _, r := range in {
		if _, ok := allowed[r]; ok {
			out = append(out, r)
		}
	}
	return out
}

// buildIntRange returns nil unless at least one bound was supplied. Missing
// bounds default to (defaultMin, defaultMax).
func buildIntRange(min, max *int, defaultMin, defaultMax int) *domain.IntRange {
	if min == nil && max == nil {
		return nil
	}
	r := &domain.IntRange{Min: defaultMin, Max: defaultMax}
	if min != nil {
		r.Min = *min
	}
	if max != nil {
		r.Max = *max
	}
	return r
}
