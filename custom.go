package jsonschema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
)

// ValueKind identifies an immutable exact JSON value kind.
type ValueKind uint8

const (
	// NullKind identifies null or an absent Value.
	NullKind ValueKind = iota
	// BooleanKind identifies a JSON boolean.
	BooleanKind
	// NumberKind identifies an exact JSON number.
	NumberKind
	// StringKind identifies a JSON string.
	StringKind
	// ArrayKind identifies a JSON array.
	ArrayKind
	// ObjectKind identifies a JSON object.
	ObjectKind
)

// Value is a read-only, retainable view of decoded schema or instance JSON.
// A zero Value reports NullKind but is absent from Index/Lookup failures.
// Names returns a detached sorted slice; child views share immutable storage.
// Number preserves lexical decimal text; strings use encoding/json escape
// decoding, including its replacement of unpaired surrogate escapes.
type Value struct {
	value *jsonValue
}

// Kind returns the JSON value kind.
func (value Value) Kind() ValueKind {
	if value.value == nil {
		return NullKind
	}
	return ValueKind(value.value.kind)
}

// Bool returns a boolean value and whether the kind matched.
func (value Value) Bool() (bool, bool) {
	if value.value == nil || value.value.kind != kindBoolean {
		return false, false
	}
	return value.value.boolean, true
}

// Number returns the exact JSON number text and whether the kind matched.
func (value Value) Number() (string, bool) {
	if value.value == nil || value.value.kind != kindNumber {
		return "", false
	}
	return value.value.number, true
}

// String returns a string value and whether the kind matched.
func (value Value) String() (string, bool) {
	if value.value == nil || value.value.kind != kindString {
		return "", false
	}
	return value.value.text, true
}

// Len returns the array item or object member count, or zero otherwise.
func (value Value) Len() int {
	if value.value == nil {
		return 0
	}
	switch value.value.kind {
	case kindArray:
		return len(value.value.array)
	case kindObject:
		return len(value.value.object)
	default:
		return 0
	}
}

// Index returns one array item.
func (value Value) Index(index int) (Value, bool) {
	if value.value == nil || value.value.kind != kindArray ||
		index < 0 || index >= len(value.value.array) {
		return Value{}, false
	}
	return Value{value: value.value.array[index]}, true
}

// Lookup returns one object member.
func (value Value) Lookup(name string) (Value, bool) {
	if value.value == nil || value.value.kind != kindObject {
		return Value{}, false
	}
	item, exists := value.value.object[name]
	return Value{value: item}, exists
}

// Names returns a sorted copy of object member names.
func (value Value) Names() []string {
	if value.value == nil || value.value.kind != kindObject {
		return nil
	}
	return sortedStringKeys(value.value.object)
}

// KeywordCompiler compiles one keyword value into an immutable evaluator.
type KeywordCompiler interface {
	// Compile borrows a read-only exact keyword Value and returns a nonnil
	// evaluator whose state supports concurrent reuse. Calls may overlap across
	// compilations; return promptly under context cancellation. The registry
	// does not deep-copy callbacks or roll back their side effects.
	Compile(context.Context, Dialect, Value) (KeywordEvaluator, error)
}

// KeywordCompilerFunc adapts a function to KeywordCompiler.
type KeywordCompilerFunc func(context.Context, Dialect, Value) (KeywordEvaluator, error)

// Compile directly invokes the function. Nil functions/direct callback panics
// propagate here; Compiler supplies containment when compiling registered
// keywords, not on arbitrary calls to this adapter.
func (compiler KeywordCompilerFunc) Compile(
	ctx context.Context,
	dialect Dialect,
	value Value,
) (KeywordEvaluator, error) {
	return compiler(ctx, dialect, value)
}

// KeywordEvaluator evaluates a compiled custom keyword.
type KeywordEvaluator interface {
	// Evaluate borrows a read-only exact instance view. It may be called more
	// than once for validation, tracking and output and concurrently across
	// requests; results must be deterministic for an input and cancellation
	// cooperative. Error is operational failure, not invalidity; callback
	// panics are contained only through Schema entry points, not this interface.
	Evaluate(context.Context, Value) (KeywordResult, error)
}

// KeywordEvaluatorFunc adapts a function to KeywordEvaluator.
type KeywordEvaluatorFunc func(context.Context, Value) (KeywordResult, error)

// Evaluate directly invokes the function without result copying or panic
// containment. Schema entry points supply that boundary; nil functions panic
// when called directly.
func (evaluator KeywordEvaluatorFunc) Evaluate(
	ctx context.Context,
	value Value,
) (KeywordResult, error) {
	return evaluator(ctx, value)
}

// KeywordResult reports custom assertion validity and an optional exact JSON
// annotation. Nil Annotation means no annotation; explicit JSON null creates
// an annotation unit, but OutputUnit currently decodes it to nil and omits
// the annotation field during JSON marshaling. This is a wire limitation,
// not a different specification interpretation.
type KeywordResult struct {
	// Valid is the assertion result when Evaluate returns no error.
	Valid bool
	// Annotation is borrowed JSON for retained output; keep bytes unchanged while consumed. Nil omits it; explicit JSON null decodes to nil and its output field is omitted. Flag validation ignores it.
	Annotation json.RawMessage
}

type registeredVocabulary struct {
	keywords map[string]KeywordCompiler
}

// WithVocabulary registers an owned map of borrowed keyword compilers when
// the option runs. Keep the supplied map unchanged through NewCompiler. An
// absolute fragment-free URI, unique nonempty non-built-in keyword names and
// nonnil compilers are required. Modern activation follows the resource
// vocabulary declaration; historical drafts activate registered keywords.
// See JSONSCHEMA-DEC-001/002; registration is not external IO authorization.
func WithVocabulary(
	identifier string,
	keywords map[string]KeywordCompiler,
) Option {
	return func(config *compilerConfig) error {
		parsed, err := url.Parse(identifier)
		if err != nil || !parsed.IsAbs() || parsed.Fragment != "" {
			return fmt.Errorf("%w: invalid vocabulary identifier %q", ErrInvalidSchema, identifier)
		}
		if _, exists := config.vocabularies[identifier]; exists {
			return fmt.Errorf("%w: duplicate vocabulary %q", ErrInvalidSchema, identifier)
		}
		registered := registeredVocabulary{keywords: make(map[string]KeywordCompiler, len(keywords))}
		for name, compiler := range keywords {
			if name == "" || isKnownKeyword(name) || interfaceIsNil(compiler) {
				return fmt.Errorf("%w: invalid custom keyword %q", ErrInvalidSchema, name)
			}
			for _, vocabulary := range config.vocabularies {
				if _, duplicate := vocabulary.keywords[name]; duplicate {
					return fmt.Errorf("%w: duplicate custom keyword %q", ErrInvalidSchema, name)
				}
			}
			registered.keywords[name] = compiler
		}
		config.vocabularies[identifier] = registered
		return nil
	}
}

func interfaceIsNil(value any) bool {
	if value == nil {
		return true
	}
	reflected := reflect.ValueOf(value)
	switch reflected.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice:
		return reflected.IsNil()
	default:
		return false
	}
}

func cloneVocabularies(
	source map[string]registeredVocabulary,
) map[string]registeredVocabulary {
	result := make(map[string]registeredVocabulary, len(source))
	for identifier, vocabulary := range source {
		keywords := make(map[string]KeywordCompiler, len(vocabulary.keywords))
		for name, compiler := range vocabulary.keywords {
			keywords[name] = compiler
		}
		result[identifier] = registeredVocabulary{keywords: keywords}
	}
	return result
}
