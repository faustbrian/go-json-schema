# Source documentation audit

[Documentation map](index.md) · [Repository README](../README.md)
[Next: API guide](api.md)

## Pre-edit scope and smallest sufficient proof

Base `fa376a219063445f0256e4775e33943105d0b457`, detached main-source
worktree. This worker owns comments/docs-only Tier A in this repository; original
`chore/go-1-27-toolchain` and its untracked workspace sum remain untouched.
Outcome: truthful package/API/private rationale for schema compilation,
validation, callbacks, resource retrieval and diagnostic/output boundaries.
Direct consumers are the Bowtie command and pinned published-v2 comparison
harness; neither gains algorithm, dependency or conformance changes.

Reproduce with `git ls-files '*.go'`, `modules.json`, `packages.json` and native
`source-docs inventory/check --json`: 85 tracked Go files, two modules, three
root catalog packages plus one test-only comparison directory. Initial 43
exported omissions are structural input, not a quality score. Important risks
are falsely unconditional immutability/concurrent safety with borrowed custom
callbacks, pre-versus-post-allocation bounds, cooperative cancellation, URI
confinement versus filesystem capabilities, exact JSON/number behavior, dialect
scopes, diagnostic privacy and output traversal multiplicative budgets.

## Complete tracked Go inventory

- `.`: `a_uri_template_termination_test.go`, `annotations_test.go`, `api_edge_test.go`, `benchmark_test.go`, `callback.go`, `callback_internal_test.go`, `compiler.go`, `compiler_internal_test.go`, `compiler_malformed_internal_test.go`, `compiler_test.go`, `conformance_annotation_keyword_test.go`, `conformance_array_annotation_test.go`, `conformance_array_applicator_test.go`, `conformance_cardinality_test.go`, `conformance_combinator_test.go`, `conformance_dependencies_test.go`, `conformance_draft3_applicator_test.go`, `conformance_dynamic_reference_test.go`, `conformance_equality_test.go`, `conformance_format_annotation_test.go`, `conformance_full_test.go`, `conformance_loop_detection_test.go`, `conformance_numeric_test.go`, `conformance_object_applicator_test.go`, `conformance_optional_core_test.go`, `conformance_optional_format_test.go`, `conformance_optional_regex_test.go`, `conformance_pattern_test.go`, `conformance_reference_test.go`, `conformance_type_test.go`, `conformance_unevaluated_properties_test.go`, `conformance_unique_items_test.go`, `conformance_upstream_review_test.go`, `conformance_vocabulary_test.go`, `custom.go`, `dialect.go`, `doc.go`, `documentation_contract_test.go`, `errors.go`, `evaluate.go`, `evaluate_index_error_internal_test.go`, `evaluate_mutation_internal_test.go`, `evaluation_internal_test.go`, `example_test.go`, `formats.go`, `formats_internal_test.go`, `fuzz_test.go`, `idna_normalization_dependency_test.go`, `json.go`, `json_equal.go`, `json_equal_internal_test.go`, `json_internal_test.go`, `keyword_test.go`, `limits.go`, `loader.go`, `loader_internal_test.go`, `loader_test.go`, `metaschema.go`, `metaschema_test.go`, `mutation_edges_internal_test.go`, `output.go`, `output_internal_test.go`, `output_mutation_internal_test.go`, `output_test.go`, `regex.go`, `regex_internal_test.go`, `regex_test.go`, `regexp_unicode_dependency_test.go`, `resolution_internal_test.go`, `resource_identity_internal_test.go`, `resource_identity_test.go`, `specification_manifest_test.go`, `upstream_recent_review_test.go`, `uri.go`, `uri_internal_test.go`, `uri_template_expression_test.go`, `value.go`, `value_test.go`.
- `benchmarks/comparison`: `comparison_test.go`, `differential_test.go`, `public_release_test.go`.
- `cmd/bowtie-json-schema`: `main.go`, `main_test.go`.
- `internal/cmd/conformance-manifest`: `main.go`, `main_test.go`.

## Missing documentation versus open source questions

Repair ResourceLoader/KeywordCompiler/KeywordEvaluator/FormatChecker obligations,
limit/error/output/result field contracts and stale Limits text claiming future
evaluator budgets. Clarify copy/borrow and callback cooperation from source,
not passing concurrency tests alone. Audit every compiler/evaluator/output
branch, loader/value/format/regex/URI helper and all test/fuzz/benchmark/example
variants. Question any blanket immutable, bounded-allocation or privacy claim
that crosses injected callbacks/encoders/loaders; do not invent stronger policy.

Drafts 3,4,6,7,2019-09,2020-12, 15 existing specification decisions, ECMA
regular-expression and RFC URI identity choices have pinned authority and
fixture references. Preserve those official bytes, source locks, dated review
records, decision digests and published-peer observations. This pass may name
existing limitations but does not reinterpret normative requirements or run a
new conformance/differential/performance campaign.

## Generation, variants and selected checks

No generated Go, platform build split or debt markers were found; embedded
official meta-schema directive remains intact. Generation owners are existing
conformance-manifest scripts/command; no official/corpus regeneration is needed
for prose. Commands have existing canonical package comments; test-only
comparison documentation is audited without creating a production package.

The configured root API baseline is full native go-doc, ordered `.`,
`./cmd/bowtie-json-schema`, `./internal/cmd/conformance-manifest`; concatenate
unmodified `go doc -all` stdout and normalize only trailing newlines to one.
The legacy check-api-compat script adds PACKAGE headers and is not the current
Tools recipe. Producer pin is Go1.27.2; module minimum stays1.27.0.

After an explicit serialized Go grant: formatting, all85 exact comment-free
AST/signatures/literals/directives and significant tokens (report raw deltas;
normalize only optional implicit semicolons before closing braces if needed),
parsed/rendered comments, all3 native GoDocs and stable configured API refresh,
native documentation script/snippet selector and existing Examples if relevant.
No peer benchmark module execution, race/fuzz/service/release qualification.
Cheap source/offline-spec/local-link/marker checks, pinned changed-added/staged
secret scan and ordinary complete Tier A truth/usefulness review precede an
explicit-path Intent/ACK safe-wrapper commit. All3 task-owned Go caches,
helpers, generated snippets/renders and secret snapshots are cleaned after
bounded runs on success/failure. Parent owns publication/CI/goal custody.

## Source-grounded contracts and limits

Registry maps and decoded plans are owned/read-only, but callback implementations
are borrowed. Compiler caches are per Compile, not shared runtime caches;
unfinished nodes are cached before reference recursion to preserve graph
identity. Schema callback purity/concurrency is required for repeated validity,
evaluated-location, annotation and diagnostic passes. No callback shutdown or
forced context interruption is owned. Constructor/bundled meta-schema checking
uses default limits independently from configured request budgets.

JSON ingestion keeps decimal text, checks raw UTF-8 and rejects duplicate
members; escaped-surrogate replacement follows encoding/json. Exact symbolic
digits/scale avoid exponent expansion and hash collisions fall back to exact
structural equality. Compile/reference identity and dynamic scope remain
resource/dialect-specific, with existing decision/test ownership named next to
those constraints. Budget counters are not total instruction/heap bounds.
Resource bytes are admitted after loader return; encoder output admission does
not bound encoding/json internal buffers; FS URI confinement is not OS
symlink/device capability confinement.

Existing defects, not standards interpretations: some reference failures quote
raw reference text despite maintained broad redaction claims; explicit-null
annotations decode to nil and disappear under OutputUnit's omitempty tag.
Both are documented narrowly, with no runtime repair or decision/digest change.
The requested disposable observations will distinguish source findings from
actual executed output. Callback errors retain raw causes for errors.Is/As;
JSON parse names, output locations and annotations can also be sensitive.

Source-review coverage includes all 85 tracked paths: ingestion/equality,
compiler/index/reference/vocabulary, evaluation/tracking, output hierarchy/
accounting, format/URI/regex, callback containment, command IO/generation and
every test/helper family in the inventory. Four root fuzz targets, six benchmark
entrypoints (four root, two published-peer) and four root Examples retain their
bodies/Output blocks. Root conformance wrappers bind mandatory/optional keyword
lanes, annotation/dynamic/reference/equality/format corpora and recent reviewed
regressions; they are source-audited, not newly executed as a campaign. Private
capitalized helper fields receive no mechanical boilerplate.

Separate source-audited harness concern: FuzzRawAndValueValidationAgree filters
UTF-8/trailing input but not duplicate member names; standard Go decoding can
collapse duplicates before the value path while raw validation refuses them.
Its agreement assertion can therefore reject an intentional parser-policy
difference. No new fuzz campaign or harness behavior fix belongs in this pass.

Open rationale: no supported source explains why explicit-null annotations are
omitted or why some reference errors bypass identifier redaction. These require
separately scoped behavior/security repair, not invented normative resolutions.

## Completed evidence and explicit remaining boundaries

- Current directory-aware source check: zero findings after useful repairs to
  the initial 43 omissions. This structural result does not replace source
  truth/usefulness review. No generated Go or unsupported private-member
  boilerplate was introduced.
- Go1.27.2 formatting and native scanner/parser against the base: all85
  comment-free AST/signatures/literals/directives and raw significant tokens
  identical, zero raw deltas and no normalization needed; 182 nonempty comment
  groups HTML-rendered. Final five public files plus two private pinned-label
  clarifications received fresh preservation and 68 affected HTML renders.
- Three native `go doc -all` entrypoints rendered. Configured full-prose API
  generation was stable twice in package order; final root rerender/baseline
  assembly reused unchanged command blocks and was stable. Native `golib api
  check` passed all three configured package outputs; comparison is explicitly
  not applicable. Public declarations remain identical; only prose/spacing
  updates the API snapshot.
- Native `scripts/check-docs.sh` passed: Markdown links/anchors/references,
  actionable routes and Go fences, fragment compilation, two standalone
  programs with designated output, four existing Examples and API equality.
  Final comment-only refresh reused unchanged snippet/Example results.
  Native docs-contract Python tests passed, and final static link/fence
  validation passed with task-owned output removed. One auxiliary invocation
  omitted its required output-directory creation; it made no source changes
  and the corrected owned invocation/native script passed.
- Offline `specification check` validated one module/15 unchanged decisions.
  Native official-suite/meta-schema provenance scripts verified 558/19 files.
  These are integrity checks, not a new conformance or peer qualification.
- Disposable read-only observations imported this worktree's actual root:
  reference failure exposed synthetic URI user-info/query/fragment (reported
  only booleans); a default-null annotation produced one unit with nil Go
  Annotation and no explicit-null wire field. These observations establish
  the two existing defects, not repaired behavior or broad conformance.
- Final staged changed-added/new-file pinned-secret check and ordinary complete
  Tier A review apply to the exact explicit paths. Runtime bodies, dependencies,
  corpus/meta-schema inputs, specification authorities/decisions/digests and
  historical peer evidence remain unchanged. No runtime, race/leak, active
  fuzz, benchmark, official-suite, Bowtie or public-consumer campaign ran.

All three exact owned Go caches, snippets/native temp, helpers, probes, renders
and secret snapshots are removed after bounded sessions, success/failure
included. The clean committed detached worktree is retained for parent delivery;
local proof is not remote-main/hosted-CI or release qualification.

An additional source-only interpretation question remains: compile's format
lookup/standardFormatSupported can leave plan.format nil for unrecognized
formats, while JSONSCHEMA-DEC-003 describes mandatory support when an assertion
vocabulary is required. That immutable expected contract is preserved; this
pass does not assert a new standards interpretation, run another probe or
repair the evaluator. The raw/value fuzz duplicate-admission concern above
also remains source-audited rather than executed fuzz evidence.

## Objective Enforcement

The checksum-pinned public Tools v2.4.0 verifier now supplies
`make source-docs-check` as a prerequisite of the existing local and CI Make
routes. Ordinary shared CI selects strict source-documentation checking and
requires its result. Native prerequisites, operation budgets, security and
release controls are retained. Hosted exact-head evidence is recorded by the
coordination ledger; configuration alone does not establish green delivery.
