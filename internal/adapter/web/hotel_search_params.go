package web

// HotelSearchParams is the raw, typed input collected from the /hotels query
// string. Nil pointers / empty strings / nil slices mean "absent". The
// translation to the normalised domain.HotelSearchQuery is buildQuery's job.
type HotelSearchParams struct {
	UniqueID *int64

	HotelName  string // empty/blank == absent
	SellStatus *bool

	StarRatings []string
	Types       []string

	CountryCodes         []string
	CityCodes            []string
	RegionCodes          []string
	TouristicRegionCodes []string
	NonAdminCityCodes    []string
	PoiCodes             []string
	NeighbourhoodCodes   []string

	ChainCodes    []string
	FacilityCodes []string
	BadgeCodes    []string

	ReviewScoreMin     *int
	ReviewScoreMax     *int
	NumberOfReviewsMin *int
	NumberOfReviewsMax *int

	CreationDateFrom string // raw yyyy-MM-dd or empty
	CreationDateTo   string

	Page int
}
