package jsonschema

import (
	"errors"
	"time"

	"github.com/dlclark/regexp2/v2"
)

type ecmaPattern struct {
	compiled          *regexp2.Regexp
	backtrackingLimit int
	matchMilliseconds int
}

func compilePattern(pattern string) (*ecmaPattern, error) {
	return compilePatternWithLimits(pattern, DefaultLimits())
}

func compilePatternWithLimits(pattern string, limits Limits) (*ecmaPattern, error) {
	compiled, err := regexp2.Compile(
		pattern,
		regexp2.ECMAScript|regexp2.Unicode,
		regexp2.OptionMaxBacktrackingStackSize(limits.MaxRegexBacktracking),
	)
	if err != nil {
		return nil, err
	}
	compiled.MatchTimeout = time.Duration(limits.MaxRegexMatchMilliseconds) * time.Millisecond
	return &ecmaPattern{
		compiled:          compiled,
		backtrackingLimit: limits.MaxRegexBacktracking,
		matchMilliseconds: limits.MaxRegexMatchMilliseconds,
	}, nil
}

// JSONSCHEMA-DEC-011 binds ECMA-262 16th edition to the pinned regexp2
// ECMAScript+Unicode engine. Matching has
// engine stack/time guards, not a context-interruptible call. Replace its
// input-bearing errors with typed limits; regex_internal_test.go guards this.
func (pattern *ecmaPattern) matchString(value string) (bool, error) {
	matched, err := pattern.compiled.MatchString(value)
	if err == nil {
		return matched, nil
	}
	if errors.Is(err, regexp2.ErrBacktrackingStackLimit) {
		return false, &LimitError{
			Resource: "regular expression backtracking",
			Limit:    pattern.backtrackingLimit,
		}
	}
	return false, &LimitError{
		Resource: "regular expression match milliseconds",
		Limit:    pattern.matchMilliseconds,
	}
}
