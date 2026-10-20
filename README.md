# CoinW Go SDK

Use CoinW from Go without writing your own request signing, response decoding,
rate limiting, or WebSocket connection loop. The SDK brings Spot, Futures, and
account APIs into one client, with typed requests and responses and decimal
strings for financial values.

You can start with public market data and add private account access when you
need it. Trading, transfers, and withdrawals happen only when your code calls
the corresponding method. Creating a client or reconnecting a stream does not
place an order or move funds.

The implementation follows the [English CoinW documentation](https://www.coinw.com/api-doc/en)
snapshot dated **2026-10-20**. Its catalogue contains **86 REST operations** and
**26 WebSocket subscription variants** across five session profiles. Three
operations have conflicting routes in the documentation and remain blocked;
see [known documentation conflicts](#known-documentation-conflicts).

This is an unofficial community SDK, released under the [MIT license](LICENSE).
It is not affiliated with CoinW.

## Find your starting point

- [Install the module](#installation)
- [Make your first public request](#your-first-public-request)
- [Connect a private account](#private-account-access)
- [Choose a service](#where-to-find-an-operation)
- [Work with typed responses and money](#responses-and-decimal-values)
- [Use filters and instrument types](#requests-filters-and-instruments)
- [Submit an order explicitly](#submitting-an-order)
- [Receive WebSocket data](#websocket-streams)
- [Handle failures and rate limits](#errors-retries-and-rate-limits)
- [Run tests or contribute](#testing-and-contributing)

The [API index](docs/API.md) lists every mapped method and subscription, with
links to its CoinW source page and example fixture. If you are upgrading from
the earlier raw-response API, start with the [migration guide](docs/MIGRATION.md).

## Installation

```sh
go get github.com/tigusigalpa/coinw-go
```

Use **Go 1.26 or 1.27 with the latest patch**. The module declares Go 1.26 as
its minimum, and the main CI workflow tests both series. The choice follows
the maintained Go release series at the documentation snapshot; the reasoning
and update policy are in [CONTRIBUTING.md](CONTRIBUTING.md).

The only runtime dependency is
[`gorilla/websocket v1.5.3`](https://github.com/gorilla/websocket/releases/tag/v1.5.3),
licensed under BSD-2-Clause. It supplies the WebSocket transport. The SDK adds
separate adapters for CoinW's raw JSON and Socket.IO protocols and manages
serialized reads and writes around the connection.

## Your first public request

Public market data does not require an API key. This program fetches recent
Spot trades for `BTC_USDT` and prints their IDs, prices, and quantities:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    coinw "github.com/tigusigalpa/coinw-go"
)

func main() {
    client, err := coinw.NewClient()
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    response, err := client.SpotMarket().GetRecentTrades(ctx, "BTC_USDT", 0)
    if err != nil {
        log.Fatal(err)
    }

    for _, trade := range response.Data {
        fmt.Printf("trade=%s price=%s amount=%s\n",
            trade.ID.String(), trade.Price.String(), trade.Amount.String())
    }
}
```

The final `0` omits the optional `end` filter. You receive a slice of typed
trade records in `response.Data`; there is no JSON map to unpack first.

Reuse a client across requests. Give each operation an appropriate context,
usually the context of your HTTP handler, worker job, or application lifecycle.
The context also bounds time spent waiting for a rate-limit slot or a retry.

## Private account access

Spot and Futures have separate authentication profiles. Configure the profile
used by your application, or configure both if you need both markets. Public
requests remain credential-free on a client that also has private access.

This example reads a Spot account balance. It prints only the number of
returned currencies, leaving the balances themselves available to your code:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    coinw "github.com/tigusigalpa/coinw-go"
)

func main() {
    client, err := coinw.NewClient(coinw.WithSpotCredentials(
        os.Getenv("COINW_SPOT_API_KEY"),
        os.Getenv("COINW_SPOT_SECRET_KEY"),
    ))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    response, err := client.SpotAccount().GetSpotAccountBalance(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Loaded balances for %d currencies\n", len(response.Data))
}
```

For Futures, use `WithFuturesCredentials` with your Futures API key and secret.
The SDK returns an `AuthenticationError` if a private operation has no
credentials for its required profile.

| API family | Credentials used |
|---|---|
| Public Spot and Futures market data | None, except market operations whose source page explicitly requires authentication |
| Spot private REST, Common accounts and sub-accounts | Spot key and secret |
| Futures private REST and Affiliate | Futures key and secret |
| Spot private WebSocket | Spot key and secret in the documented login frame |
| Futures private WebSocket | Futures key and secret in its separate login frame |
| Spot public Socket.IO | Public token fetched by the SDK; no private API key |

Keep credentials in your environment or secret store. SDK error strings redact
known credential fields, but private response objects and raw envelopes are
still private data. Avoid dumping them into application logs. The
[security policy](SECURITY.md) explains what is retained in responses and errors.

### What the SDK signs for you

Spot private REST uses the documented uppercase MD5 signature over sorted
parameters and the secret. Futures private REST uses HMAC-SHA256/Base64 over
the timestamp, HTTP method, endpoint, and query or serialized body. Its
authentication is carried in the `sign`, `api_key`, and `timestamp` headers.

These profiles are independent. Account and sub-account calls use the Spot
profile; Affiliate calls use the Futures profile. You do not need to assemble
signatures yourself when using the service methods.

## Where to find an operation

The client groups methods by the same areas you will find in CoinW's docs:

| What you need | Client services |
|---|---|
| Spot instruments, tickers, candles, trades, order books | `SpotMarket()` |
| Spot balances | `SpotAccount()` |
| Spot order status and trading history | `SpotCheck()` |
| Place or cancel Spot orders | `SpotTrade()` |
| Futures instruments, prices, funding data, order books | `FuturesMarket()` |
| Futures balances, fees, coupons, margin settings | `FuturesAssets()` |
| Futures leverage information | `FuturesCommon()` |
| Futures order and SL/TP queries | `FuturesCheck()` |
| Current and historical Futures positions | `FuturesPosition()` |
| Futures orders, closes, margin changes, SL/TP | `FuturesTrade()` |
| Deposit and withdrawal information, withdrawals, asset transfers | `Account()` |
| Sub-account keys, balances available for transfer, transfers | `SubAccounts()` |
| Invited users, affiliate activity and team data | `Affiliate()` |

Method presence includes the registered operations whose routes are blocked
pending clarification. Check the [API index](docs/API.md) for their status and
the exact request and response types.

## Responses and decimal values

A successful REST call returns `*transport.Response[T]`. The type parameter
describes that operation's payload: a record, slice, map, page, or explicit
compatibility union.

| Field | What it contains |
|---|---|
| `response.Data` | The validated, operation-specific payload |
| `response.EndpointID` | The registered operation ID |
| `response.HTTPStatus` | HTTP response status |
| `response.Result.Data` | The original JSON `data` value |
| `response.Envelope` | The original response body |

Check the error before reading the response. The SDK validates HTTP status,
the business result, and the typed payload before returning a successful
response. Missing required structure, unknown typed fields, malformed tuples,
and invalid decimal tokens become `ProtocolError` values. A changed upstream
schema needs a contract update; it is not silently decoded into an arbitrary map.

### Money stays in decimal form

Prices, quantities, fees, margins, PnL, and funding rates use validated decimal
strings. That keeps values such as `0.000000000000000001` intact through JSON
serialization and avoids rounding through `float64`.

Use `coinw.NewDecimal` for values coming from your application:

```go
func encodeAmount(text string) ([]byte, error) {
    amount, err := coinw.NewDecimal(text)
    if err != nil {
        return nil, err
    }
    return json.Marshal(amount)
}
```

This helper uses `encoding/json` and the root `coinw` import. Passing
`"0.1000"` produces the JSON string `"0.1000"`. Scientific notation such as
`"1e-4"` is not a valid request `Decimal`; provide `"0.0001"` instead.
`MustDecimal` is useful for constants and tests, but panics on invalid input.

CoinW's response examples sometimes use JSON numbers or exponent notation for
the same financial field. `model.WireDecimal` handles that documented conflict:

| Wire value | `WireDecimal.String()` |
|---|---|
| `"0.1000"` | `0.1000` |
| `2.964e-05` | `0.00002964` |
| `"1.044E+5"` | `104400` |

`WireDecimal` serializes as a plain decimal string. Its `RawJSON()` method
returns a copy of the original token when you need to inspect its wire form.
The SDK does not provide arithmetic or choose rounding rules for you; keep
calculations in decimal-aware application code.

### When the documentation gives two shapes

Some differences are more than formatting. An ID may be a string or an integer;
a table may describe a scalar while its example wraps that value in an object.
The model makes those alternatives visible:

| Type | Documented difference it preserves |
|---|---|
| `StringOrInteger` | String or integer IDs, timestamps, and conflicting status fields |
| `ValueOrScalar[T]` | A scalar or a `{ "value": ..., "ts": ... }` wrapper |
| `ObjectOrArray[T]` | An object or an array |
| `JSONEncoded[T]` | JSON carried directly or inside a JSON string |
| `MaskedDecimal` | An affiliate balance that may be hidden as `"-"` or `""` |
| `LeverageOptions` | String or integer-array leverage options |

These adapters preserve their raw form for explicit inspection. They do not
infer timestamp units, convert contract quantities, or turn a hidden balance
into zero. The [contract notes](docs/CONTRACTS.md) explain the individual cases,
including the positional Futures REST candle tuple.

## Requests, filters, and instruments

Use a short service method when its arguments cover your request. When you
need additional documented filters, use its `WithRequest` method and request
record from `model`. Every REST operation outside the Futures trade family has
that form; Futures trade methods use their specific request records directly.

Optional pointer fields distinguish an omitted filter from an explicit zero.
For example, a pointer to `PositionModel(0)` means isolated margin, while a
nil optional field leaves the filter out. In the shorter historical-position
and transaction-detail methods, a negative `positionModel` omits that filter.

### Choose the right instrument spelling

CoinW uses different identifiers across its APIs. Keep the form expected by
the method you are calling:

| Context | Example |
|---|---|
| Spot REST trading pair | `BTC_USDT` |
| Spot public Socket.IO symbol | `BTC-USDT` |
| Spot raw public subscription | A pair code such as `"78"` in the documented examples |
| Futures USDT instrument | Base identifier such as `BTC` |
| Futures USDC instrument | Full pair such as `BTC_USDC` |

Get the current Spot pair-code mapping from
`SpotMarket().Get24hTradeSummaryForAllInstruments`; the examples' numeric code
is not a substitute for looking up the pair you want.

For typed Futures requests, construct `model.USDTInstrument("BTC")` or
`model.USDCInstrument("BTC_USDC")`. `model.InstrumentList` carries multiple
typed instruments for the documented comma-separated batch fields.

Using the published [K-Line request parameters](https://www.coinw.com/api-doc/en/futures-trading/market/get-k-line-of-an-instrument),
this helper requests 100 one-minute Futures candles in UTC. It uses the
`model` and `transport` packages alongside the root `coinw` import:

```go
func loadCandles(
    ctx context.Context,
    client *coinw.Client,
    instrument model.Instrument,
) (*transport.Response[[]model.FuturesCandleTuple], error) {
    limit := int64(100)
    timezone := "0"

    return client.FuturesMarket().GetKLineOfAnInstrumentWithRequest(ctx,
        model.FuturesMarketGetKLineOfAnInstrumentRequest{
            CurrencyCode: instrument,
            Granuality:   "0",
            Limit:        &limit,
            KlineType:    &timezone,
        })
}
```

`Granuality` is spelled that way in CoinW's request table. The SDK sends the
published field name. The REST candle response keeps five decimal slots after
its timestamp because the table and example disagree about high/open ordering;
read the contract note before assigning OHLC meaning to those slots.

## Submitting an order

Account-changing methods are explicit calls: placing or modifying orders,
cancelling them, closing or reversing positions, adjusting margin, configuring
SL/TP, transferring assets, resetting keys, and withdrawing funds. None run
as part of client construction or stream recovery.

The following helper **submits a Futures limit order when you call it**. Its
price, quantity, instrument, and client reference come from the caller; it has
no automatic invocation. It uses `context`, `coinw`, and `model` imports:

```go
func placeLimitOrder(
    ctx context.Context,
    client *coinw.Client,
    instrument model.Instrument,
    priceText, quantityText, clientOrderID string,
) (string, error) {
    price, err := coinw.NewDecimal(priceText)
    if err != nil {
        return "", err
    }
    quantity, err := coinw.NewDecimal(quantityText)
    if err != nil {
        return "", err
    }

    response, err := client.FuturesTrade().PlaceAnOrder(ctx,
        model.FuturesOrderRequest{
            Instrument:    instrument,
            Direction:     model.Long,
            Leverage:      1,
            QuantityUnit:  model.ContractQuantity,
            Quantity:      quantity,
            PositionModel: model.IsolatedMargin,
            PositionType:  model.LimitOrder,
            OpenPrice:     &price,
            ThirdOrderID:  clientOrderID,
        })
    if err != nil {
        return "", err
    }
    return response.Data.Value.String(), nil
}
```

Here the quantity is in **contracts**, as selected by `ContractQuantity`.
Choose `QuoteQuantity` or `BaseQuantity` only when that is the unit you intend;
the SDK does not convert between them. The returned ID acknowledges an order,
not a fill. For batch requests, inspect each item's `MsgCode` as well as the
outer result: individual items can fail.

A limit order requires `OpenPrice`. Other documented checks include
`closeNum` versus `closeRate` for closing a position, and `ClosePiece` for a
partial batch SL/TP request. An invalid request fails before network I/O.

Futures can hold long and short positions at the same time. An order in the
opposite direction can create a hedge; use the documented close methods when
you intend to close a position. Funding windows may also prevent trading.
After a mutation times out, check state through the corresponding read API
before deciding whether to submit it again.

## WebSocket streams

Choose the session factory for the protocol and account profile you need:

| Factory | Protocol and access | Example subscription constructor |
|---|---|---|
| `SpotPublicWS()` | Raw JSON, public Spot Method 2 | `ws.SpotTicker(pairCode)` |
| `SpotSocketIOWS()` | Socket.IO, public Spot Method 1 | `ws.SpotTickerSocketIO("BTC-USDT")` |
| `SpotPrivateWS()` | Raw JSON, private Spot | `ws.SpotAssets()` |
| `FuturesPublicWS()` | Raw JSON, public Futures on `/perpum` | `ws.FuturesTicker(instrument)` |
| `FuturesPrivateWS()` | Raw JSON, private Futures on `/perpum` | `ws.FuturesAssets()` |

Import the stream package as
`ws "github.com/tigusigalpa/coinw-go/websocket"`. A typed subscription binds
the requested channel to its payload type, so you decode using the same
subscription value that you passed to subscribe.

### A public Futures feed

This complete example watches the BTC USDT ticker until you press Ctrl+C.
Transport failures may reconnect; other SDK errors stop the loop so the
application can decide what to do with them.

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log"
    "os"
    "os/signal"

    coinw "github.com/tigusigalpa/coinw-go"
    coinwerr "github.com/tigusigalpa/coinw-go/errors"
    "github.com/tigusigalpa/coinw-go/model"
    ws "github.com/tigusigalpa/coinw-go/websocket"
)

func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    if err := watchTicker(ctx); err != nil && !errors.Is(err, context.Canceled) {
        log.Fatal(err)
    }
}

func watchTicker(ctx context.Context) error {
    client, err := coinw.NewClient()
    if err != nil {
        return err
    }
    instrument, err := model.USDTInstrument("BTC")
    if err != nil {
        return err
    }

    subscription := ws.FuturesTicker(instrument)
    session := client.FuturesPublicWS()
    if err := session.Connect(ctx); err != nil {
        return err
    }
    defer session.Close()
    if err := ws.SubscribeStream(ctx, session, subscription); err != nil {
        return err
    }

    var decodeErr error
    runErr := session.Run(ctx,
        func(message *ws.Message) bool {
            if !message.IsEvent() {
                return true
            }
            ticker, err := ws.DecodeStream(message, subscription)
            if err != nil {
                decodeErr = err
                return false
            }
            fmt.Println(ticker.Last.String())
            return true
        },
        func(err error) bool {
            log.Printf("ticker stream: %v", err)
            var transportErr *coinwerr.TransportError
            return errors.As(err, &transportErr)
        })
    if decodeErr != nil {
        return decodeErr
    }
    return runErr
}
```

The event callback returns `true` to keep reading and `false` to stop. The
error callback returns `true` to allow recovery or `false` to return the error.
`Run` handles the documented heartbeat, bounded reconnect backoff, login where
needed, and restoration of the registered subscriptions. Server rejections,
including `29001`, are delivered as typed errors.

Callbacks run synchronously, which keeps delivery from building an unbounded
queue. Keep them short and context-aware; heartbeats run independently. A
session supports one active `Run`. Cancelling the connection's context closes
it, and a read timeout requires a new connection.

### Spot Socket.IO has its own adapter

For Method 1, use `SpotSocketIOWS`, `SubscribeSocketIO`, and `DecodeSocketIO`.
The SDK fetches the public token and completes the Engine.IO/Socket.IO
handshake before subscribing. The default transport URL uses
`/socket.io/?EIO=4&transport=websocket&token=...`.

Several Spot payloads carry JSON inside a string. Their typed decoder returns
`model.JSONEncoded[T]`; read the decoded payload through `.Value`, for example
`payload.Value.Last.String()` for a Spot ticker. `RawJSON()` retains the
original form. Passing a Socket.IO subscription to the raw adapter, or the
reverse, returns a validation error.

### Private streams need a starting snapshot

A private push does not promise an initial account snapshot. Supply
`WithOnConnect` to fetch the matching REST state after login and before
subscriptions are restored. Typed private subscriptions require this hook.

For a Spot asset stream, the setup can look like this. The caller provides
`applySnapshot` to store the baseline in application state:

```go
func newSpotAssetsSession(
    client *coinw.Client,
    applySnapshot func(map[string]model.WireDecimal),
) (*ws.RawJsonSession, error) {
    return client.SpotPrivateWS(ws.WithOnConnect(func(ctx context.Context) error {
        response, err := client.SpotAccount().GetSpotAccountBalance(ctx)
        if err != nil {
            return err
        }
        applySnapshot(response.Data)
        return nil
    }))
}
```

Then connect the session and subscribe with `ws.SpotAssets()`. The hook runs
again after each successful login; if it fails, that connection is closed.
Choose a snapshot whose scope matches the subscription, especially when its
REST method requires instruments or other filters. Your application remains
responsible for reconciling events with that baseline; timestamps alone do
not establish delivery order or an atomic snapshot boundary.

### Incremental Spot depth needs a sequence anchor

Use `ws.DepthSequenceTracker` with the documented `depth_snapshot` stream:

1. Receive a snapshot and call `Reset(snapshotSeq)` to anchor the tracker.
2. Check incremental updates with `Track(startSeq, endSeq)`.
3. Apply `DepthApplied` updates and discard `DepthDuplicate` updates.
4. On `DepthNeedsSnapshot` or `DepthGap`, request a fresh snapshot and wait for
   its anchor before applying further increments.

`TrackAndResync` calls your snapshot-request callback for an absent anchor or
gap. Requesting a snapshot does not itself reset the tracker: call `Reset`
only when the snapshot arrives. A gap never advances the accepted sequence.
See [snapshot and ordering policy](docs/CONTRACTS.md#snapshot-and-ordering-policy)
for the contract details.

## Errors, retries, and rate limits

Use `errors.Is` for context cancellation or deadlines and `errors.As` for SDK
error classes. If you wrap an error in your own code, use `%w` to keep that
classification available.

The following helper uses standard `context`, `errors`, and `log` imports,
plus `coinwerr "github.com/tigusigalpa/coinw-go/errors"`:

```go
func reportRequestError(err error) {
    var rateLimit *coinwerr.RateLimitError
    var apiError *coinwerr.APIError

    switch {
    case errors.Is(err, context.Canceled):
        log.Print("CoinW operation cancelled")
    case errors.Is(err, context.DeadlineExceeded):
        log.Print("CoinW deadline reached; check state before resubmitting a mutation")
    case errors.As(err, &rateLimit):
        log.Printf("CoinW rate limit: %v", err)
    case errors.As(err, &apiError):
        log.Printf("CoinW operation %s failed with code %s: %v",
            apiError.EndpointID, apiError.Code, err)
    default:
        log.Printf("CoinW request failed: %v", err)
    }
}
```

Other classes include `ValidationError`, `AuthenticationError`,
`HTTPStatusError`, `ProtocolError`, `SerializationError`, `TransportError`,
`AmbiguousDocumentationError`, and `CancelledError`. Error strings are
redacted; avoid logging arbitrary dumps of wrapped causes or private payloads.

### Retry behavior

Only eligible read-only **GET** requests may retry after a transport or
rate-limit failure. The default allows two retries, with bounded exponential
backoff and jitter. Set `WithMaxReadRetries(0)` to disable automatic retries.

**POST, PUT, and DELETE are never retried automatically**, including reads
that CoinW implements as POST. A timeout on a mutation does not tell you
whether the exchange accepted it. Read the relevant state before making a
resubmission decision; do not wrap order or withdrawal calls in a generic retry loop.

### Limits apply together

| Scope | Documented global budget |
|---|---|
| Spot | 100 requests/s per IP and 300 requests/s per UID |
| Futures Group A: `klines`, `tickers`, `ticker`, `trades`, `depth` | 30 requests/s per IP |
| Other Futures operations | 100 requests/s per UID |
| Affiliate | 5 requests/s per IP and UID, with its documented interface limits |

Per-interface limits also apply, including fractional rates such as one
request per five seconds. IP and UID scopes are tracked separately. HTTP 429
and business code `29001` produce `RateLimitError`.

The default limiter lives in one client instance. Without a configured UID,
calls share a conservative local UID budget. To coordinate several clients
in one process, inject the same limiter into them. For several application
instances, provide a shared implementation of `ratelimit.Limiter` through
`WithRateLimiter`. Creating more clients does not increase CoinW's allowance.

`DisableRateLimiter()` is available for controlled test setups or applications
that manage these budgets elsewhere.

## Configuring the client

The defaults are the documented REST base `https://api.coinw.com`, a
15-second HTTP attempt timeout, two eligible read retries, and the in-memory
limiter. An outer context deadline can bound the complete operation, including
limiter waits and retry delays.

You can supply an HTTP client and adjust those settings with functional
options. This helper uses `net/http`, `time`, and the root `coinw` import:

```go
func newHTTPClient() (*coinw.Client, error) {
    return coinw.NewClient(
        coinw.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
        coinw.WithRequestTimeout(5*time.Second),
        coinw.WithMaxReadRetries(1),
    )
}
```

The SDK copies configuration and HTTP client settings during construction.
Injected transports, limiters, clocks, and callbacks should be safe for
concurrent use. Pass `WithUID` when you have the UID for the limit scope.

For tests, proxies, or application-specific infrastructure, the options also
include `WithRestBaseURL`, `WithWebSocketURLs`, `WithRateLimiter`, `WithClock`,
`WithDialer`, and `WithRetryPolicy`. The retry hook changes timing, while the
restriction to eligible read-only GET requests still applies.

## Known documentation conflicts

Some CoinW pages disagree about the route to call. These operations have typed
models and contract fixtures, but return `AmbiguousDocumentationError` without
sending a request until their route is verified:

| Operation | Published conflict |
|---|---|
| `common.subAccount.getMaxTransferableAmount` | `/api/v1/private` versus a link containing `/api/v1/api/v1/private` |
| `spot.check.getBatchHistoricalOrders` | Conflicting host/path, including `/v1/private` without the usual `/api` prefix |
| `futures.trade.modifyStopLossTakeProfitInBatches` | Endpoint `/v1/perpum/updateTpsl` versus example `/v1/perpum/TPSL` |

Spot `CancelAllOrders` also requires an explicit instrument. Both the short
method and its typed request form block an empty filter because the docs do
not consistently define cancellation across all markets.

A few other details can look surprising when you first encounter them:

- Affiliate routes retain the literal `patener` spelling from CoinW's paths.
- `GetHistoricalPublicTrades` and `GetMarginRequirementsOfAllInstruments`
  use Futures private authentication because their pages require it.
- The Affiliate trading-volume list sends a JSON params object even when
  there are no filters, so the empty body is `{}`.
- Spot K-Line uses the published endpoint and interval table despite the
  conflicting hyperlink example.
- Conflicting money formats and response shapes use the named adapters
  described above; the SDK does not infer their missing business meaning.

The full list is in [CONTRACTS.md](docs/CONTRACTS.md). Check the
[CoinW changelog](https://www.coinw.com/api-doc/en/changelog/change_log) when
reviewing an upstream change or preparing a release.

## Testing and contributing

The normal test suite uses local HTTP servers, fake transports, and scripted
WebSocket peers. It needs no live credentials:

```sh
go test ./...
go test -race -timeout 2m ./...
go vet ./...
```

It exercises every catalogue row, strict decoding, fixed signature vectors,
authentication separation, mutation retry rules, cancellation, reconnect,
snapshot hooks, and sequence-gap handling. The
[coverage catalogue](contracts/catalogue.json) records each row's source and
fixture. Passing an offline fixture establishes that documented contract;
live exchange behavior needs its own integration check.

The main [CI workflow](.github/workflows/ci.yml) runs the Go 1.26.x and 1.27.x
matrix, race checks, builds, formatting, vet, pinned Staticcheck and
golangci-lint, module verification, and vulnerability checks. The
[local verification record](docs/VERIFICATION.md) states which checks were
actually run and which live checks were left disabled.

### Optional live integration tests

Without opt-in, all `TestIntegration*` tests skip. To enable the public REST
check, explicitly set the flag and the API base you intend to use:

```sh
export COINW_INTEGRATION=1
export COINW_TEST_BASE_URL=https://api.coinw.com
go test -run '^TestIntegrationPublicREST$' -v .
```

The equivalent PowerShell setup is:

```powershell
$env:COINW_INTEGRATION = '1'
$env:COINW_TEST_BASE_URL = 'https://api.coinw.com'
go test -run '^TestIntegrationPublicREST$' -v .
```

Those commands make real requests to the specified API. Private read checks
need separate `COINW_TEST_SPOT_*` or `COINW_TEST_FUTURES_*` credentials;
WebSocket checks need explicit URLs for their selected adapters. The one-order
cancellation check additionally needs `COINW_TEST_MUTATIONS=1` and
`COINW_TEST_CANCEL_ORDER_ID`, and always skips when `CI` is set.

Follow [INTEGRATION.md](docs/INTEGRATION.md) for the complete variable names and
test selection. Keep live credentials and account responses out of fixtures
and commits, and remove test environment variables when you are finished.

### Project guide

| Document | What to use it for |
|---|---|
| [API index](docs/API.md) | Find a method, subscription, source page, or fixture |
| [Contract notes](docs/CONTRACTS.md) | Understand response conflicts, ordering, and snapshot policy |
| [Migration guide](docs/MIGRATION.md) | Update earlier code to typed responses and requests |
| [Integration guide](docs/INTEGRATION.md) | Run selected live checks explicitly |
| [Verification record](docs/VERIFICATION.md) | Review the recorded local check results |
| [Contributing](CONTRIBUTING.md) | Set up development tools and change a contract |
| [Release checklist](docs/RELEASING.md) | Review versioning, upstream changes, and publication steps |
| [Changelog](CHANGELOG.md) | See what changed in the SDK |
| [Security policy](SECURITY.md) | Report a vulnerability and handle sensitive data |
| [Code of conduct](CODE_OF_CONDUCT.md) | Participate in the project |

### Package layout

Most application code starts with `coinw.Client` and uses `model` for request
records, enums, and values. The remaining packages keep the implementation
parts available for testing and more specialized integrations:

| Package | Responsibility |
|---|---|
| Root `coinw` | Client construction, options, `Decimal`, service and session factories |
| `services` | REST methods grouped by domain |
| `model` | Request and response records, instrument types, decimal and compatibility values |
| `transport` | HTTP execution, signing dispatch, response validation, eligible retries |
| `auth` | Separate Spot and Futures signing helpers |
| `endpoints` | Registered routes, methods, authentication and limit metadata |
| `ratelimit` | Context-aware in-process budgets and the injectable limiter interface |
| `websocket` | Raw JSON and Socket.IO sessions, typed subscriptions, recovery and sequence tracking |
| `errors` | Typed SDK errors and credential redaction |
