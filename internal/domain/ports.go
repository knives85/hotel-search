package domain

import "context"

// SearchPort is the read side of hotel search, backed by OpenSearch.
type SearchPort interface {
	Search(ctx context.Context, q HotelSearchQuery) (SearchResult, error)
	SidebarFilterCounts(ctx context.Context, q HotelSearchQuery) (SidebarFilterCounts, error)
}

// HotelStatsPort exposes the aggregate stats, backed by OpenSearch.
type HotelStatsPort interface {
	Stats(ctx context.Context, q HotelSearchQuery) (HotelStats, error)
}

// GeoLocationReadPort serves geo autocomplete (cities, countries, regions,
// neighbourhoods, POIs), backed by Postgres.
//
// TODO: add the remaining suggesters (neighbourhood, admin/touristic region,
// non-admin city, POI) as they are ported.
type GeoLocationReadPort interface {
	SuggestCities(ctx context.Context, term string) ([]Suggestion, error)
	SuggestCountries(ctx context.Context, term string) ([]Suggestion, error)
}

// FacilityReadPort serves facility autocomplete, backed by Postgres.
type FacilityReadPort interface {
	SuggestFacilities(ctx context.Context, term string) ([]Suggestion, error)
}

// ChainReadPort serves hotel-chain autocomplete, backed by Postgres.
type ChainReadPort interface {
	SuggestChains(ctx context.Context, term string) ([]Suggestion, error)
}

// InventoryListReadPort reads saved inventory lists, backed by Postgres.
type InventoryListReadPort interface {
	ListInventoryLists(ctx context.Context) ([]InventoryList, error)
}

// JobReadPort reads background jobs, backed by Postgres.
type JobReadPort interface {
	ListJobs(ctx context.Context) ([]Job, error)
}

// JobArtifactReadPort downloads a job's output artifact, backed by S3.
type JobArtifactReadPort interface {
	Download(ctx context.Context, id string) ([]byte, error)
}
