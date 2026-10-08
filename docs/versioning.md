# Compatibility, deprecation, and versioning

The [specification decision register](specification-decisions.md) is part of
the observable compatibility contract. A change to a resolved parsing,
resolution, validation, annotation, or output decision requires compatibility
review even when the earlier behavior was not documented elsewhere.

The module follows semantic versioning from `v1.0.0`. Public API
changes remain possible but require changelog entries, migration notes, and
executable contract updates.

Released dialect semantics are normative behavior and are not changed under a
minor version to follow a newer draft. New released dialects require explicit
constants and lanes. Experimental unreleased drafts, if introduced, use a
separate opt-in API and are never included in stable compliance totals.

Removing a public option, loader, error classification, output field, or
extension interface requires a major version after v1. Deprecations remain
documented for at least one minor release when security does not require
immediate removal. Tightening a default security limit is documented as an
operational compatibility change.

The pre-v1 URI identity correction normalizes equivalent resource and loader
keys. A custom loader that keyed resources by a non-normalized spelling must
either normalize its registry or index the normalized identifier passed to
`Load`. Two `MapLoader` keys that previously coexisted but normalize to one URI
now fail construction with `ErrResourceUnavailable`.

Module releases use repository tags such as `v1.0.0`. The release process is
in [RELEASING.md](../RELEASING.md).

## Migrating to v2

The v2 module requires Go 1.27 and uses the import path
`github.com/faustbrian/go-json-schema/v2`. Upgrade the toolchain and update
imports before adopting `v2.0.0`; no version-specific source directory is used.
The exported API is unchanged, but stored regular expressions should be
reviewed against the [ECMAScript migration guidance](dialects.md#regular-expression-migration).
URI-template grammar and canonical hostname handling now enforce the selected
normative contracts.

The internal comparison harness retains its published v1 dependency during
root release preparation. That is a v1 consumer baseline, not v2 qualification;
the harness will adopt v2 through public module resolution after publication.

Pending direct-consumer adoption targets also include CloudEvents
`adapters/golib` (through `adapters/jsonschema`), OpenAPI root and
`interoperability`, and Schema Registry root and `providers/confluent`.
They currently consume v1; no v2 reverse dependency is declared until that
consumer has adopted and qualified the public new major.
