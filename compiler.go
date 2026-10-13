package jsonschema

import (
	"context"
	"fmt"
	"reflect"
)

// Option configures a Compiler without mutating shared global state.
type Option func(*compilerConfig) error

type compilerConfig struct {
	dialect       Dialect
	limits        Limits
	loader        ResourceLoader
	assertFormats bool
	assertContent bool
	formats       map[string]FormatChecker
	vocabularies  map[string]registeredVocabulary
}

// ResourceLoader retrieves schema resources under application authorization
// and IO policy. The compiler borrows the loader without closing it.
type ResourceLoader interface {
	// Load receives a normalized resource URI without a reference fragment.
	// Return bytes stable until the consuming Compile call finishes; loading may
	// overlap across compilations and must cooperate with a nonnil context.
	// The compiler checks byte budgets only after Load returns. Classify misses
	// as ErrResourceNotFound for composite fallback; other errors stop it.
	// Package callback wrappers contain panics and hide error text, but preserve
	// raw error causes and cannot interrupt arbitrary blocking.
	Load(context.Context, string) ([]byte, error)
}

// ResourceLoaderFunc adapts a function to ResourceLoader.
type ResourceLoaderFunc func(context.Context, string) ([]byte, error)

// Load directly invokes the function without copying bytes or containing
// panic. A nil function panics on direct use; Compiler callback containment
// applies only when the adapter is invoked through Compile.
func (loader ResourceLoaderFunc) Load(ctx context.Context, identifier string) ([]byte, error) {
	return loader(ctx, identifier)
}

// WithDialect selects the dialect used to compile schemas.
func WithDialect(dialect Dialect) Option {
	return func(config *compilerConfig) error {
		if err := dialect.validate(); err != nil {
			return err
		}
		config.dialect = dialect

		return nil
	}
}

// WithLimits replaces the entire limits value; every field must be positive.
// Start from DefaultLimits rather than a partially populated zero value.
func WithLimits(limits Limits) Option {
	return func(config *compilerConfig) error {
		if err := limits.validate(); err != nil {
			return err
		}
		config.limits = limits

		return nil
	}
}

// WithResourceLoader installs a borrowed loader; nil and typed-nil values
// reject. Bundled official meta-schemas resolve locally before this loader.
// Registration is not scheme/host authorization or a network byte/time limit.
func WithResourceLoader(loader ResourceLoader) Option {
	return func(config *compilerConfig) error {
		if resourceLoaderIsNil(loader) {
			return fmt.Errorf("%w: nil resource loader", ErrResourceUnavailable)
		}
		config.loader = loader

		return nil
	}
}

// WithFormatAssertion enables format validation for recognized formats.
func WithFormatAssertion() Option {
	return func(config *compilerConfig) error {
		config.assertFormats = true

		return nil
	}
}

// WithContentAssertion enables Draft 7 validation of recognized content
// encodings and media types. Later dialects keep content as annotations.
func WithContentAssertion() Option {
	return func(config *compilerConfig) error {
		config.assertContent = true

		return nil
	}
}

// WithFormat registers or replaces a borrowed format checker in an owned
// registry. Nil/typed-nil checkers and empty names reject. Registries are copied,
// not checker state; callers own concurrency, cancellation and lifetime.
func WithFormat(name string, checker FormatChecker) Option {
	return func(config *compilerConfig) error {
		if name == "" || formatCheckerIsNil(checker) {
			return fmt.Errorf("%w: invalid format registration %q", ErrInvalidSchema, name)
		}
		config.formats[name] = customFormatChecker{checker: checker}

		return nil
	}
}

func resourceLoaderIsNil(loader ResourceLoader) bool {
	if loader == nil {
		return true
	}
	value := reflect.ValueOf(loader)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func formatCheckerIsNil(checker FormatChecker) bool {
	if checker == nil {
		return true
	}
	value := reflect.ValueOf(checker)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

// Compiler owns immutable registry maps and resource policy, not the state of
// borrowed loader/format/keyword callbacks. Concurrent Compile calls create
// independent caches but require those callbacks to support concurrent use.
// A zero Compiler is not a configured compiler; construct it with NewCompiler.
type Compiler struct {
	dialect       Dialect
	limits        Limits
	loader        ResourceLoader
	assertFormats bool
	assertContent bool
	formats       map[string]FormatChecker
	metaSchema    *Schema
	vocabularies  map[string]registeredVocabulary
}

// NewCompiler constructs an isolated compiler. Draft 2020-12 is the explicit
// default when no dialect option is supplied. Options run synchronously, in
// order; nil options reject and option panics are not contained. Construction
// prepares the bundled meta-schema with default budgets/background context,
// without application IO. Application limits apply to later compile/evaluation
// work rather than constructor meta-schema preparation.
func NewCompiler(options ...Option) (*Compiler, error) {
	config := compilerConfig{
		dialect:      Draft202012,
		limits:       DefaultLimits(),
		formats:      standardFormats(),
		vocabularies: make(map[string]registeredVocabulary),
	}

	for index, option := range options {
		if option == nil {
			return nil, fmt.Errorf("compiler option %d is nil", index)
		}
		if err := option(&config); err != nil {
			return nil, fmt.Errorf("compiler option %d: %w", index, err)
		}
	}
	return newCompiler(config, compileOfficialMetaSchema)
}

func newCompiler(
	config compilerConfig,
	prepareMetaSchema func(Dialect) (*Schema, error),
) (*Compiler, error) {
	formats := cloneFormats(config.formats)
	applyStandardFormatLimits(formats, config.limits)
	compiler := &Compiler{
		dialect:       config.dialect,
		limits:        config.limits,
		loader:        config.loader,
		assertFormats: config.assertFormats,
		assertContent: config.assertContent,
		formats:       formats,
		vocabularies:  cloneVocabularies(config.vocabularies),
	}
	metaSchema, err := prepareMetaSchema(config.dialect)
	if err != nil {
		return nil, fmt.Errorf("prepare %s meta-schema: %w", config.dialect, err)
	}
	compiler.metaSchema = metaSchema

	return compiler, nil
}

// Schema holds an immutable plan and per-call evaluation state. Concurrent
// validation also requires borrowed format checkers and custom evaluators to
// be safe for concurrent use. It owns no Close/shutdown or worker lifecycle.
// A zero Schema rejects validation with ErrInvalidSchema.
type Schema struct {
	dialect Dialect
	limits  Limits
	plan    *schemaPlan
}

// Compile parses one document and resolves reachable references into a
// per-call cache. Input bytes are borrowed during parsing, not retained by the
// plan, and must not change during the call. Nil context rejects; cancellation
// is cooperative, not preemption of callbacks or every compilation step.
// Configured dialect is the root fallback; recognized embedded dialects and
// modern vocabulary policy are resource-scoped (JSONSCHEMA-DEC-001/008).
// Schema meta-validation uses separate default budgets. Operational/budget
// failures return no plan; callback side effects are not rolled back.
// URI diagnostics redact user information, query and fragment components.
// Parse/schema errors can still contain schema names;
// inspect typed classifications and keep untrusted diagnostics out of logs.
func (compiler *Compiler) Compile(ctx context.Context, raw []byte) (*Schema, error) {
	if compiler == nil {
		return nil, fmt.Errorf("%w: nil compiler", ErrInvalidSchema)
	}
	if len(raw) > compiler.limits.MaxTotalSchemaBytes {
		return nil, &LimitError{
			Resource: "total schema bytes",
			Limit:    compiler.limits.MaxTotalSchemaBytes,
		}
	}

	value, err := decodeJSON(ctx, raw, compiler.limits)
	if err != nil {
		return nil, fmt.Errorf("compile schema: %w", err)
	}
	if compiler.metaSchema != nil {
		state := schemaEvaluationState(
			ctx,
			compiler.metaSchema.plan,
			DefaultLimits(),
		)
		valid, err := compiler.metaSchema.plan.evaluate(
			value,
			compiler.dialect,
			&state,
		)
		if err != nil {
			return nil, fmt.Errorf("validate schema: %w", err)
		}
		if !valid {
			return nil, fmt.Errorf(
				"%w: schema does not satisfy the %s meta-schema",
				ErrInvalidSchema,
				compiler.dialect,
			)
		}
	}

	schemaCompiler := newSchemaCompiler(
		ctx,
		value,
		compiler.dialect,
		compiler.limits,
		compiler.loader,
		len(raw),
		compiler.assertFormats,
		compiler.assertContent,
		compiler.formats,
		compiler.vocabularies,
	)
	if err := schemaCompiler.configureVocabularies(value); err != nil {
		return nil, err
	}
	plan, err := schemaCompiler.compile(value)
	if err != nil {
		return nil, err
	}

	return &Schema{
		dialect: compiler.dialect,
		limits:  compiler.limits,
		plan:    plan,
	}, nil
}

func schemaEvaluationState(
	ctx context.Context,
	plan *schemaPlan,
	limits Limits,
) evaluationState {
	state := evaluationState{ctx: ctx, limits: limits}
	if plan.resource != nil {
		state.dynamicScope = append(state.dynamicScope, plan.resource)
	}
	return state
}
