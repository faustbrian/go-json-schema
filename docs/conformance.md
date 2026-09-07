# Conformance

The authoritative executable evidence is pinned to JSON Schema Test Suite
revision `c0b038ad7244712cf73650f44e90d0bc5704e8c7`.
Fixture results are interpreted through the explicit
[specification decision register](specification-decisions.md); passing a corpus
does not silently resolve behavior the corpus leaves open.

The separately pinned upstream range
`c9510e3bf8a896c3cba4e08509cf752b4f30dff8...f6fd52a0a95472e079cbfc6ef7f089702b80e045`
adds six malformed URI-template percent-encoding vectors. They are executed
with explicit format assertion across every applicable supported released
dialect; the future v1 copy is retained only as provenance.

| Evidence | Current result |
| --- | ---: |
| Released dialects | 6 |
| Mandatory and optional fixture files | 354 |
| Cases | 8,505 |
| Passes | 8,505 |
| Skips | 0 |
| Failures | 0 |
| Official meta-schema resources | 19 |

`TestOfficialMandatoryFixtures` and `TestOfficialOptionalFixtures` discover
the complete released-dialect trees rather than maintaining an exclusion
list. `TestOfficialAnnotationFixtures` runs every compatible official
annotation case for all six dialects. `TestOfficialBasicOutputFixtures`
validates output against the official 2019-09 and 2020-12 constraints.
`TestOfficialMetaSchemasCompileAgainstTheirDialect` self-compiles all pinned
meta-schema and vocabulary resources.

`specification/official-suite-results.tsv` records revision, dialect, file,
group count, case count, pass, skip, failure, and checksum for every fixture.
`make provenance` regenerates and compares it offline, verifies all 558 suite
files and symlinks, and verifies all 19 meta-schema checksums.

`TestReviewedOfficialFormatVectors` also executes the eight format vectors from
the separately pinned upstream range
`3c25e5f709192aadf67cf7f2eb19771a57131fec...55e23729473f4b629fd9266614280f355cd1b4fc`.
Those cases pass in every applicable released dialect; this focused review does
not claim that the complete official corpus has been repinned from
`c0b038ad7244712cf73650f44e90d0bc5704e8c7`.

`TestReviewedOfficialConformanceVectors` executes the eight distinct behaviors
from the separately pinned upstream range
`55e23729473f4b629fd9266614280f355cd1b4fc...c9510e3bf8a896c3cba4e08509cf752b4f30dff8`.
The fixture records final hashes for all 21 changed test files. Applicable
Draft 3 through Draft 2020-12 cases pass without a behavior change; future v1
copies are provenance-only and are not executed.

`TestReviewedOfficialURITemplatePercentEncodingVectors` executes the six
malformed URI-template percent-encoding behaviors from the separately pinned
one-commit range through `f6fd52a0a95472e079cbfc6ef7f089702b80e045`.
The fixture records final hashes for all five changed test files and preserves
the future v1 copy without executing it.

“Full suite compatibility” means this exact pinned corpus has zero failures
and zero unexplained skips. It does not replace normative review, hostile-input
testing, output correctness, coverage, fuzz, mutation, Bowtie, or release
gates. The stable v1 release claim remains bounded by this pinned evidence and
the gates in [RELEASING.md](../RELEASING.md).

## Decision conformance matrix

| Decision | Authority | Executable evidence | Differential status |
| --- | --- | --- | --- |
| JSONSCHEMA-DEC-001 | `json-schema-drafts-source` | `TestOfficialMetaSchemasCompileAgainstTheirDialect`, `TestOfficialVocabularyFixtures`, `TestCompoundResourcesUseTheirOwnVocabulary`, `TestCompileRejectsUnknownRequiredVocabulary`, `TestVocabularyPolicyHandlesOptionalAndPartialDeclarations`, `TestDialectFeaturePoliciesAreExact` | [Assessed](../specification/differential/maintained-peers.json); the minimized Draft 3 dialect disagreement is a deliberate policy difference. |
| JSONSCHEMA-DEC-002 | `json-schema-drafts-source` | `TestOfficialOptionalCoreFixtures`, `TestRegisteredVocabularyCompilesAndEvaluatesCustomKeyword`, `TestOfficialAnnotationFixtures` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-003 | `json-schema-drafts-source` | `TestFormatAssertionIsExplicitAndCompilerOwned`, `TestOfficialFormatAnnotationFixtures`, `TestOfficialOptionalCoreFormatFixtures`, `TestOfficialFormatAssertionVocabularyFixtures`, `TestStandardFormatsDoNotLeakAcrossDialects`, `TestReviewedOfficialFormatVectors`, `TestReviewedOfficialConformanceVectors`, `TestReviewedOfficialURITemplatePercentEncodingVectors` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-004 | `json-schema-drafts-source` | `TestContentKeywordsAreAnnotationsByDefault`, `TestContentAssertionIsLimitedToDraft7`, `TestOfficialDraft7ContentAssertionFixtures`, `TestContentValidationCoversPermissiveAndStrictBranches`, `TestContentValidationSeparatesSyntaxAndResourceFailures` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-005 | `json-schema-drafts-source` | `TestOfficialAnnotationFixtures`, `TestOfficialArrayAnnotationFixtures`, `TestOfficialAnnotationKeywordFixtures`, `TestAnnotationAndOutputTraversalSkipUnappliedSchemas`, `TestApplicableAnnotationsContinueAfterInapplicableKeywords` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-006 | `json-schema-drafts-source` | `TestOfficialUnevaluatedPropertiesFixtures`, `TestUnevaluatedAndPatternKeywordsPropagateTrackingFailures`, `TestUnevaluatedOutputContinuesAfterEvaluatedEntries` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-007 | `json-schema-drafts-source` | `TestOfficialDynamicReferenceFixtures`, `TestBasicOutputPreservesDynamicReferenceEvaluationPath`, `TestDynamicAnchorCompilationContinuesToMatchingAnchor`, `TestReferenceDepthIsRestoredOnSuccessAndFailure`, `TestAnchorKeywordsDoNotLeakIntoEarlierDialects` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-008 | `json-schema-drafts-source` | `TestReplacingReferenceIgnoresSiblingIdentifier`, `TestOfficialReferenceAndDefinitionFixtures`, `TestVerboseReferenceTraversalContinuesWithSiblingKeywords` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-009 | `rfc8259-source` | `TestCompileRejectsInvalidAndAmbiguousJSON` | [Assessed](../specification/differential/maintained-peers.json); duplicate-member parser disagreements are deliberate policy differences. |
| JSONSCHEMA-DEC-010 | `json-schema-drafts-source` | `TestValidateUsesExactNumberSemantics`, `TestExactNumberComparisonSignEdges`, `TestOfficialNumericFixtures`, `TestOfficialOptionalCoreFixtures` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-011 | `ecma262-source` | `TestPatternUsesECMAScriptLookaroundAndBackreferences`, `TestPatternBacktrackingIsBounded`, `TestOfficialPatternFixtures`, `TestOfficialOptionalRegexFixtures`, `TestRegexFormatCompilationUsesConfiguredByteLimit` | [Assessed](../specification/differential/maintained-peers.json); ECMA-262 engine disagreements are deliberate policy differences. |
| JSONSCHEMA-DEC-012 | `rfc3986-source` | `TestNormalizeURLAppliesRFCIdentityRules`, `TestCompileRejectsEquivalentDuplicateResourceIdentifiers`, `TestMapLoaderUsesNormalizedResourceIdentity`, `TestRemoveDotSegmentsPreservesURIPathStructure` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-013 | `json-schema-drafts-source` | `TestOfficialRemoteReferenceFixtures`, `TestLoaderPanicsAreContainedAndRedacted`, `TestLoaderErrorsAreRedactedAndPreserved`, `TestFSLoaderConfinesResourcesToItsBase`, `TestCompositeLoaderFallsThroughOnlyForMissingResources`, `TestResolutionErrorsRedactURISecrets` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-014 | `json-schema-output-source` | `TestOfficialBasicOutputFixtures`, `TestBasicOutputPreservesReferenceEvaluationPath`, `TestVerboseOutputIncludesEveryEvaluatedKeyword`, `TestVerboseOutputRetainsAnnotationResults`, `TestOutputBoundaryHelpersAreExact` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
| JSONSCHEMA-DEC-015 | `json-schema-test-suite-source` | `TestOfficialMandatoryFixtures`, `TestOfficialOptionalFixtures`, `TestOfficialOptionalCoreFixtures`, `TestOfficialOptionalCoreFormatFixtures`, `TestOfficialOptionalRegexFixtures`, `TestReviewedOfficialFormatVectors`, `TestReviewedOfficialConformanceVectors`, `TestReviewedOfficialURITemplatePercentEncodingVectors` | Not assessed; maintained-peer evidence is tracked separately from official-corpus conformance. |
