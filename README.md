# CoinW Go SDK

`github.com/tigusigalpa/coinw-go` — a Go SDK for the CoinW exchange API
covering the **entire documented surface**: Spot & Futures REST, Common
(accounts, sub-accounts, affiliate), and all five WebSocket variants.

Documentation source of truth: <https://www.coinw.com/api-doc/en> (English
sitemap snapshot 2026-10-20). Nothing is guessed — endpoints whose official
documentation is internally contradictory are registered but **gated**:
calling them returns `errors.AmbiguousDocumentationError` instead of
sending an unverified request.

## Install

```bash
go get github.com/tigusigalpa/coinw-go
```

Requires Go 1.21+. Only dependency: `gorilla/websocket`.

## Quick start

```go
package main

import (
    "context"
    "fmt"
    "os"

    coinw "github.com/tigusigalpa/coinw-go"
)

func main() {
    c, err := coinw.NewClient(
        coinw.WithSpotCredentials(os.Getenv("COINW_SPOT_API_KEY"), os.Getenv("COINW_SPOT_SECRET_KEY")),
        coinw.WithFuturesCredentials(os.Getenv("COINW_FUTURES_API_KEY"), os.Getenv("COINW_FUTURES_SECRET_KEY")),
        coinw.WithUID(os.Getenv("COINW_UID")), // optional — UID-scoped limiter buckets
    )
    if err != nil { panic(err) }

    ctx := context.Background()
    res, err := c.SpotMarket().GetRecentTrades(ctx, "BTC_USDT", 0)
    if err != nil { panic(err) }
    fmt.Println(string(res.Data))
}
```

## Architecture

| Package | Role |
|---|---|
| `coinw` (root) | `Client` entry point, `ClientOption`s, `Decimal`, WS factories |
| `transport` | request pipeline: ambiguity gate → limiter → signer → HTTP → retry → validate; `Result` envelope |
| `endpoints` | registry of every documented route with method/path/auth/rate-group/ambiguity metadata |
| `auth` | `SpotSigner` (MD5) + `FuturesSigner` (HMAC-SHA256) — isolated profiles |
| `ratelimit` | documented global groups + per-interface token buckets, ctx-aware |
| `services` | typed REST clients per domain |
| `websocket` | raw-JSON sessions (Spot M2 / private / Futures), Socket.IO session (Spot M1), `DepthSequenceTracker`, reconnect |
| `errors` | typed error taxonomy with secret redaction |

### Service map

| `Client` method | Domain |
|---|---|
| `SpotMarket()` `SpotAccount()` `SpotCheck()` `SpotTrade()` | Spot REST |
| `FuturesMarket()` `FuturesAssets()` `FuturesCommon()` `FuturesCheck()` `FuturesPosition()` `FuturesTrade()` | Futures REST |
| `Account()` `SubAccounts()` `Affiliate()` | Common REST |
| `SpotPublicWS()` `SpotPrivateWS()` `SpotSocketIOWS()` | Spot WebSocket |
| `FuturesPublicWS()` `FuturesPrivateWS()` | Futures WebSocket |

## Authentication

- **Spot private REST**: `api_key` + sorted `k=v&` pairs + `secret_key=…`,
  `sign = UPPER(MD5(…))` — everything rides in the URL.
- **Futures private REST**: `payload = ts + METHOD + path + (?query|body)`,
  `sign = base64(HMAC-SHA256(secret, payload))`, headers `sign`/`api_key`/`timestamp`.
- **WS private**: `{"event":"login","params":{"api_key":…,"passphrase":…}}`
  (the field is literally `passphrase` per docs).
- **Spot Socket.IO public**: token from `GET …/pusher/public-token`, then
  `wss://ws.futurescw.info?token=…` with Engine.IO/Socket.IO framing
  (`42["subscribe",{"args":"spot/match:BTC-USDT"}]`).

## WebSocket

```go
// Futures public stream
ws := c.FuturesPublicWS()
if err := ws.Connect(ctx); err != nil { panic(err) }
_ = ws.Subscribe(map[string]interface{}{
    "biz": "perpetual", "type": "depth", "pairCode": "BTC_USDT",
})
_ = ws.Run(ctx,
    func(m *ws.Message) bool {
        fmt.Println(m.Type, string(m.Data))
        return true // false stops the loop
    },
    func(err error) bool {
        fmt.Println("ws error:", err)
        return true // false = fatal
    })
```

`Run` handles `{"event":"ping"}` keepalive, bounded-backoff reconnect,
re-login and resubscribe of the registered set. Private pushes carry no
initial snapshot — pull the matching REST read after every (re)connect.

```go
// Spot incremental order book (Method 2, type=depth) — gap detection
tracker := &ws.DepthSequenceTracker{}
// inside onEvent:
//   if m.Type == "depth" { …tracker.Track(&startSeq, &endSeq) }
//   DepthGap → re-fetch snapshot via depth_snapshot, then tracker.Reset(lastSeq)
```

## Safety model

- **Mutations are never retried automatically** — after a transport
  timeout verify state via the matching read endpoint before resubmitting.
- Withdrawals/trades are never invoked implicitly.
- Monetary values are decimal **strings** (`Decimal`, `decimalArg`
  validates at the boundary; floats are rejected).
- Secrets/signatures/tokens are redacted from error strings (`errors.Redact`).
- Documented invariants are enforced client-side (exactly one of
  `closeNum`/`closeRate`, `closePiece` required when `stopType=1`,
  affiliate `partnerUids ≤ 100`, pageSize caps, K-line period table).

## Ambiguity gates

These operations return `errors.AmbiguousDocumentationError` — the official
docs contradict themselves:

| Endpoint | Conflict |
|---|---|
| `common.subAccount.getMaxTransferableAmount` | path `/api/v1/private` vs hyperlinked `/api/v1/api/v1/private` |
| `futures.trade.modifyStopLossTakeProfitInBatches` | `/v1/perpum/updateTpsl` vs snippet `/v1/perpum/TPSL` |
| `spot.check.getBatchHistoricalOrders` | `/v1/private` (no `/api` prefix) vs every other spot route |

Preserved oddities: affiliate base literally contains `patener`;
`futures.market.getHistoricalPublicTrades` uses the Futures HMAC profile
because its doc page requires authentication.

## Rate limits (documented)

- Spot global: 100 req/s per IP **and** 300 req/s per UID.
- Futures market data (`/v1/perpumPublic/*`): 30 req/s per IP.
- Futures other endpoints: 100 req/s per UID.
- Affiliate: 5 req/s per IP and UID; per-interface caps also enforced.
- Server violations → code `29001` → `errors.RateLimitError`.

Default limiter is in-process — inject a distributed `ratelimit.Limiter`
via `WithRateLimiter` for multi-worker deployments; `DisableRateLimiter()`
turns it off.

## Errors

Typed errors under `errors/`: `ValidationError`, `AuthenticationError`,
`RateLimitError`, `HTTPStatusError`, `APIError`, `ProtocolError`,
`TransportError`, `AmbiguousDocumentationError`, `CancelledError` —
classify with `errors.As`.

## Development

```bash
make test    # go test ./...
make vet     # go vet ./...
make fmt     # gofmt -w .
```

No live credentials needed — tests use scripted fake transports/dialers.

## Disclaimer

Unofficial community SDK. Not affiliated with CoinW. Recheck the changelog
(<https://www.coinw.com/api-doc/en/changelog/change_log>) before releases.
