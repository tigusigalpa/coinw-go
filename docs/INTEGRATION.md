# Optional integration tests

Normal `go test ./...` runs offline and skips every `TestIntegration*` test.
The local contract suite covers every mapped REST operation and WS variant;
it does not prove that a live CoinW deployment matches the documentation.

For read-only public REST checks, set these explicitly in your own shell:

```sh
export COINW_INTEGRATION=1
export COINW_TEST_BASE_URL=https://api.coinw.com
go test -run '^TestIntegrationPublicREST$' -v .
```

On PowerShell, use `$env:COINW_INTEGRATION='1'` and
`$env:COINW_TEST_BASE_URL='https://api.coinw.com'` before the same test command.
The example URL is the documented production API; use an approved test/proxy
environment where available. Opt-in grants permission to make the selected
real network calls. The suite does not infer a testnet route.

`TestIntegrationPrivateReadOnly` additionally uses
`COINW_TEST_SPOT_API_KEY` / `COINW_TEST_SPOT_SECRET_KEY` and/or
`COINW_TEST_FUTURES_API_KEY` / `COINW_TEST_FUTURES_SECRET_KEY`.
Each profile is independent; missing profiles skip their subtest. Use dedicated
keys with the least available permissions and never paste values into logs.

`TestIntegrationPublicWebSockets` requires each selected adapter's URL:

| Variable | Documented endpoint / purpose |
|---|---|
| `COINW_TEST_SPOT_WS_URL` | `wss://ws.futurescw.com`, Spot raw JSON |
| `COINW_TEST_FUTURES_WS_URL` | `wss://ws.futurescw.com/perpum`, Futures raw JSON |
| `COINW_TEST_SOCKETIO_URL` | `wss://ws.futurescw.info`, Spot Socket.IO |
| `COINW_TEST_TOKEN_URL` | `https://www.coinw.com/pusher/public-token` |

The Socket.IO subtest requires both its socket URL and token URL. These tests
subscribe to public data, decode a typed event, and close within a deadline.
They do not open private streams or mutate account state.

## Mutation test: separate opt-in

`TestIntegrationMutationCancelSpecifiedOrder` can cancel exactly one existing
Futures order explicitly selected by the caller. In addition to the general
opt-in, base URL and Futures credentials, it requires `COINW_TEST_MUTATIONS=1`
and `COINW_TEST_CANCEL_ORDER_ID`. It always skips if `CI` is set. Invoke it only
after selecting an expendable order and understanding that cancellation changes
live state. The suite never places an order to prepare this test.

There are no live withdrawal, transfer, margin, order-creation, close-position
or reverse-position integration tests. Their request/response and no-retry
contracts are tested offline. No live tests were needed to produce fixtures.
Remove all test environment variables after the session.
