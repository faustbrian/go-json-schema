# Normative and conformance inputs

[`manifest.json`](manifest.json) is the structured provenance index for the
official suite, published meta-schemas, and Bowtie evidence. Its source pins
and referenced local digests are checked by
`TestSpecificationManifestPinsEveryConformanceSource`.

## Upstream review history

### 2026-09-03

- JSON Schema specification `main` from
  `499eba5749b0a22940e15660dafe50b74df05cb9` through
  `0932747f3f3128758f3166e0d3e23e0b8d1025ee` changes future v1 meta-test
  infrastructure and removes an illustrative trailing comma. The immutable
  Draft 3 through Draft 2020-12 sources are unchanged, so the review is
  behavior-neutral for every decision bound to those sources.
- JSON Schema Test Suite
  `3c25e5f709192aadf67cf7f2eb19771a57131fec...55e23729473f4b629fd9266614280f355cd1b4fc`
  adds the eight format vectors pinned under `testdata/regressions`. Their
  successful explicit-assertion result binds the format-assertion and
  optional-suite decisions without changing the complete-corpus pin.

## Decision conformance matrix

| Decision | Authority | Executable evidence | Differential status |
| --- | --- | --- | --- |
| JSONSCHEMA-DEC-001 | `json-schema-drafts-source` | `TestOfficialMetaSchemasCompileAgainstTheirDialect`, `TestOfficialVocabularyFixtures`, `TestCompoundResourcesUseTheirOwnVocabulary`, `TestCompileRejectsUnknownRequiredVocabulary`, `TestVocabularyPolicyHandlesOptionalAndPartialDeclarations`, `TestDialectFeaturePoliciesAreExact` | [Assessed](differential/maintained-peers.json); the minimized Draft 3 dialect disagreement is a deliberate policy difference. |
| JSONSCHEMA-DEC-002 | `json-schema-drafts-source` | `TestOfficialOptionalCoreFixtures`, `TestRegisteredVocabularyCompilesAndEvaluatesCustomKeyword`, `TestOfficialAnnotationFixtures` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-003 | `json-schema-drafts-source` | `TestFormatAssertionIsExplicitAndCompilerOwned`, `TestOfficialFormatAnnotationFixtures`, `TestOfficialOptionalCoreFormatFixtures`, `TestOfficialFormatAssertionVocabularyFixtures`, `TestStandardFormatsDoNotLeakAcrossDialects`, `TestReviewedOfficialFormatVectors` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-004 | `json-schema-drafts-source` | `TestContentKeywordsAreAnnotationsByDefault`, `TestContentAssertionIsLimitedToDraft7`, `TestOfficialDraft7ContentAssertionFixtures`, `TestContentValidationCoversPermissiveAndStrictBranches`, `TestContentValidationSeparatesSyntaxAndResourceFailures` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-005 | `json-schema-drafts-source` | `TestOfficialAnnotationFixtures`, `TestOfficialArrayAnnotationFixtures`, `TestOfficialAnnotationKeywordFixtures`, `TestAnnotationAndOutputTraversalSkipUnappliedSchemas`, `TestApplicableAnnotationsContinueAfterInapplicableKeywords` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-006 | `json-schema-drafts-source` | `TestOfficialUnevaluatedPropertiesFixtures`, `TestUnevaluatedAndPatternKeywordsPropagateTrackingFailures`, `TestUnevaluatedOutputContinuesAfterEvaluatedEntries` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-007 | `json-schema-drafts-source` | `TestOfficialDynamicReferenceFixtures`, `TestBasicOutputPreservesDynamicReferenceEvaluationPath`, `TestDynamicAnchorCompilationContinuesToMatchingAnchor`, `TestReferenceDepthIsRestoredOnSuccessAndFailure`, `TestAnchorKeywordsDoNotLeakIntoEarlierDialects` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-008 | `json-schema-drafts-source` | `TestReplacingReferenceIgnoresSiblingIdentifier`, `TestOfficialReferenceAndDefinitionFixtures`, `TestVerboseReferenceTraversalContinuesWithSiblingKeywords` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-009 | `rfc8259-source` | `TestCompileRejectsInvalidAndAmbiguousJSON` | [Assessed](differential/maintained-peers.json); duplicate-member parser disagreements are deliberate policy differences. |
| JSONSCHEMA-DEC-010 | `json-schema-drafts-source` | `TestValidateUsesExactNumberSemantics`, `TestExactNumberComparisonSignEdges`, `TestOfficialNumericFixtures`, `TestOfficialOptionalCoreFixtures` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-011 | `ecma262-source` | `TestPatternUsesECMAScriptLookaroundAndBackreferences`, `TestPatternBacktrackingIsBounded`, `TestOfficialPatternFixtures`, `TestOfficialOptionalRegexFixtures`, `TestRegexFormatCompilationUsesConfiguredByteLimit` | [Assessed](differential/maintained-peers.json); ECMA-262 engine disagreements are deliberate policy differences. |
| JSONSCHEMA-DEC-012 | `rfc3986-source` | `TestNormalizeURLAppliesRFCIdentityRules`, `TestCompileRejectsEquivalentDuplicateResourceIdentifiers`, `TestMapLoaderUsesNormalizedResourceIdentity`, `TestRemoveDotSegmentsPreservesURIPathStructure` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-013 | `json-schema-drafts-source` | `TestOfficialRemoteReferenceFixtures`, `TestLoaderPanicsAreContainedAndRedacted`, `TestLoaderErrorsAreRedactedAndPreserved`, `TestFSLoaderConfinesResourcesToItsBase`, `TestCompositeLoaderFallsThroughOnlyForMissingResources`, `TestResolutionErrorsRedactURISecrets` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-014 | `json-schema-output-source` | `TestOfficialBasicOutputFixtures`, `TestBasicOutputPreservesReferenceEvaluationPath`, `TestVerboseOutputIncludesEveryEvaluatedKeyword`, `TestVerboseOutputRetainsAnnotationResults`, `TestOutputBoundaryHelpersAreExact` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-015 | `json-schema-test-suite-source` | `TestOfficialMandatoryFixtures`, `TestOfficialOptionalFixtures`, `TestOfficialOptionalCoreFixtures`, `TestOfficialOptionalCoreFormatFixtures`, `TestOfficialOptionalRegexFixtures`, `TestReviewedOfficialFormatVectors` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |


The official JSON Schema Test Suite is pinned to commit
`c0b038ad7244712cf73650f44e90d0bc5704e8c7`, committed upstream on
2026-07-14. The upstream repository is
<https://github.com/json-schema-org/JSON-Schema-Test-Suite> and is licensed
under the MIT License included in the vendored tree.

`scripts/sync-official-suite.sh` retrieves the immutable GitHub archive,
checks its SHA-256 digest, and imports it without modifying fixture bytes.
The generated `official-suite.sha256` records every vendored file. The
default `make provenance` command runs entirely offline and rejects missing,
added, or changed files.

There are no local deviations from the pinned suite. Project regression
fixtures MUST be stored under `testdata/regressions`, never in
`testdata/official`.

Reviewed upstream changes that do not justify repinning the complete corpus
remain isolated under `testdata/regressions`. The fixture
`json-schema-test-suite-3c25e5f-to-55e2372.json` pins the exact two-commit
upstream range, changed-source digests, decision bindings, and eight new format
vectors reviewed on 2026-09-03. `TestReviewedOfficialFormatVectors` executes
the seven RFC 5321 mailbox cases for Drafts 4, 6, 7, and 2019-09 and the valid
Draft 7 IDN-hostname case with explicit format assertion.

`official-suite-results.tsv` inventories every released-dialect mandatory and
optional fixture. Each row records its group and case count, checksum, and the
zero-skip, zero-failure result enforced by `TestOfficialMandatoryFixtures` and
`TestOfficialOptionalFixtures`. The offline provenance check regenerates the
manifest and rejects silent case-count reductions.

The released-dialect meta-schemas and their vocabulary meta-schemas are
pinned by immutable dialect URI and SHA-256 digest in
`official-meta-schemas.sources.tsv` and `official-meta-schemas.sha256`.
They were retrieved from the official JSON Schema publication endpoints on
2026-07-19 without modifying their bytes. Those specification artifacts use
the JSON Schema specification project's BSD 3-Clause or Academic Free
License 3.0 terms. The upstream license is published at
<https://github.com/json-schema-org/json-schema-spec/blob/main/LICENSE>.

`scripts/check-official-meta-schemas.sh` verifies the complete bundle
offline. To update it, retrieve every URI in the source manifest, review the
published dialect and license changes, replace the corresponding files, and
regenerate the checksum manifest. A checksum update MUST NOT be used to hide
a local modification or conformance failure.

## Updating the pin

1. Review upstream changes between the old and proposed revisions.
2. Update the revision and archive digest in `official-suite.env`.
3. Remove the existing vendored suite intentionally.
4. Run `scripts/sync-official-suite.sh` with network access.
5. Review fixture and case-count changes before committing them.
6. Run every per-dialect conformance lane and `make provenance`.

The pin MUST NOT be updated merely to remove or avoid a failing case.
