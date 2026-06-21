# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this repo is

A Go reimplementation of the **`app`** module of Hotel Search (originally Kotlin/Spring Boot). It is the read/query side that powers hotel search, filtering and faceted badge counts over **OpenSearch**, with a **PostgreSQL** companion for registries (geo, facilities, chains), saved inventory lists and jobs, and **S3** for job artifacts.

The existing HTMX UI from the Kotlin module is meant to stay unchanged — the HTTP routes and URLs mirror the original Kotlin controllers exactly so the same fragments/swaps keep working once handlers are implemented.

This is a skeleton port: routes are wired but most handlers return `501 Not Implemented` until their adapter is built. Treat `domain.ErrNotImplemented` and the `// TODO:` markers as the work-list.

> Project-specific note: the global `~/.claude/CLAUDE.md` defaults (Kotlin, Spring Boot, Mockk/Kotest, Maven, spotless) do **not** apply here — this module is Go. Use `go fmt`/`go vet`, standard `testing` package, `httptest`, and (per roadmap) `testcontainers-go` for integration tests.

## Common commands

```bash
make run          # go run ./cmd/hotel-search
make test         # go test ./...
make build        # binary into ./bin
make fmt          # go fmt ./...
make vet          # go vet ./...
make tidy         # go mod tidy

# Run a single test
go test ./internal/adapter/web -run TestHealthz -v

# Smoke check (server on :8080 by default)
curl -i localhost:8080/healthz                       # 200 ok
curl -i localhost:8080/hotel-search/hotels           # 501 until ported
```

Module path: `github.com/knives85/hotel-search`. Go version: `1.22` (uses the method-prefixed `ServeMux` patterns like `"GET /healthz"` and `r.PathValue("id")`, both 1.22+).

## Architecture

Hexagonal (ports & adapters). The dependency rule points inward: `adapter` and `usecase` depend on `domain`, **never the reverse**.

```
cmd/hotel-search/           entrypoint (config + HTTP server + graceful shutdown)
internal/
  config/                   PM_* env vars, defaults mirror the Kotlin application.yml
  domain/                   models + port interfaces, no framework deps
  usecase/                  thin orchestrators, depend only on domain ports
  adapter/
    opensearch/             SearchPort + HotelStatsPort  (target: opensearch-go/v4 + SigV4)
    postgres/               geo / facility / chain / inventory / jobs (target: jackc/pgx/v5)
    s3/                     JobArtifactReadPort           (target: aws-sdk-go-v2)
    web/                    net/http ServeMux, handlers, templates (HTMX)
```

Key wiring conventions to preserve:

- Each adapter file ends with compile-time port-conformance assertions like `var _ domain.SearchPort = (*Repository)(nil)`. Keep these when adding/removing methods so port mismatches fail at build time.
- Use cases are intentionally thin (one struct per route, holding the ports it needs, with a single `Execute`). Don't bake business logic into handlers or adapters.
- Routes are registered in `internal/adapter/web/server.go` under `s.contextPath` (default `/hotel-search`), except `/healthz` which is intentionally **outside** the prefix so probes don't depend on it.
- Handler stubs call `notImplemented(w, name)` — replace these per-handler as their adapters land, mirroring the Kotlin controller's response shape (HTML fragment for HTMX swaps, JSON only where the Kotlin endpoint returned JSON).

## Templates / UI

The UI is server-rendered HTMX (Thymeleaf in Kotlin). To keep it unchanged, port the HTML fragments into `internal/adapter/web/templates/` and render with Go's `html/template`. The HTML barely changes — mostly translating Thymeleaf `th:*` attributes into `{{ ... }}` actions. See `internal/adapter/web/templates/README.md` for the Kotlin→Go fragment mapping.

## Configuration

All config comes from env vars (see `internal/config/config.go`):
`PM_ADDR`, `PM_CONTEXT_PATH`, `PM_POSTGRES_DSN`, `PM_OPENSEARCH_ENDPOINT`, `PM_OPENSEARCH_REGION`, `PM_S3_BUCKET`, `PM_S3_REGION`. Defaults are dev-friendly (localhost Postgres, localhost OpenSearch, `eu-central-1`).

## Notes for porting

- **The `usecase` package has no test file yet.** When porting a use case, add a table-driven test alongside it; the existing `internal/adapter/web/server_test.go` is the style reference (plain `testing` + `httptest`, no third-party assertion libs).
- The parent working directory is still named `portfolio-manager-go/` from an earlier rename — this does not affect the Go module (which is `github.com/knives85/hotel-search`) but the user may want to rename it externally at some point.

## Porting roadmap (from README)

1. `hotels` read path first: domain `HotelSearchQuery`, OpenSearch query + `terms` aggregations, results/stats/filter-counts handlers.
2. Integration tests with `testcontainers-go` (OpenSearch + Postgres).
3. Postgres registries + suggesters.
4. HTMX templates.
5. Optional latency benchmark vs the Kotlin version.

When porting a route, the loop is: extend `domain` types/ports → add/extend a `usecase` → implement the adapter method → wire concrete adapter into `main.go` → replace the `notImplemented` call with the real handler + template render. Keep handlers thin: parse query → call use case → render fragment.
