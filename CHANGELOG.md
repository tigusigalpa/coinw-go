# Changelog

## [Unreleased]

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
