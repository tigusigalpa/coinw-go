# Migrating to typed responses

This unreleased change is an API transition and needs a minor version increment
while the module remains pre-v1. Go 1.26 is now the supported minimum.

REST service methods return `*transport.Response[T]`. Read typed fields directly:

```go
result, err := client.SpotMarket().GetRecentTrades(ctx, "BTC_USDT", 0)
if err != nil { return err }
for _, trade := range result.Data {
    fmt.Println(trade.Price.String(), trade.Amount.String())
}
```

The former raw payload is `result.Result.Data`; the full raw body remains
`result.Envelope`. `result.DecodeData` is retained as a low-level escape hatch,
but the service has already validated its typed schema. Raw private data must
not be logged. Code converting `result.Data` to a string must be updated.

All documented filters are exposed by the corresponding `WithRequest` method
and endpoint request DTO. Optional pointer fields preserve explicit zero values:

```go
instrument, _ := model.USDCInstrument("BTC_USDC")
result, err := client.FuturesMarket().GetKLineOfAnInstrumentWithRequest(ctx,
    model.FuturesMarketGetKLineOfAnInstrumentRequest{
        CurrencyCode: instrument,
        Granuality: "0", // spelling from the official request table
    })
```

`PlaceAnOrder`, `ModifyAnOrder`, batch orders/closes and SL/TP methods now accept
specific request structs in `model`. Arbitrary extra maps are no longer public
request APIs. `WithdrawalRequest` has an explicit `Chain` and `Decimal` amount;
memo is sent even when empty, without inventing memo semantics. Missing required
fields now fail before network I/O. Spot order lookup/cancel parameters use
`orderNumber`.

Use `Decimal` for strict request money. Responses use `WireDecimal` only where
the published string/numeric representations conflict. Its `String()` and JSON
serialization produce plain validated decimal strings; `RawJSON()` preserves
the original token, including exponent notation. `StringOrInteger`,
`ValueOrScalar`, `ObjectOrArray`, `JSONEncoded`, `MaskedDecimal` and
`LeverageOptions` preserve their specific documented alternatives.

Typed WS constructors and `SubscribeStream`/`DecodeStream` or
`SubscribeSocketIO`/`DecodeSocketIO` bind a subscription to its payload type.
The low-level raw adapters remain available for explicit dispatch. Public and
private business profiles must not be mixed. Initialize incremental depth from
a documented snapshot; never use an unanchored increment as the initial book.
Typed private subscriptions require `WithOnConnect` with the caller's matching
REST snapshot read, which the session reruns after each login.

The published interval/origin lists are now copy-returning functions
`services.SpotKLinePeriods()` and `services.FuturesOriginTypes()`; changing a
returned slice cannot mutate SDK validation rules.
