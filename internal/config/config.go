// Package config loads the app-module configuration from environment
// variables, with defaults that mirror the Kotlin application.yml.
package config

import "os"

// Config holds the runtime configuration of the app module.
type Config struct {
	Addr        string // HTTP listen address, e.g. ":8080"
	ContextPath string // base path, e.g. "/hotel-search"

	// PostgresDSN points at the relational store for registries
	// (geo, facilities, chains), inventory lists and jobs.
	PostgresDSN string

	// OpenSearch backs hotel search, filtering and aggregations.
	OpenSearchEndpoint string
	OpenSearchRegion   string // AWS region for SigV4 signing

	// S3 stores job artifacts (e.g. CSV exports).
	S3Bucket string
	S3Region string
}

// Load reads configuration from the environment, applying defaults.
func Load() Config {
	return Config{
		Addr:               env("PM_ADDR", ":8080"),
		ContextPath:        env("PM_CONTEXT_PATH", "/hotel-search"),
		PostgresDSN:        env("PM_POSTGRES_DSN", "postgres://hotel:hotel@localhost:5432/hotelsearch?sslmode=disable"),
		OpenSearchEndpoint: env("PM_OPENSEARCH_ENDPOINT", "https://localhost:9200"),
		OpenSearchRegion:   env("PM_OPENSEARCH_REGION", "eu-central-1"),
		S3Bucket:           env("PM_S3_BUCKET", "hotel-search-qa"),
		S3Region:           env("PM_S3_REGION", "eu-central-1"),
	}
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}
