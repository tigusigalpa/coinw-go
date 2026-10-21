# CoinW Golang Client/SDK/Library

![CoinW go/golang sdk library package client](https://i.postimg.cc/7Lc7tVMf/coinw-go-sdk-hero-banner.jpg)

A Go client for the [CoinW](https://www.coinw.com) exchange API. It takes care
of the tedious parts — request signing, response decoding, rate limits, and
keeping WebSocket connections alive — so you can focus on your own logic.

What you get:

- Spot, Futures, and account APIs behind a single client
- Typed requests and responses instead of `map[string]any`
- Prices and amounts as decimal strings, never `float64`
- A built-in rate limiter that follows CoinW's published budgets
- WebSocket sessions that reconnect, log in again, and resubscribe on their own

Nothing touches your money unless you ask for it. Creating a client or
reconnecting a stream never places an order or moves funds — only an explicit
method call does.

> **Unofficial.** This is a community project under the [MIT license](LICENSE).
> It is not affiliated with or endorsed by CoinW.

**API coverage.** The SDK is built against the
[English CoinW docs](https://www.coinw.com/api-doc/en) as of **2026-10-20**:
**86 REST operations** and **26 WebSocket subscription variants** across five
session profiles. Three operations are blocked because the docs give
conflicting routes for them — see
[known documentation conflicts](#known-documentation-conflicts).

Every method and subscription is listed in the [API index](docs/API.md), with a
link to its CoinW page and its test fixture. Upgrading from the old
raw-response API? Start with the [migration guide](docs/MIGRATION.md).

## Contents

- [Installation](#installation)
- [Quick start](#quick-start)
- [Private account access](#private-account-access)
- [Finding the right service](#finding-the-right-service)
- [Responses and money](#responses-and-money)
- [Requests, filters, and instruments](#requests-filters-and-instruments)
- [Placing orders](#placing-orders)
- [WebSocket streams](#websocket-streams)
- [Errors, retries, and rate limits](#errors-retries-and-rate-limits)
- [Configuration](#configuration)
- [Known documentation conflicts](#known-documentation-conflicts)
- [Testing and contributing](#testing-and-contributing)

## Installation

```sh
go get github.com/tigusigalpa/coinw-go
```

You'll need **Go 1.26 or 1.27** with the latest patch release. Go 1.26 is the
minimum, and CI tests both. Why these two, and when that changes, is explained
in [CONTRIBUTING.md](CONTRIBUTING.md).

There is exactly one dependency:
[`gorilla/websocket v1.5.3`](https://github.com/gorilla/websocket/releases/tag/v1.5.3)
(BSD-2-Clause). It handles the socket itself. On top of it, the SDK adds its own
adapters for CoinW's raw JSON and Socket.IO protocols and makes sure reads and
writes on a connection never overlap.

## Quick start

Public market data doesn't need an API key. This program prints recent Spot
trades for `BTC_USDT`:

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

The trailing `0` means "no `end` filter". `response.Data` is already a slice of
typed trades, so there's no JSON to dig through.

Two habits worth picking up early:

- **Create one client and reuse it.**
- **Pass a real context into every call** — your HTTP handler's, your worker's,
  or your application's. It also caps how long a call may wait for a rate-limit
  slot or a retry.

## Private account access

Spot and Futures use separate API keys. Configure the one you need, or both.
Public calls keep working without credentials, even on a client that has them.

Reading a Spot balance looks like this. The example only prints how many
currencies came back; the balances themselves are in `response.Data`:

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

For Futures, use `WithFuturesCredentials`. If you call a private method without
the credentials it needs, you get an `AuthenticationError`.

Which key goes where:

| API family | Credentials |
|---|---|
| Public Spot and Futures market data | None — except the few market endpoints whose docs page requires auth |
| Spot private REST, Common accounts, sub-accounts | Spot key and secret |
| Futures private REST, Affiliate | Futures key and secret |
| Spot private WebSocket | Spot key and secret, sent in the login frame |
| Futures private WebSocket | Futures key and secret, sent in its own login frame |
| Spot public Socket.IO | A public token the SDK fetches for you; no API key |

Keep your keys in environment variables or a secret store. The SDK redacts
known credential fields from its error messages, but responses still contain
private account data — don't dump them into your logs. The
[security policy](SECURITY.md) covers exactly what ends up in responses and
errors.

### Signing

You never build signatures by hand; the service methods do it for you. If
you're curious how:

- **Spot private REST** — uppercase MD5 over the sorted parameters plus the
  secret.
- **Futures private REST** — HMAC-SHA256, Base64-encoded, over the timestamp,
  HTTP method, endpoint, and query string or body. Sent in the `sign`,
  `api_key`, and `timestamp` headers.

Account and sub-account calls sign with the Spot key. Affiliate calls sign with
the Futures key.

## Finding the right service

Methods are grouped the same way CoinW's docs are:

| What you need | Service |
|---|---|
| Spot instruments, tickers, candles, trades, order books | `SpotMarket()` |
| Spot balances | `SpotAccount()` |
| Spot order status and trade history | `SpotCheck()` |
| Placing and cancelling Spot orders | `SpotTrade()` |
| Futures instruments, prices, funding, order books | `FuturesMarket()` |
| Futures balances, fees, coupons, margin settings | `FuturesAssets()` |
| Futures leverage info | `FuturesCommon()` |
| Futures order and SL/TP queries | `FuturesCheck()` |
| Current and historical Futures positions | `FuturesPosition()` |
| Futures orders, closes, margin changes, SL/TP | `FuturesTrade()` |
| Deposits, withdrawals, asset transfers | `Account()` |
| Sub-account keys, transferable balances, transfers | `SubAccounts()` |
| Invited users, affiliate activity, team data | `Affiliate()` |

The blocked operations still have methods; they just return an error without
sending anything. The [API index](docs/API.md) shows each operation's status
and its exact request and response types.

## Responses and money

Every successful REST call returns `*transport.Response[T]`, where `T` is that
operation's payload: a record, a slice, a map, a page, or one of the union
types described [below](#when-the-docs-disagree-on-shape).

| Field | What's in it |
|---|---|
| `response.Data` | The validated, typed payload |
| `response.EndpointID` | The operation ID |
| `response.HTTPStatus` | The HTTP status code |
| `response.Result.Data` | The original JSON `data` value |
| `response.Envelope` | The original response body |

Always check `err` first. By the time you get a response, the SDK has already
checked the HTTP status, CoinW's business result code, and the payload shape.
Missing fields, unknown fields, malformed tuples, and invalid decimals all come
back as a `ProtocolError`. If CoinW changes a schema, you get an error — not a
half-decoded map that quietly looks fine. Fixing it means updating the
contract.

### Money stays decimal

Prices, quantities, fees, margins, PnL, and funding rates are validated decimal
strings. A value like `0.000000000000000001` survives a JSON round-trip intact,
because it never passes through `float64`.

To turn your own input into a decimal, use `coinw.NewDecimal`:

```go
func encodeAmount(text string) ([]byte, error) {
    amount, err := coinw.NewDecimal(text)
    if err != nil {
        return nil, err
    }
    return json.Marshal(amount)
}
```

`"0.1000"` encodes as the JSON string `"0.1000"`, trailing zeros included.
Scientific notation such as `"1e-4"` isn't accepted in requests — write
`"0.0001"` instead. `MustDecimal` is handy for constants and tests, but it
panics on invalid input.

CoinW itself is less strict. Its responses sometimes send the same money field
as a JSON number, or in exponent form. `model.WireDecimal` accepts all of these
and normalizes them:

| On the wire | `WireDecimal.String()` |
|---|---|
| `"0.1000"` | `0.1000` |
| `2.964e-05` | `0.00002964` |
| `"1.044E+5"` | `104400` |

It serializes back as a plain decimal string, and `RawJSON()` gives you a copy
of the original token if you need it. There's no arithmetic and no rounding
policy built in: do your math with whatever decimal library you trust.

### When the docs disagree on shape

Sometimes the difference is more than formatting. An ID is a string on one page
and an integer on another, or a table describes a plain value while the example
wraps it in an object. Instead of guessing, the model gives each of these cases
an explicit type:

| Type | What it covers |
|---|---|
| `StringOrInteger` | IDs, timestamps, and status fields that appear as either a string or an integer |
| `ValueOrScalar[T]` | A plain value or a `{ "value": ..., "ts": ... }` wrapper |
| `ObjectOrArray[T]` | An object or an array |
| `JSONEncoded[T]` | JSON sent directly or packed inside a string |
| `MaskedDecimal` | An affiliate balance that may be hidden as `"-"` or `""` |
| `LeverageOptions` | Leverage options as a string or an array of integers |

Each type keeps the raw value so you can inspect it. None of them try to be
clever: they don't guess timestamp units, convert contract quantities, or turn
a hidden balance into zero. The [contract notes](docs/CONTRACTS.md) go through
every case, including the positional Futures REST candle tuple.

## Requests, filters, and instruments

The short service methods take the common arguments directly. When you need
more filters, use the `...WithRequest` variant with a request struct from
`model`. Every REST operation has one, except the Futures trade family — those
methods take their request structs directly.

Optional fields are pointers, so "not set" and "zero" mean different things. A
pointer to `PositionModel(0)` means isolated margin; `nil` leaves the filter
out. One exception: in the short historical-position and transaction-detail
methods, pass a negative `positionModel` to skip that filter.

### Instrument names

CoinW spells instruments differently depending on the API. Use whatever the
method you're calling expects:

| Where | Example |
|---|---|
| Spot REST trading pair | `BTC_USDT` |
| Spot public Socket.IO symbol | `BTC-USDT` |
| Spot raw public subscription | A pair code, such as `"78"` in CoinW's examples |
| Futures USDT instrument | Just the base, such as `BTC` |
| Futures USDC instrument | The full pair, such as `BTC_USDC` |

Don't hard-code that `"78"` — it's only an example. Look up the real pair code
with `SpotMarket().Get24hTradeSummaryForAllInstruments`.

For typed Futures requests, build instruments with `model.USDTInstrument("BTC")`
or `model.USDCInstrument("BTC_USDC")`. For batch fields that take a
comma-separated list, use `model.InstrumentList`.

Here's a helper that fetches 100 one-minute Futures candles in UTC, using the
[K-Line request parameters](https://www.coinw.com/api-doc/en/futures-trading/market/get-k-line-of-an-instrument)
(imports: `coinw`, `model`, `transport`):

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

Yes, it's `Granuality`. That's how CoinW spells the field, and the SDK sends it
exactly as published.

Also note that each REST candle has five decimal slots after the timestamp, and
the SDK doesn't name them. CoinW's table and example disagree on whether high
or open comes first, so read the contract note before treating them as OHLC.

## Placing orders

Anything that changes your account is an explicit method call: placing,
modifying, or cancelling orders; closing or reversing positions; adjusting
margin; setting SL/TP; transferring assets; resetting keys; withdrawing funds.
None of it ever happens behind your back during client setup or stream
recovery.

This helper **places a real Futures limit order every time you call it**. The
caller supplies the instrument, price, quantity, and client order ID
(imports: `context`, `coinw`, `model`):

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

Things to keep in mind:

- **Quantity is in contracts** here, because of `ContractQuantity`. Switch to
  `QuoteQuantity` or `BaseQuantity` only if that's really the unit you mean;
  the SDK never converts between them.
- **The returned ID means "accepted", not "filled".**
- **Batch requests can partly fail.** Check each item's `MsgCode`, not just the
  outer result.
- **Some mistakes are caught locally.** A limit order without `OpenPrice`, a
  close that breaks the `closeNum`/`closeRate` rules, or a partial batch SL/TP
  without `ClosePiece` fails before anything is sent.
- **Futures allows hedging.** You can hold a long and a short at the same time,
  so an order in the opposite direction opens a hedge instead of closing your
  position. To close, use the close methods. Trading can also be unavailable
  during funding windows.
- **If a mutation times out, look before you resend.** Query the order or
  position first — the exchange may have accepted it.

## WebSocket streams

Pick the session factory for the protocol and account you need:

| Factory | Protocol and access | Example subscription |
|---|---|---|
| `SpotPublicWS()` | Raw JSON, public Spot (Method 2) | `ws.SpotTicker(pairCode)` |
| `SpotSocketIOWS()` | Socket.IO, public Spot (Method 1) | `ws.SpotTickerSocketIO("BTC-USDT")` |
| `SpotPrivateWS()` | Raw JSON, private Spot | `ws.SpotAssets()` |
| `FuturesPublicWS()` | Raw JSON, public Futures on `/perpum` | `ws.FuturesTicker(instrument)` |
| `FuturesPrivateWS()` | Raw JSON, private Futures on `/perpum` | `ws.FuturesAssets()` |

Import the package as `ws "github.com/tigusigalpa/coinw-go/websocket"`.
Subscriptions are typed: the value you subscribe with is the same value you
decode with, so the decoder always knows what payload to expect.

### Example: a public Futures ticker

This program prints the BTC USDT ticker until you press Ctrl+C. Transport
errors trigger a reconnect; any other SDK error stops the loop and comes back
to you.

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

How the two callbacks work:

- The **event callback** returns `true` to keep reading, `false` to stop.
- The **error callback** returns `true` to let the session recover, `false` to
  stop and return the error.

`Run` handles heartbeats, reconnects with bounded backoff, logs in when needed,
and restores your subscriptions. Server rejections, including `29001`, arrive
as typed errors. If the server refuses a single subscription, the connection
stays up: the error goes to your error callback, the other streams keep
running, and the refused stream is dropped so reconnects don't send it again.

The session remembers the subscriptions you asked for, even while it's
disconnected. `Subscribe` and `Unsubscribe` work offline, and every connect
restores the current set. Calling `Close` stops an active `Run`, which then
returns `nil`.

Callbacks run synchronously, so there's no hidden queue that can grow forever.
The flip side is that a slow callback slows down reading, so keep them short
and respect the context. Heartbeats run separately and won't be held up.

A few more rules: one session runs one `Run` at a time; cancelling the
connection's context closes it; after a read timeout you need a new connection
(`Run` takes care of that).

### Spot over Socket.IO

Spot Method 1 speaks Socket.IO, so it has its own set of functions:
`SpotSocketIOWS`, `SubscribeSocketIO`, and `DecodeSocketIO`. The SDK fetches the
public token and completes the Engine.IO/Socket.IO handshake before it
subscribes. The default URL is
`/socket.io/?EIO=4&transport=websocket&token=...`.

Several Spot payloads are JSON packed inside a string. Their decoders return
`model.JSONEncoded[T]`, so you reach the data through `.Value` — for a Spot
ticker, that's `payload.Value.Last.String()`. `RawJSON()` still gives you the
original. Mixing the two adapters (a Socket.IO subscription on a raw session,
or the other way round) returns a validation error.

### Private streams: start from a snapshot

CoinW doesn't promise that a private stream begins with your current account
state. That's what `WithOnConnect` is for: it loads the state over REST on
every connection, after login and after your subscriptions are sent, but before
the first event is read. Events pushed while the snapshot loads wait in the
socket, so nothing falls into the gap between the two. Typed private
subscriptions require this hook.

Here's a Spot assets session. `applySnapshot` is your own function that stores
the baseline:

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

Then subscribe with `ws.SpotAssets()` and call `Run`. The hook runs once per
connection, and if it fails, that connection is closed. Subscribe before you
start reading: a subscription added later on the same connection isn't covered
by its snapshot.

Make sure the snapshot covers the same scope as the subscription — this matters
most when the REST method needs instruments or other filters. Merging live
events into the snapshot is your job. Some early events may already be
reflected in it, and timestamps alone don't tell you the delivery order or
where exactly the snapshot ends.

### Incremental Spot depth

For the `depth_snapshot` stream, `ws.DepthSequenceTracker` keeps your order
book consistent:

1. When a snapshot arrives, call `Reset(snapshotSeq)`.
2. For each incremental update, call `Track(startSeq, endSeq)`.
3. Apply `DepthApplied` updates; drop `DepthDuplicate` ones.
4. On `DepthNeedsSnapshot` or `DepthGap`, request a new snapshot and don't apply
   anything until it arrives.

`TrackAndResync` does the same, and calls your snapshot-request callback when
there's no anchor yet or a gap shows up. Note that requesting a snapshot
doesn't reset the tracker — call `Reset` when the snapshot actually arrives. A
gap never moves the accepted sequence forward. The details are in the
[snapshot and ordering policy](docs/CONTRACTS.md#snapshot-and-ordering-policy).

## Errors, retries, and rate limits

Use `errors.Is` for context cancellation and deadlines, and `errors.As` for SDK
error types. If you wrap errors yourself, use `%w` so this keeps working.

For example (imports: `context`, `errors`, `log`, and
`coinwerr "github.com/tigusigalpa/coinw-go/errors"`):

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

The other error types are `ValidationError`, `AuthenticationError`,
`HTTPStatusError`, `ProtocolError`, `SerializationError`, `TransportError`,
`AmbiguousDocumentationError`, and `CancelledError`. Error messages are
redacted, but still avoid logging raw wrapped causes or private payloads
wholesale.

### Retries

Only read-only **GET** requests are retried automatically, after a transport
error or a rate-limit hit. By default that's up to two retries with exponential
backoff and jitter. `WithMaxReadRetries(0)` turns retries off.

**POST, PUT, and DELETE are never retried** — and that includes the read
endpoints CoinW happens to implement as POST. When a mutation times out, you
simply don't know whether the exchange took it. Check the state before you
resend, and never wrap order or withdrawal calls in a generic retry loop.

### Rate limits

CoinW enforces several budgets at once:

| Scope | Global budget |
|---|---|
| Spot | 100 requests/s per IP and 300 requests/s per UID |
| Futures Group A: `klines`, `tickers`, `ticker`, `trades`, `depth` | 30 requests/s per IP |
| Other Futures operations | 100 requests/s per UID |
| Affiliate | 5 requests/s per IP and UID, plus per-interface limits |

Per-interface limits apply on top of these, and some are slow — as low as one
request every five seconds. IP and UID budgets are tracked separately. The
built-in limiter uses a sliding window, so it never lets more than N requests
through in any one-second span, even right after a quiet period. HTTP 429
and business code `29001` both come back as a `RateLimitError`.

By default, each client has its own in-memory limiter. Without a configured
UID, all calls share one conservative local UID budget. To share limits:

- **between clients in one process**, give them the same limiter;
- **across processes or machines**, implement `ratelimit.Limiter` on top of
  shared storage and pass it in with `WithRateLimiter`.

Spinning up more clients won't get you a bigger allowance — CoinW counts per IP
and UID, not per client.

`DisableRateLimiter()` is there for tests, or for applications that enforce
these limits somewhere else.

## Configuration

Out of the box you get the REST base `https://api.coinw.com`, a 15-second
timeout per HTTP attempt, two read retries, and the in-memory limiter. A
deadline on your context caps the whole operation, limiter waits and retry
delays included.

Change any of it with functional options (imports: `net/http`, `time`,
`coinw`):

```go
func newHTTPClient() (*coinw.Client, error) {
    return coinw.NewClient(
        coinw.WithHTTPClient(&http.Client{Timeout: 10 * time.Second}),
        coinw.WithRequestTimeout(5*time.Second),
        coinw.WithMaxReadRetries(1),
    )
}
```

The client copies its configuration and HTTP client settings when it's
created, so changing them afterwards has no effect. Anything you inject —
transports, limiters, clocks, callbacks — must be safe for concurrent use. If
you know your UID, pass `WithUID` so the limiter uses the right scope.

For tests, proxies, or custom infrastructure there are also `WithRestBaseURL`,
`WithWebSocketURLs`, `WithRateLimiter`, `WithClock`, `WithDialer`, and
`WithRetryPolicy`. A custom retry policy only changes the timing — it can't
make anything other than read-only GET requests retryable.

## Known documentation conflicts

A few CoinW pages disagree about which route to call. These operations have
typed models and fixtures, but until their routes are confirmed they return
`AmbiguousDocumentationError` without sending anything:

| Operation | What the docs say |
|---|---|
| `common.subAccount.getMaxTransferableAmount` | `/api/v1/private` in one place, a link to `/api/v1/api/v1/private` in another |
| `spot.check.getBatchHistoricalOrders` | Conflicting host and path, including `/v1/private` without the usual `/api` prefix |
| `futures.trade.modifyStopLossTakeProfitInBatches` | The endpoint says `/v1/perpum/updateTpsl`, the example says `/v1/perpum/TPSL` |

Spot `CancelAllOrders` also requires an explicit instrument. Both the short
method and the request form reject an empty filter, because the docs don't
clearly say what "cancel all" means across markets.

Other quirks you might trip over:

- Affiliate routes really do say `patener`. That's CoinW's spelling.
- `GetHistoricalPublicTrades` and `GetMarginRequirementsOfAllInstruments` use
  Futures private auth, because their docs pages require it.
- The Affiliate trading-volume list always sends a JSON params object, so with
  no filters the body is `{}`.
- Spot K-Line follows the published endpoint and interval table, not the
  conflicting hyperlink example.
- Conflicting money formats and response shapes go through the adapter types
  described above. The SDK doesn't invent business meaning the docs leave out.

The full list is in [CONTRACTS.md](docs/CONTRACTS.md). When something changes
upstream, start with CoinW's
[changelog](https://www.coinw.com/api-doc/en/changelog/change_log).

## Testing and contributing

The test suite runs entirely offline, against local HTTP servers, fake
transports, and scripted WebSocket peers. No credentials needed:

```sh
go test ./...
go test -race -timeout 2m ./...
go vet ./...
```

It covers every catalogue entry, strict decoding, fixed signature vectors,
separation of Spot and Futures auth, the no-retry rule for mutations,
cancellation, reconnects, snapshot hooks, and sequence gaps. The
[coverage catalogue](contracts/catalogue.json) records the source page and
fixture behind each entry.

Keep in mind what passing tests actually prove: the SDK matches the docs. Whether
the live exchange matches the docs is a separate question, and that's what the
integration tests are for.

[CI](.github/workflows/ci.yml) runs on Go 1.26.x and 1.27.x with the race
detector, builds, gofmt, vet, pinned Staticcheck and golangci-lint, module
verification, and govulncheck. The
[verification record](docs/VERIFICATION.md) lists which checks were actually
run locally and which live checks were left off.

### Live integration tests (optional)

All `TestIntegration*` tests are skipped unless you opt in. To run the public
REST check, set the flag and the API base URL explicitly:

```sh
export COINW_INTEGRATION=1
export COINW_TEST_BASE_URL=https://api.coinw.com
go test -run '^TestIntegrationPublicREST$' -v .
```

Or in PowerShell:

```powershell
$env:COINW_INTEGRATION = '1'
$env:COINW_TEST_BASE_URL = 'https://api.coinw.com'
go test -run '^TestIntegrationPublicREST$' -v .
```

These send real requests. Private read checks need separate
`COINW_TEST_SPOT_*` or `COINW_TEST_FUTURES_*` credentials, and WebSocket checks
need explicit URLs for the adapters you want to test. The single
order-cancellation check also requires `COINW_TEST_MUTATIONS=1` and
`COINW_TEST_CANCEL_ORDER_ID`, and is always skipped when `CI` is set.

[INTEGRATION.md](docs/INTEGRATION.md) lists every variable and test. Never
commit live credentials or real account responses as fixtures, and unset the
test variables when you're done.

### Project docs

| Document | Use it to |
|---|---|
| [API index](docs/API.md) | Find a method, subscription, source page, or fixture |
| [Contract notes](docs/CONTRACTS.md) | Understand response conflicts, ordering, and snapshot policy |
| [Migration guide](docs/MIGRATION.md) | Move older code to typed responses and requests |
| [Integration guide](docs/INTEGRATION.md) | Run selected live checks |
| [Verification record](docs/VERIFICATION.md) | See the recorded local check results |
| [Contributing](CONTRIBUTING.md) | Set up dev tools and change a contract |
| [Release checklist](docs/RELEASING.md) | Versioning, upstream changes, and publishing |
| [Changelog](CHANGELOG.md) | See what changed in the SDK |
| [Security policy](SECURITY.md) | Report a vulnerability and handle sensitive data |
| [Code of conduct](CODE_OF_CONDUCT.md) | How we work together |

### Package layout

Most code only needs `coinw.Client` plus `model` for request structs, enums, and
values. The other packages are exposed for testing and more specialized setups:

| Package | What it does |
|---|---|
| Root `coinw` | Client construction, options, `Decimal`, service and session factories |
| `services` | REST methods grouped by domain |
| `model` | Request and response records, instrument types, decimal and union types |
| `transport` | HTTP execution, signing dispatch, response validation, read retries |
| `auth` | Separate Spot and Futures signers |
| `endpoints` | Route registry: methods, auth profiles, rate-limit metadata |
| `ratelimit` | Context-aware in-process budgets and the pluggable `Limiter` interface |
| `websocket` | Raw JSON and Socket.IO sessions, typed subscriptions, recovery, sequence tracking |
| `errors` | Typed SDK errors and credential redaction |
