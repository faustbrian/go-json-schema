package jsonschema

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidJSON classifies malformed or ambiguous JSON input.
	ErrInvalidJSON = errors.New("invalid JSON")
	// ErrInvalidSchema classifies schemas that are invalid for their dialect.
	ErrInvalidSchema = errors.New("invalid schema")
	// ErrLimitExceeded classifies work rejected by an explicit resource limit.
	ErrLimitExceeded = errors.New("resource limit exceeded")
	// ErrResourceUnavailable classifies a schema resource that could not be loaded.
	ErrResourceUnavailable = errors.New("schema resource unavailable")
	// ErrResourceNotFound classifies an identifier absent from a loader.
	ErrResourceNotFound = errors.New("schema resource not found")
	// ErrUnsupportedDialect classifies an unknown stable dialect.
	ErrUnsupportedDialect = errors.New("unsupported dialect")
	// ErrUnsupportedVocabulary classifies an unknown required vocabulary.
	ErrUnsupportedVocabulary = errors.New("unsupported vocabulary")
	// ErrCallbackPanic classifies a recovered application callback panic.
	ErrCallbackPanic = errors.New("callback panic")
)

// JSONError describes a JSON ingestion failure without retaining input bytes.
type JSONError struct {
	// Offset is the decoder byte offset, or zero for failures before decoder creation.
	Offset int64
	// Kind is the broad invalid-JSON or limit classification used by errors.Is.
	Kind error
	// Cause retains the diagnostic cause; it may include member names or decoder text and is not a secrecy guarantee.
	Cause error
}

// LimitError reports which deterministic work budget was exhausted.
type LimitError struct {
	// Resource names the exhausted budget, not the offending input value.
	Resource string
	// Limit is the configured bound for that resource.
	Limit int
}

// Error implements error.
func (err *LimitError) Error() string {
	return fmt.Sprintf("%v: %s exceeds %d", ErrLimitExceeded, err.Resource, err.Limit)
}

// Unwrap classifies the error as ErrLimitExceeded.
func (err *LimitError) Unwrap() error {
	return ErrLimitExceeded
}

// Error implements error.
func (err *JSONError) Error() string {
	if err.Offset > 0 {
		return fmt.Sprintf("%v at byte %d: %v", err.Kind, err.Offset, err.Cause)
	}

	return fmt.Sprintf("%v: %v", err.Kind, err.Cause)
}

// Unwrap exposes both the classification and underlying cause.
func (err *JSONError) Unwrap() []error {
	return []error{err.Kind, err.Cause}
}
