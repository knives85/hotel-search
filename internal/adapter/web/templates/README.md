# Templates

The existing UI is server-rendered HTMX (Thymeleaf in the Kotlin module).
To keep the UI unchanged, port the HTML fragments here and render them with
Go's standard `html/template`.

Mapping from the Kotlin module (`app/src/main/resources/templates/`):

- `hotels/index.html`              → full search page
- `hotels/fragments/results-*.html`→ results table / rows (HTMX swaps)
- `hotels/fragments/*-suggestions.html` → autocomplete dropdowns
- `hotels/fragments/stats-oob.html`, `sidebar-counts-oob.html` → OOB swaps
- `inventory-lists/…`, `jobs/…`    → respective pages and fragments

The HTML itself barely changes: HTMX only needs the same endpoints to return
the same markup. Mostly it is translating Thymeleaf `th:*` attributes into
`{{ ... }}` template actions.
