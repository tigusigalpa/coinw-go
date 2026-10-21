# Changelog

## [Unreleased]

### Changed
- Go 1.26 is the maintained minimum; CI covers Go 1.26 and 1.27.
- REST services now return `transport.Response[T]` with strict typed data and
  preserved raw envelopes. See `docs/MIGRATION.md` for API changes.
- Futures order, batch, modification and SL/TP APIs use specific request DTOs;
  withdrawal fields are explicit and monetary values use validated decimals.
- Incremental depth needs a snapshot before applying updates. Resync callbacks
  request the documented snapshot on an absent anchor or a sequence gap.

### Contract coverage and release preparation
- Added exact decimal/exponent, identifier, scalar/object, object/array,
  encoded-JSON, masked-balance and leverage compatibility adapters.
- Added complete request-table DTOs and `WithRequest` methods for published
  filters, typed USDT/USDC instruments, enums and immutable typed WS streams.
- Added a source-linked catalogue for 86 REST and 26 WS rows, dedicated
  published/sanitized fixtures, per-row HTTP success/error tests, mock-WS
  subscription/ack/event tests and fixed independent signature golden vectors.
- Added tests for precision, malformed/missing fields, snapshot/relogin/resubscribe,
  schema conflicts, safe decoding errors, context propagation and race safety.
- Server WS rejections are delivered to the error callback as typed authentication,
  rate-limit or API errors, including code 29001.
- Added environment-gated public/private read integration checks and a separately
  opted-in one-order cancellation check that is always disabled in CI.
- Added executable examples, SECURITY, CONTRIBUTING, CODE_OF_CONDUCT, contract
  notes, integration instructions, semantic versioning and release checklist.
- Added pinned lint/static/security tooling, dependency checks, and injectable
  bounded retry backoff/sleep hooks. Runtime dependencies remain unchanged.
- Fixed Spot `orderNumber` request keys, Futures `granuality` spelling and its
  1..1500 limit, both previously unavailable interface-rate entries, required
  request fields and conditional ID/quantity rules.

### Fixed
- Terminating, idempotent credential redaction in all SDK error strings;
  transport causes remain available through error unwrapping.
- Per-attempt HTTP timeouts, encoded Futures signing queries, strict JSON
  envelopes/content types, bounded response reads and public-token validation.
- No automatic retries for POST reads, matching the original safety contract.
- Missing documented interface rates, IP/UID scope handling, fractional-rate
  bucket capacity and the Group B classification of the batch ticker route.
- WebSocket cancellation, independent heartbeat, serialized writes, fatal
  read timeouts, failed-login cleanup, immutable resubscriptions and Socket.IO
  handshake ordering; added a connection hook for caller-provided snapshots.
- Money validation in flexible request fields and batches; extra fields can
  no longer replace validated arguments or withdrawal destinations.
- Isolated-margin filters (`positionModel=0`) are sent; negative values omit
  the filter. Optional instrument/name filters can now be omitted where verified.
- README WebSocket variable shadowing and explicit implementation limitations.
- Unfiltered Spot cancellation is gated instead of assuming an all-market scope;
  batch SL/TP now requires stopFrom and permits the documented optional instrument.
- A refused WebSocket subscription no longer closes the connection: it goes to
  the error callback, the other streams keep running, and the refused stream is
  dropped from the restore set. Reconnect backoff resets only after a healthy
  frame, so connections that are always rejected now reach `MaxAttempts`.
- `Close` now stops an active `Run`, which returns nil, instead of letting it
  reconnect.
- The rate limiter uses a sliding window, so "N requests/s" can no longer
  burst to about 2N within one second.
- `WithOnConnect` now runs after the subscriptions are sent and before the
  first frame is read, so no event falls between the snapshot and the stream.
- `Subscribe` and `Unsubscribe` change the desired set while disconnected, and
  `Connect` restores it. A failed send keeps the subscription registered.
- Fixed a nil-pointer panic when the connection was dropped between `Connect`
  and the restore of subscriptions.

### Added
- Initial implementation of the CoinW Go SDK (`github.com/tigusigalpa/coinw-go`,
  Go 1.21+, only dependency gorilla/websocket).
- Complete documented REST coverage: Spot (market/account/check/trade),
  Futures (market/assets/common/check/position/trade), Common
  (accounts/sub-accounts/affiliate) as typed services over a shared
  request pipeline.
- Two isolated auth profiles: Spot MD5 signature, Futures HMAC-SHA256
  (headers sign/api_key/timestamp).
- All five documented WS variants: Spot raw-JSON public (Method 2), Spot
  Socket.IO public (Method 1, token fetch + Engine.IO framing), Spot
  private, Futures public, Futures private — login ack, ping keepalive,
  bounded-backoff reconnect + resubscribe.
- Endpoint registry with method/path/auth/rate-group/ambiguity metadata;
  ambiguity gates return AmbiguousDocumentationError.
- Documented rate-limit policy (spot IP+UID, futures groups A/B, affiliate,
  per-interface buckets) — context-aware, injectable/distributed-ready.
- DepthSequenceTracker for the spot incremental order book.
- Decimal-string monetary model; typed error taxonomy with redaction.
- Unit/integration tests with scripted fake transports and dialers, CI.
