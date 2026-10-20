# Security policy

Only the latest tagged SDK release receives security fixes. Use the latest
patch of one of the Go versions covered by CI (currently 1.26 and 1.27).

Report a suspected vulnerability privately through the repository's
[GitHub security advisory form](https://github.com/tigusigalpa/coinw-go/security/advisories/new).
If private reporting is unavailable, contact the repository owner through
their GitHub profile without disclosing exploit details publicly. There is no
promised response SLA. Never include API credentials, signatures, wallet
addresses, account payloads, or unredacted request URLs in an issue or report.

Include the SDK commit/version, Go version, affected operation ID, a synthetic
reproduction, and expected versus observed behavior. Redact evidence before
sharing it. Revoke exposed exchange credentials at CoinW immediately; updating
the SDK does not revoke a key.

The SDK separates Spot and Futures credentials. Public calls receive neither.
No mutation runs during construction, reconnect, or test setup. Mutations are
never retried automatically. After a timeout, use the corresponding documented
read operation to determine state before deciding whether to resubmit.

`Response.Result.Data`, `Result.Envelope`, WS raw frames, and compatibility
`RawJSON()` values can contain private data. They are for explicit application
inspection and must not be logged. Typed decoder errors contain contract
locations instead of payload values. Error causes retained for `errors.Is/As`
can themselves contain upstream details; log the SDK error string, not an
arbitrary dump of its wrapped fields.

CI uses public synthetic fixtures, checksum verification, static analysis and
the Go vulnerability database. Live integration tests require explicit local
environment variables. Mutation integration tests are additionally disabled
when `CI` is set. Vulnerability database availability is required for the
security job; failures must be investigated rather than suppressed.
