# Versioning and release checklist

This module follows [Semantic Versioning 2.0.0](https://semver.org/).
Until v1.0, breaking changes require a minor version increment and migration
notes; patch releases preserve the public API and behavior. After v1.0, breaking
Go API, serialization, compatibility or supported-toolchain changes require a
major version and Go module import-path migration where applicable.
Treat a stricter decoder that rejects previously accepted documented responses
as a compatibility change. Release candidates use `vX.Y.Z-rc.N`.

The current working tree is unreleased. No tag, registry publication or remote
release is implied by passing local checks.

## Before tagging

- Recheck the [CoinW sitemap](https://www.coinw.com/api-doc/en/sitemap.xml) and
  [changelog](https://www.coinw.com/api-doc/en/changelog/change_log). Review every
  source linked from `contracts/catalogue.json`, including all WS variants.
- Confirm each route, verb, authentication profile, rate scope, request field,
  response model and fixture. Keep the three route conflicts gated until CoinW
  confirms a contract in writing; record the evidence before changing a gate.
- Review the distinctions in `docs/CONTRACTS.md`; never infer OHLC ordering,
  hidden balances, quantity conversion, fills, or undocumented pagination.
- Regenerate or edit fixtures only from verified documentation. Record syntax
  repairs and synthetic substitutions; scan the diff for secrets/private data.
- Review dependency changes, module checksums, licenses and security advisories.
  Verify that the maintained Go baseline and full CI matrix are still correct.
- Run all commands in `CONTRIBUTING.md` and obtain green CI on both Go series.
  Confirm examples compile, formatting is clean, and integration tests skip
  without opt-in. Live mutation tests must remain absent from release CI.
- If an authorized test account is available, run the separately configured
  read-only integration checks and record date, environment, operation IDs and
  pass/fail only. Record skipped checks honestly; never claim live conformance
  from offline fixtures. Mutations need an explicit local opt-in and order ID.
- Move relevant `CHANGELOG.md` entries into the intended version and date;
  include `docs/MIGRATION.md` for this typed-response API transition. Recheck the
  MIT license, SECURITY and contributor documents.

## Publish and verify

After maintainer approval of a concrete release commit, create the reviewed
`vX.Y.Z` tag and publish release notes describing changes, migration steps,
documented ambiguity gates, tested Go versions and any skipped live checks.
Verify the module can be resolved and built by its version from a clean checkout.
If a release is faulty, publish a corrective version; do not move published tags.

Review the first reports after release, triage security reports privately, and
update the supported-version statement when a new release supersedes it.
