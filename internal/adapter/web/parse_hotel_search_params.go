package web

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ErrInvalidParam wraps any failure to coerce a query-string value into its
// expected Go type. The handler layer can map this to HTTP 400. Use
// errors.Is(err, ErrInvalidParam) to detect.
var ErrInvalidParam = errors.New("invalid query parameter")

// parseHotelSearchParams extracts the typed inputs of the /hotels endpoints
// from r.URL.Query() (or any url.Values). Coercion-only: no domain-level
// normalisation happens here — that is buildQuery's job.
//
// Multi-valued parameters (lists of codes) accept both the repeated-key form
// `countryCodes=IT&countryCodes=FR` and the comma-separated form
// `countryCodes=IT,FR`; both are flattened, trimmed and stripped of empty
// entries.
func parseHotelSearchParams(v url.Values) (HotelSearchParams, error) {
	uniqueID, err := optionalInt64(v, "uniqueId")
	if err != nil {
		return HotelSearchParams{}, err
	}
	sellStatus, err := optionalBool(v, "sellStatus")
	if err != nil {
		return HotelSearchParams{}, err
	}
	reviewScoreMin, err := optionalInt(v, "reviewScoreMin")
	if err != nil {
		return HotelSearchParams{}, err
	}
	reviewScoreMax, err := optionalInt(v, "reviewScoreMax")
	if err != nil {
		return HotelSearchParams{}, err
	}
	numberOfReviewsMin, err := optionalInt(v, "numberOfReviewsMin")
	if err != nil {
		return HotelSearchParams{}, err
	}
	numberOfReviewsMax, err := optionalInt(v, "numberOfReviewsMax")
	if err != nil {
		return HotelSearchParams{}, err
	}
	page, err := intOrDefault(v, "page", 0)
	if err != nil {
		return HotelSearchParams{}, err
	}

	return HotelSearchParams{
		UniqueID: uniqueID,

		HotelName:  v.Get("hotelName"),
		SellStatus: sellStatus,

		StarRatings: stringList(v["starRatings"]),
		Types:       stringList(v["types"]),

		CountryCodes:         stringList(v["countryCodes"]),
		CityCodes:            stringList(v["cityCodes"]),
		RegionCodes:          stringList(v["regionCodes"]),
		TouristicRegionCodes: stringList(v["touristicRegionCodes"]),
		NonAdminCityCodes:    stringList(v["nonAdminCityCodes"]),
		PoiCodes:             stringList(v["poiCodes"]),
		NeighbourhoodCodes:   stringList(v["neighbourhoodCodes"]),

		ChainCodes:    stringList(v["chainCodes"]),
		FacilityCodes: stringList(v["facilityCodes"]),
		BadgeCodes:    stringList(v["badgeCodes"]),

		ReviewScoreMin:     reviewScoreMin,
		ReviewScoreMax:     reviewScoreMax,
		NumberOfReviewsMin: numberOfReviewsMin,
		NumberOfReviewsMax: numberOfReviewsMax,

		CreationDateFrom: v.Get("creationDateFrom"),
		CreationDateTo:   v.Get("creationDateTo"),

		Page: page,
	}, nil
}

// stringList flattens a list of query-string values, splitting comma-separated
// entries, trimming whitespace and dropping empties. Returns nil for empty
// input so the caller can distinguish "absent" from "present-but-empty".
func stringList(vs []string) []string {
	if len(vs) == 0 {
		return nil
	}
	out := make([]string, 0, len(vs))
	for _, v := range vs {
		for _, part := range strings.Split(v, ",") {
			p := strings.TrimSpace(part)
			if p != "" {
				out = append(out, p)
			}
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func optionalInt(v url.Values, key string) (*int, error) {
	raw := strings.TrimSpace(v.Get(key))
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %s=%q must be an integer", ErrInvalidParam, key, raw)
	}
	return &n, nil
}

func optionalInt64(v url.Values, key string) (*int64, error) {
	raw := strings.TrimSpace(v.Get(key))
	if raw == "" {
		return nil, nil
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: %s=%q must be a 64-bit integer", ErrInvalidParam, key, raw)
	}
	return &n, nil
}

func optionalBool(v url.Values, key string) (*bool, error) {
	raw := strings.TrimSpace(v.Get(key))
	if raw == "" {
		return nil, nil
	}
	b, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("%w: %s=%q must be true|false", ErrInvalidParam, key, raw)
	}
	return &b, nil
}

func intOrDefault(v url.Values, key string, def int) (int, error) {
	raw := strings.TrimSpace(v.Get(key))
	if raw == "" {
		return def, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %s=%q must be an integer", ErrInvalidParam, key, raw)
	}
	return n, nil
}
