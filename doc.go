// Package jsonschema compiles and evaluates JSON Schemas without implicit
// network access.
//
// Dialect selection, schema retrieval, extension registries, format policy,
// output shape, and resource limits are explicit compiler configuration. A
// compiled plan is immutable; concurrent calls additionally require borrowed
// loaders, custom format checkers and keyword evaluators to support concurrency.
// NewCompiler and Compile are normal entry points, followed by Schema.Validate.
// Supported dialects are Drafts 3, 4, 6, 7, 2019-09 and 2020-12; unknown dialects
// and required unknown vocabularies fail. Optional retrieval never adds implicit
// networking, and this module owns no callback/provider shutdown.
//
// The module provides a stable v2 API. Compliance claims are made only by the
// generated conformance evidence committed with the module.
package jsonschema
