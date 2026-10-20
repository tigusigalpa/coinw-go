# Local verification — 2026-10-20

Verified on Windows amd64 with live integrations disabled explicitly.

| Check | Result |
|---|---|
| `go test ./...`, Go 1.26.10 | Passed |
| `go test -race ./...`, Go 1.27.3 | Passed |
| Changed service/HTTP/WS contracts on both series | Passed |
| `go vet ./...` | Passed |
| Staticcheck v0.8.2 | Passed |
| golangci-lint v2.15.0 | 0 issues |
| Module checksum verification | Passed |
| `go mod tidy -diff` | No changes |
| `go build ./...` | Passed |
| Formatting and whitespace | Clean |
| govulncheck v1.9.0 on Go 1.27.3 | No vulnerabilities found |
| Integration tests without opt-in | All four suites skipped |

Contract coverage is 86 REST operations and 26 WS subscription variants, with
86 REST success fixtures and 48 WS event/ack fixtures. Common error fixtures
are exercised for every sendable REST operation. Ambiguous routes are tested to
make zero HTTP requests. Source links and fixture provenance are in
`contracts/catalogue.json`; the readable index is `docs/API.md`.

Windows test runs used serial package builds and stripped test debug symbols
where needed; race instrumentation remained enabled. CI is configured to run
the full race/test/vet/build matrix on Go 1.26.x and 1.27.x, plus pinned lint,
static analysis, checksum and vulnerability checks. Remote CI was not triggered
as part of this local work.

No live exchange integration, order cancellation, trading, withdrawal or
transfer was executed. Live conformance is unverified until an authorized
operator runs the opt-in tests. No tag or remote release was published. The
three documented route ambiguity gates remain in place as required by the
specification. Use `docs/RELEASING.md` before publishing.
