# Contributing

Use Go 1.26 or 1.27 with the latest available patch. The `go 1.26.0` directive
selects the oldest maintained release series on the 2026-10-20 documentation
snapshot. CI tests both maintained series; older Go versions are unsupported.
Review this policy at each release using the
[Go support policy](https://go.dev/doc/devel/release#policy) and
[download list](https://go.dev/dl/). New dependencies require a module, license,
maintenance and necessity review.

The only runtime dependency is `github.com/gorilla/websocket v1.5.3`, under
BSD-2-Clause. It supplies the RFC 6455 transport; our separate raw JSON and
Socket.IO/Engine.IO adapters implement the CoinW protocols. Its reader/writer
constraints are enforced by the session wrapper. It is not a Socket.IO client
by itself. See the [upstream release](https://github.com/gorilla/websocket/releases/tag/v1.5.3)
and [license](https://github.com/gorilla/websocket/blob/v1.5.3/LICENSE).

## Contract changes

1. Read the linked English CoinW page and changelog. Record conflicting table,
   endpoint and example statements in `docs/CONTRACTS.md`.
2. Update the endpoint, typed request/response or stream contract, source link
   and fixture mapping in `contracts/catalogue.json` together. Every catalogue
   row needs a dedicated fixture and an exercised public method/subscription.
3. Use sanitized published examples or explicitly label a synthetic fixture.
   Remove omission markers and repair printed literal syntax explicitly; never
   capture live private responses or credentials. Preserve numeric tokens and
   decimal scale without parsing through `float64`.
4. Add a meaningful regression for the changed behavior, including errors and
   cancellation where relevant. Do not silently relax a schema or invent an
   exchange rule. Keep unresolved routes gated and preserve raw conflict forms.
5. Update migration notes and the changelog for public API changes.

Request DTOs expose all published filters; optional pointer fields distinguish
omission and zero. `services/base.go` validates the request contract before the
transport. Response DTOs validate known fields and required structure. New
upstream fields require an explicit contract update. Monetary models must use
`Decimal` or a named exact compatibility adapter, never binary floats.

## Local checks

```sh
gofmt -w .
go test ./...
go test -race -timeout 2m ./...
go vet ./...
go run honnef.co/go/tools/cmd/staticcheck@v0.8.2 ./...
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.15.0 run
go mod verify
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

Run on a current patched toolchain so security scans include accurate standard
library results. All unit/contract tests use local fake peers or `httptest`.
See `docs/INTEGRATION.md` for optional live checks. Never add credentials to CI,
enable live mutation tests in PR/release jobs, or run trade/withdraw/transfer
actions merely to refresh fixtures. Inject shared limiters, clocks, retry and
WS transport dependencies; supplied implementations must be concurrency safe.

Follow `CODE_OF_CONDUCT.md` and `docs/RELEASING.md` before publishing a release.
