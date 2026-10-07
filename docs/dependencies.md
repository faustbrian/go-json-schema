# Dependency review

## `golang.org/x/net/idna`

- Purpose: IDNA2008, UTS #46 mapping, registration validation, bidi rules,
  contextual rules, Punycode canonicalization, and DNS-length enforcement for
  `hostname`, `idn-hostname`, `email`, URI, and IRI formats.
- Necessity: the standard library does not provide equivalent IDNA processing;
  a local table implementation would duplicate security-sensitive Unicode
  standards and update work.
- Ownership and maintenance: maintained by the Go project.
- License: BSD 3-Clause; recorded in `NOTICE`.
- Security: covered by dependency review, `govulncheck`, hostile official IDN
  fixtures, and the package's format budgets.
- Replacement: prefer a standard-library implementation if one becomes
  available with equivalent conformance; otherwise update deliberately after
  reviewing Unicode and behavior changes.

`golang.org/x/text` is a transitive requirement of `x/net/idna` under the same
Go project maintenance and BSD license model.

## `github.com/dlclark/regexp2/v2`

- Purpose: execute schema patterns with ECMAScript lookaround,
  backreferences, character classes, and Unicode behavior unavailable in
  Go's RE2 engine.
- Necessity: JSON Schema defines its pattern syntax in terms of ECMA-262;
  translating only the RE2-compatible subset creates observable validation
  divergences.
- Ownership and maintenance: an independently maintained pure-Go project,
  pinned by `go.mod` and reviewed on update.
- License: MIT for the engine, plus Unicode License v3 for supplementary
  generated Unicode data; recorded in `NOTICE` and the accompanying
  [complete Unicode notice](../UNICODE-LICENSE.txt).
- Security: every compiled expression has a caller-configurable backtracking
  stack bound and match timeout. Pattern count and bytes are bounded before
  compilation.
- Replacement: prefer a standard-library ECMAScript engine if one becomes
  available with equivalent Unicode and resource-limit behavior.

### Supplementary Unicode data provenance

The pinned regexp2/v2 v2.8.2 source revision is
`9b134016f0c3527b1534d3b8fccdcb6e956fc987`. Its generator supplements
Go-standard-library Unicode tables with aliases, binary properties, and
Script_Extensions from the official
[Unicode Character Database 17.0.0 archive](https://www.unicode.org/Public/17.0.0/ucd/UCD.zip).
The generator verifies the archive SHA-256:
`2066d1909b2ea93916ce092da1c0ee4808ea3ef8407c94b4f14f5b7eb263d28e`.

The upstream generated file references a notice absent from the v2.8.2
module archive; [upstream issue 120](https://github.com/dlclark/regexp2/issues/120)
tracks that packaging defect. This distribution carries the complete official
Unicode License v3 notice in `UNICODE-LICENSE.txt`, while retaining the engine's
MIT attribution. The notice snapshot was retrieved on 2026-10-07 from
[Unicode's official license](https://www.unicode.org/license.txt), SHA-256
`e7a93b009565cfce55919a381437ac4db883e9da2126fa28b91d12732bc53d96`.
The UCD input files retain their Unicode, Inc. 2025 data attribution; the
current license notice is not represented as a historical 2025 snapshot.
Any distribution carrying this supplementary data must retain the Unicode
notice or its associated documentation. This owned notice does not repair
the upstream archive's missing file.

No dependency may add implicit network work or a mutable behavior registry.

## Automated review

`make dependencies` verifies module content, tidiness, and the complete build
graph. `make license` checks dependency licenses with a pinned `go-licenses`,
and `make secrets` scans the package tree with a pinned, redacting Gitleaks.
`make workflows` validates the owned workflows and rejects mutable action
references, missing top-level permission declarations, and
`pull_request_target`. Pull requests also run GitHub's dependency-review
action at an immutable commit. `make supply-chain` runs the dependency,
license, and secret gates together.
