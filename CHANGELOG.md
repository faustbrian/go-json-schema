# Changelog

All notable changes will be documented here. The project follows Keep a
Changelog structure and semantic versioning after v1.

## Unreleased

### Changed

- Verify the eight newly reviewed official-suite format vectors across their
  applicable released dialects, pin the exact upstream range and changed-file
  identities, and record that the review does not repin the complete corpus.
  Updated decisions:
  - JSONSCHEMA-DEC-003 sha256:d34404aa2d58558e1da69c87a0a0aeb1129289e15c4c7b1eb1ef31e8bd8b4060
  - JSONSCHEMA-DEC-015 sha256:87457420272991078cd2e6455e586dcffd861789e54e2949a115eddcd4c54a87

- Record the reviewed JSON Schema specification `main` range through
  `0932747f3f3128758f3166e0d3e23e0b8d1025ee` as behavior-neutral for the
  immutable Draft 3 through Draft 2020-12 contracts.

- Publish complete schema-v2 cohesion metadata for the public JSON Schema
  module, plus engineering classification and versioned ecosystem navigation
  for its comparison harness.

- Adopt the checksum-verified `go-library-tools` v1.3.0 CLI, add the local
  `make cohesion` validation entry point, and pin reusable-workflow cohesion
  enforcement to its final immutable revision.

- Adopt the released `go-library-tools` v1.2.0 CLI and immutable merged
  workflow at `1f9629e5f27418600460b55a50a5b2fc81697fab` so local and
  hosted verification enforce the same specification-governance contract.

### Documentation

- Make all 15 entries in the
  [specification decision register](docs/specification-decisions.md)
  machine-auditable with exact authority,
  conformance, fuzz, change-control, and source-monitoring bindings. The
  maintained-peer differential lanes for dialect selection, duplicate JSON
  members, and ECMA-262 regular expressions are assessed with reproducible
  pinned comparisons; the remaining applicable lanes stay explicitly
  unassessed.

  - JSONSCHEMA-DEC-001 sha256:147be7d9160c48ba03e887c03a04022e6c23aded214075338148c08f98d12823
  - JSONSCHEMA-DEC-002 sha256:6210f50357e06748fc6e5b93572ea3ffe0709a0ee47294e3d2f2cbf230d4e0b6
  - JSONSCHEMA-DEC-003 sha256:b55d88fdcaae555feef51f852602384722cd573f251d8ea1b206cb7876daa95f
  - JSONSCHEMA-DEC-004 sha256:029aeea98e782a554e1e3f8c569794a42be1d0f38fb9fede70181df30d437d6d
  - JSONSCHEMA-DEC-005 sha256:1f7d44b92192d7178e257f241deb8913a2b635f6c3f7b37b093a9fdd7be20870
  - JSONSCHEMA-DEC-006 sha256:2c9a0efa21e005b47b0c5127c79b7b5ba0ad93928e967f9a28887916e5d61503
  - JSONSCHEMA-DEC-007 sha256:4e19e022adfe6601d62913e585c55eee3da9a13bda5d09eb0dec02677271cdbb
  - JSONSCHEMA-DEC-008 sha256:3a1485ac1d7c2aa1588a727af61938c8fa40fa81fb09f2f34c4105f863a5f8b9
  - JSONSCHEMA-DEC-009 sha256:02e05d5543725ba19dd7fa92a2f980ead34885c3853fcf2d0d777102c8f2e07b
  - JSONSCHEMA-DEC-010 sha256:a199c275c546249c16b5f589e4fd8f83cefcf0fc5bbe31e50f1ea2e92dc78abf
  - JSONSCHEMA-DEC-011 sha256:e9b1b1947bb158f6e2c420e59a751ecf3e6bf82dd882db7b6be2f61148539b50
  - JSONSCHEMA-DEC-012 sha256:968e55e377d7c6894a75b8f9ffdfb9d83f620f0cc2cea70385e2531e09b9fd9b
  - JSONSCHEMA-DEC-013 sha256:043803367f42b081e6818ceabbcb0da1384736b739fc478af417164612612ea1
  - JSONSCHEMA-DEC-014 sha256:0532c816848d9effab3dac3d59f20bb49042cccf5015b237d7af574106f03d93
  - JSONSCHEMA-DEC-015 sha256:1d253317c7fab90f94846398ab2eed103bec0361e3f5832b72800942e4edbb39

- Link the root README to the normative specification decision register.

- Align maintained documentation with the stable v1 release contract and
  remove obsolete candidate or verdict wording.

- Replace archived monorepo links and completed planning artifacts with a
  standalone documentation index and current repository references.
- Document standalone repository tags instead of the obsolete monorepo tag
  convention.

## 1.0.0 - 2026-08-25

### Changed

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.
- Align the comparison module with the current standalone root archive
  checksum used by repository-local and CI verification.
- Make CodeQL resolve the comparison module through that same root archive.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link contributors directly to the normative specification decision register.

- Correct stale package, standalone, and authoritative-source links in public
  documentation.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-json-schema` identity while preserving its documented API and behavior.

### Documentation

- Add a machine-validated provenance manifest that binds the official test
  suite, published meta-schemas, and Bowtie interoperability reports to their
  existing checksummed evidence.
- Add an auditable specification-decision register covering dialects,
  vocabularies, annotations, references, exact data-model behavior, secure
  resource loading, output, and optional conformance policy.
- Restrict opt-in content assertion to Draft 7 so Draft 2019-09 and Draft
  2020-12 content annotations cannot change the enclosing validation result.
- Reject literal brackets in URI and IRI paths, queries, and fragments while
  retaining bracketed IP literals in authorities.
- Short-circuit Draft 3 `type` and `disallow` schema alternatives after a
  decisive match so validation and diagnostic output evaluate the same branch.
- Expose official-suite conformance as an explicit repository gate distinct
  from ordinary tests and interoperability harnesses.
- Contain and redact custom keyword compiler, keyword evaluator, and format
  checker panics behind the typed `ErrCallbackPanic` boundary.
- Extend panic containment to resource loaders, supplied filesystems, and
  caller-provided JSON marshalers.
- Redact callback error text while retaining the original error for
  `errors.Is` and `errors.As` inspection.
- Reject malformed and duplicate schema resource identifiers and ambiguous
  duplicate anchors instead of silently ignoring or overwriting them.
- Resolve references to the pinned official meta-schema bundle without
  requiring callers or Bowtie registries to duplicate those resources.
- Replace quadratic `uniqueItems` scans for large arrays with canonical,
  collision-safe hash buckets while retaining direct small-array checks.
- Normalize RFC 3986 resource identity across scheme and host case, default
  ports, dot segments, and percent-encoded unreserved characters.
- Prevent `$anchor`, `$recursiveAnchor`, and `$dynamicAnchor` semantics from
  leaking into dialects where those keywords are not defined.
- Apply `MaxRegexBytes` to asserted `regex` format values before compiling
  attacker-controlled expressions.
- Restrict every built-in format name to the dialect that defines it while
  preserving explicit application-defined format extensions.
- Execute schema regular expressions with bounded ECMAScript and Unicode
  semantics, including lookaround and backreferences.
- Add executable examples, official-fixture-backed fuzz corpora, and separate
  compile, validation, reference, and adversarial scaling benchmarks.
- Scope built-in and custom vocabulary activation to each schema resource in
  compound documents, and reject unindexed reference sources without panic.
- Add local release gates and monorepo CI for per-dialect conformance,
  coverage, race, fuzz, mutation, analysis, API, docs, and releases.
- Emit deterministic keyword diagnostics with by-reference evaluation paths,
  condensed Detailed output, complete uncondensed Verbose hierarchies, and a
  dedicated retained-annotation API.

- Added exact JSON parsing, all six released dialects, complete pinned-suite
  discovery, references and dynamic scope, vocabulary processing, official
  meta-validation, standard formats and content policy, annotations, standard
  output, custom extensions, secure loaders, explicit limits, and Bowtie
  protocol support.
- Added pinned official fixture and meta-schema provenance plus a generated
  zero-skip, zero-failure manifest for 8,505 cases.
