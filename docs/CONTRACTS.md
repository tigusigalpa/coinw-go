# Verified contract snapshot: 2026-10-20

`contracts/catalogue.json` is the machine-readable mapping for the original Go
catalogue: 86 REST operations and 26 WS subscription variants across the five
session profiles. Each row links its English documentation source, typed API,
payload model and dedicated example fixture. HTTP contract tests exercise the
public service methods against `httptest`; WS tests exercise subscriptions,
unsubscriptions, acknowledgements and typed events through a scripted peer.

Response fixtures come from published examples. Python booleans/quotes become
JSON, omission markers and redundant commas are removed, a missing closing
brace and a printed newline in a key are repaired. Masked numeric IDs use a
synthetic integer above 2^53. Address, API-key and whitelist strings use
`fixture-only`. No fixture comes from a live account. These syntax repairs do
not authorize choosing an undocumented business meaning. The Futures depth
fixture is explicitly synthetic because its printed example is malformed.

## Strict values

`Decimal` validates plain decimal strings for requests. `WireDecimal` is the
named response adapter for published string/numeric/exponent conflicts. It
expands exponents exactly using string operations, serializes as a decimal
string, and retains its original JSON token. An 8192-byte token and exponent
range -4096..4096 bound decoding work; no precision is rounded or truncated.
Out-of-range inputs are protocol errors. Monetary values never use floats.

Batch trade acknowledgements retain each item's `msgCode` and identifiers.
Inspect every item for partial failures; a successful envelope does not imply
that every batch item succeeded or that any returned order identifier is a fill.

`StringOrInteger` preserves the original token for identifiers, conflicting
statuses and timestamps. It does not interpret opaque strings or timestamp
units. `ValueOrScalar[T]`, `ObjectOrArray[T]`, `JSONEncoded[T]`, `MaskedDecimal`
and `LeverageOptions` expose specific alternatives rather than guessing a
single shape. Optional pointers distinguish absent/null fields from present
zero values. Unknown typed fields, missing required structure and malformed
tuples fail before a service returns a result. The original envelope remains
available for explicit investigation and must not be logged when private.

## Unresolved documentation conflicts

| Source / field | SDK behavior |
|---|---|
| Sub-account max transferable amount: `/api/v1/private` versus linked `/api/v1/api/v1/private` | Registered typed scalar response; sending remains gated. |
| Spot batch historical orders: host/path disagreement | Registered object/array response union; sending remains gated. |
| Futures modify batch SL/TP: `/updateTpsl` versus `/TPSL` | Typed request and acknowledgement; sending remains gated. |
| Spot K-Line endpoint/table versus hyperlink `period` example | Use `returnChartData`, published seconds enum and `currencyPair`; never rewrite to an unsupported period. |
| Spot cancel-all optional `currencyPair` versus specified-instrument description | Require an explicit instrument; the unfiltered request is gated. |
| Address declared integer, example string; memo meanings unclear | Preserve address string/integer form; pass memo literally, without inventing a format or meaning. |
| Futures transferable balance, leverage, margin rate, order acknowledgements and pending count | Preserve scalar versus `{value, ts?}` via `ValueOrScalar`; order IDs never mean fills. |
| Futures REST K-Line table high/open ordering differs from values printed in example | `FuturesCandleTuple` keeps timestamp and five positional decimal slots. WS candles have a separate documented OHLC model. |
| Futures instrument leverage string versus integer array; config JSON string versus object | `LeverageOptions` and `JSONEncoded[FuturesInstrumentConfig]` retain raw alternatives. |
| Affiliate contract balance hidden as `"-"` or `""`; historical position `profitRate` includes `%` | Separate `MaskedDecimal` and `PercentageText`; no invented balance or fraction conversion. |
| Futures depth `ask` versus `asks`; `n`/`ts` locations in malformed example | Retain both labels and outer/inner metadata independently. Never merge sides or relocate timestamps implicitly. |
| Mark-price push describes `p` as index price; funding `nt` timestamp semantics unclear | Expose literal `P` and `NT`; do not assign a new meaning or units. |
| Historical public trades page/pagination descriptions | Expose published tokens and controls only; no invented cursor arithmetic or auto-pagination. |

Only verified CoinW clarification should remove a route gate. Update models,
source evidence, fixtures and migration notes together. A changed live schema is
an error to review, not permission to silently accept arbitrary payloads.

## Snapshot and ordering policy

Private streams do not promise an initial snapshot. Supply `WithOnConnect` to
fetch the documented REST read for your selected subscription after login and
before resubscription on every connection. A failing hook aborts that connection.
Typed private subscriptions reject a missing hook. This makes snapshot selection
explicit where the REST read requires caller-selected instruments or filters.
The application chooses the scope and reconciles events using documented state
or sequence information; timestamps alone do not establish ordering.

Incremental Spot depth cannot initialize a book. `DepthNeedsSnapshot` and
`DepthGap` request the documented `depth_snapshot` channel through the caller's
`TrackAndResync` callback. Only `Reset(snapshotSeq)` anchors the tracker; stale
updates are duplicates and a gap never advances the accepted sequence. Buffer
or discard increments until a snapshot is available. Do not substitute a REST
snapshot without a compatible documented sequence anchor.

All raw/Socket.IO callbacks run synchronously for bounded delivery pressure.
Keep them short and context-aware. Heartbeats run independently. Cancellation
closes the connection; no automatic trade, transfer or withdrawal can occur on
reconnect. Custom transports and callbacks must support concurrency as described
in `CONTRIBUTING.md`.
