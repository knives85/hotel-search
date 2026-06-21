# Hotel Search — app module (Go)

A Go reimplementation of the **`app`** module of Hotel Search: the
read/query side that powers hotel search, filtering and faceted badge counts
over **OpenSearch**, with a **PostgreSQL** companion for registries (geo,
facilities, chains), saved inventory lists and jobs, and **S3** for job
artifacts. The existing HTMX UI is meant to stay unchanged — the HTTP routes
and URLs mirror the original controllers.

> This is a learning/portfolio port. It is intentionally a skeleton: the routes
> are wired but every handler returns `501 Not Implemented` until its adapter is
> built.

## Architecture

Hexagonal (ports & adapters), mirroring the Kotlin module:

```
cmd/hotel-search/           entrypoint (config + HTTP server)
internal/
  config/                   configuration from env (defaults from application.yml)
  domain/                   core types + port interfaces (no framework deps)
  usecase/                  thin application logic, depends only on ports
  adapter/
    opensearch/             SearchPort + HotelStatsPort  (TODO: opensearch-go/v4)
    postgres/               geo / facility / chain / inventory / jobs (TODO: pgx)
    s3/                     JobArtifactReadPort           (TODO: aws-sdk-go-v2)
    web/                    HTTP router, handlers, templates (HTMX)
```

The dependency rule points inward: `adapter` and `usecase` depend on `domain`,
never the reverse.

## Routes

Registered under the context path (default `/hotel-search`):

| Method + path                                   | Purpose                       |
| ----------------------------------------------- | ----------------------------- |
| `GET /hotels`                                   | search page                   |
| `GET /hotels/results`                           | results table (HTMX)          |
| `GET /hotels/stats`                             | aggregate stats (HTMX)        |
| `GET /hotels/filter-counts`                     | sidebar badge counts (HTMX)   |
| `GET /hotels/export.csv`                         | CSV export                    |
| `GET /hotels/{country,city,…}-suggest`          | autocomplete suggesters       |
| `GET /jobs`, `/{id}/row`, `/{id}/download`      | background jobs               |
| `GET /healthz`                                  | liveness probe (no prefix)    |

## Run

```bash
make run          # go run ./cmd/hotel-search
make test         # go test ./...
make build        # binary into ./bin

curl -i localhost:8080/healthz                       # 200 ok
curl -i localhost:8080/hotel-search/hotels           # 501 (route wired, not implemented)
```

Configuration is read from the environment (see `internal/config`), e.g.
`PM_ADDR`, `PM_CONTEXT_PATH`, `PM_POSTGRES_DSN`, `PM_OPENSEARCH_ENDPOINT`,
`PM_OPENSEARCH_REGION`, `PM_S3_BUCKET`.

## Roadmap

1. Port the `hotels` read path first: domain `HotelSearchQuery`, the OpenSearch
   query + `terms` aggregations, the results/stats/filter-counts handlers.
2. Add integration tests with `testcontainers-go` (OpenSearch + Postgres).
3. Port the Postgres registries and the suggesters.
4. Port the HTMX templates (`internal/adapter/web/templates`).
5. Optional: re-run a small latency benchmark vs the Kotlin version.

## Note

The module path in `go.mod` is `github.com/knives85/hotel-search`. Update it to
your real repo if needed, e.g.:

```bash
go mod edit -module github.com/<you>/<repo>
grep -rl github.com/knives85/hotel-search . | xargs sed -i '' \
  -e 's#github.com/knives85/hotel-search#github.com/<you>/<repo>#g'
```
