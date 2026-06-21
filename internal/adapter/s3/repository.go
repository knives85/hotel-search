// Package s3 adapts object storage (S3) to the job-artifact read port.
//
// TODO: wire github.com/aws/aws-sdk-go-v2 (config + service/s3).
package s3

import (
	"context"

	"github.com/knives85/hotel-search/internal/domain"
)

// Repository implements the S3-backed job-artifact read port.
type Repository struct {
	// TODO: hold the S3 client and the bucket name here.
}

// NewRepository builds an S3 repository.
func NewRepository() *Repository {
	return &Repository{}
}

// Download fetches a job's output artifact by id.
func (r *Repository) Download(ctx context.Context, id string) ([]byte, error) {
	return nil, domain.ErrNotImplemented
}

// Compile-time check that Repository satisfies the intended port.
var _ domain.JobArtifactReadPort = (*Repository)(nil)
