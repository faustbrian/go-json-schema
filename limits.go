package jsonschema

import "fmt"

// Limits supplies positive per-document, compilation and evaluation budgets.
// Counters bound selected owned operations, not every CPU instruction, total
// heap allocation or callback IO. Loader/encoder allocations can precede
// admission; cancellation and external work remain cooperative. Output phases
// share evaluation counters and can repeat callback calls.
type Limits struct {
	// MaxInputBytes caps bytes in each parsed document, default 16 MiB; value encoding checks writer output, not encoder-internal allocation.
	MaxInputBytes int
	// MaxNestingDepth caps parser value recursion with root depth one; default 256.
	MaxNestingDepth int
	// MaxTotalValues caps values per parsed document, excluding member names; default 1000000.
	MaxTotalValues int
	// MaxObjectMembers caps members of each object, not aggregate object cardinality; default 1000000.
	MaxObjectMembers int
	// MaxArrayItems caps items of each array; default 1000000.
	MaxArrayItems int
	// MaxNumberBytes caps exact number token text, including exponent, after token decoding; default 4096.
	MaxNumberBytes int
	// MaxSchemaResources caps the root plus newly loaded documents per Compile, not every embedded resource identity; default 1024.
	MaxSchemaResources int
	// MaxTotalSchemaBytes caps root plus returned loaded-document bytes per Compile; default 64 MiB, checked after loaders return.
	MaxTotalSchemaBytes int
	// MaxEvaluationOps caps counted evaluation/output visits per public call; default 1000000, not instruction or elapsed-time bounds.
	MaxEvaluationOps int
	// MaxUniqueComparisons caps pair/hash-bucket comparison accounting for uniqueItems; default 1000000, not each recursively compared value.
	MaxUniqueComparisons int
	// MaxFormatChecks caps built-in/custom asserted format invocations per evaluation/output call; default 100000.
	MaxFormatChecks int
	// MaxSchemaNodes caps cached compiled schema plans per Compile; default 100000.
	MaxSchemaNodes int
	// MaxReferenceDepth caps nested reference evaluation, not reference graph compilation depth; default 256.
	MaxReferenceDepth int
	// MaxDynamicScopeDepth caps checked reference dynamic scope during evaluation; default 256.
	MaxDynamicScopeDepth int
	// MaxCombinatorBranches caps aggregate admitted schema-array entries per Compile; default 100000.
	MaxCombinatorBranches int
	// MaxRegexCount caps compiled schema pattern occurrences per Compile; default 10000, not asserted regex-format calls.
	MaxRegexCount int
	// MaxRegexBytes caps each schema pattern or asserted regex-format string before compilation; default 1 MiB.
	MaxRegexBytes int
	// MaxRegexBacktracking sets the regexp2 backtracking stack bound per match; default 100000.
	MaxRegexBacktracking int
	// MaxRegexMatchMilliseconds sets the regexp2 approximate match timeout; default 100 milliseconds, not a context deadline.
	MaxRegexMatchMilliseconds int
	// MaxOutputUnits caps accounted diagnostic/annotation units; default 100000, not all intermediate output allocations.
	MaxOutputUnits int
	// MaxCustomKeywordCompiles caps registered keyword compiler calls per Compile; default 100000.
	MaxCustomKeywordCompiles int
	// MaxCustomKeywordCalls caps evaluator calls shared across validation/output traversals; default 100000.
	MaxCustomKeywordCalls int
	// MaxAnnotationBytes caps each retained custom annotation before JSON decoding; default 1 MiB, not built-in annotation totals.
	MaxAnnotationBytes int
}

// DefaultLimits returns conservative standalone defaults.
func DefaultLimits() Limits {
	return Limits{
		MaxInputBytes:             16 << 20,
		MaxNestingDepth:           256,
		MaxTotalValues:            1_000_000,
		MaxObjectMembers:          1_000_000,
		MaxArrayItems:             1_000_000,
		MaxNumberBytes:            4096,
		MaxSchemaResources:        1024,
		MaxTotalSchemaBytes:       64 << 20,
		MaxEvaluationOps:          1_000_000,
		MaxUniqueComparisons:      1_000_000,
		MaxFormatChecks:           100_000,
		MaxSchemaNodes:            100_000,
		MaxReferenceDepth:         256,
		MaxDynamicScopeDepth:      256,
		MaxCombinatorBranches:     100_000,
		MaxRegexCount:             10_000,
		MaxRegexBytes:             1 << 20,
		MaxRegexBacktracking:      100_000,
		MaxRegexMatchMilliseconds: 100,
		MaxOutputUnits:            100_000,
		MaxCustomKeywordCompiles:  100_000,
		MaxCustomKeywordCalls:     100_000,
		MaxAnnotationBytes:        1 << 20,
	}
}

func (limits Limits) validate() error {
	values := []struct {
		name  string
		value int
	}{
		{name: "MaxInputBytes", value: limits.MaxInputBytes},
		{name: "MaxNestingDepth", value: limits.MaxNestingDepth},
		{name: "MaxTotalValues", value: limits.MaxTotalValues},
		{name: "MaxObjectMembers", value: limits.MaxObjectMembers},
		{name: "MaxArrayItems", value: limits.MaxArrayItems},
		{name: "MaxNumberBytes", value: limits.MaxNumberBytes},
		{name: "MaxSchemaResources", value: limits.MaxSchemaResources},
		{name: "MaxTotalSchemaBytes", value: limits.MaxTotalSchemaBytes},
		{name: "MaxEvaluationOps", value: limits.MaxEvaluationOps},
		{name: "MaxUniqueComparisons", value: limits.MaxUniqueComparisons},
		{name: "MaxFormatChecks", value: limits.MaxFormatChecks},
		{name: "MaxSchemaNodes", value: limits.MaxSchemaNodes},
		{name: "MaxReferenceDepth", value: limits.MaxReferenceDepth},
		{name: "MaxDynamicScopeDepth", value: limits.MaxDynamicScopeDepth},
		{name: "MaxCombinatorBranches", value: limits.MaxCombinatorBranches},
		{name: "MaxRegexCount", value: limits.MaxRegexCount},
		{name: "MaxRegexBytes", value: limits.MaxRegexBytes},
		{name: "MaxRegexBacktracking", value: limits.MaxRegexBacktracking},
		{name: "MaxRegexMatchMilliseconds", value: limits.MaxRegexMatchMilliseconds},
		{name: "MaxOutputUnits", value: limits.MaxOutputUnits},
		{name: "MaxCustomKeywordCompiles", value: limits.MaxCustomKeywordCompiles},
		{name: "MaxCustomKeywordCalls", value: limits.MaxCustomKeywordCalls},
		{name: "MaxAnnotationBytes", value: limits.MaxAnnotationBytes},
	}

	for _, item := range values {
		if item.value <= 0 {
			return fmt.Errorf("%w: %s must be positive", ErrLimitExceeded, item.name)
		}
	}

	return nil
}
