// Package postgres adapts PostgreSQL to the domain registry and read ports
// (geo, facilities, chains, inventory lists, jobs).
//
// TODO: wire github.com/jackc/pgx/v5 and run the Flyway-equivalent migrations
// (see the Kotlin module's db/migration scripts).
package postgres

import (
	"context"

	"github.com/knives85/hotel-search/internal/domain"
)

// Repository implements the Postgres-backed read ports.
type Repository struct {
	// TODO: hold the pgx connection pool here.
}

// NewRepository builds a Postgres repository.
func NewRepository() *Repository {
	return &Repository{}
}

// SuggestCities serves the city autocomplete.
func (r *Repository) SuggestCities(ctx context.Context, term string) ([]domain.Suggestion, error) {
	return nil, domain.ErrNotImplemented
}

// SuggestCountries serves the country autocomplete.
func (r *Repository) SuggestCountries(ctx context.Context, term string) ([]domain.Suggestion, error) {
	return nil, domain.ErrNotImplemented
}

// SuggestFacilities serves the facility autocomplete.
func (r *Repository) SuggestFacilities(ctx context.Context, term string) ([]domain.Suggestion, error) {
	return nil, domain.ErrNotImplemented
}

// SuggestChains serves the hotel-chain autocomplete.
func (r *Repository) SuggestChains(ctx context.Context, term string) ([]domain.Suggestion, error) {
	return nil, domain.ErrNotImplemented
}

// ListInventoryLists returns the saved inventory lists.
func (r *Repository) ListInventoryLists(ctx context.Context) ([]domain.InventoryList, error) {
	return nil, domain.ErrNotImplemented
}

// ListJobs returns the background jobs.
func (r *Repository) ListJobs(ctx context.Context) ([]domain.Job, error) {
	return nil, domain.ErrNotImplemented
}

// Compile-time checks that Repository satisfies the intended ports.
var (
	_ domain.GeoLocationReadPort   = (*Repository)(nil)
	_ domain.FacilityReadPort      = (*Repository)(nil)
	_ domain.ChainReadPort         = (*Repository)(nil)
	_ domain.InventoryListReadPort = (*Repository)(nil)
	_ domain.JobReadPort           = (*Repository)(nil)
)
