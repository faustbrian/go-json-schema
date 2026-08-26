# json-schema

[![CI](https://github.com/faustbrian/go-json-schema/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/faustbrian/go-json-schema/actions/workflows/ci.yml)
[![CodeQL](https://img.shields.io/badge/CodeQL-required-blue)](https://github.com/faustbrian/go-json-schema/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Mutation](https://img.shields.io/badge/mutation-100%25_required-blue)](CONTRIBUTING.md#verification)
[![Documentation](https://img.shields.io/badge/docs-checked_in_CI-blue)](docs/)
[![Go Reference](https://pkg.go.dev/badge/github.com/faustbrian/go-json-schema.svg)](https://pkg.go.dev/github.com/faustbrian/go-json-schema)
[![Release](https://img.shields.io/github/v/release/faustbrian/go-json-schema?sort=semver)](https://github.com/faustbrian/go-json-schema/releases)
[![Go](https://img.shields.io/badge/go-1.26.6-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

`json-schema` is an exact-number, dialect-aware JSON Schema compiler and
validator for Go. It supports Draft 3, Draft 4, Draft 6, Draft 7, Draft
2019-09, and Draft 2020-12 without implicit network access or global mutable
registries.

The pinned official suite currently passes 8,505 cases across 354 mandatory
and optional fixture files with zero skips and zero failures. See
[conformance](docs/conformance.md) for the supported surface and evidence
policy.

The minimum supported toolchain is Go 1.26.6.

## Quick start

```go
compiler, err := jsonschema.NewCompiler(
    jsonschema.WithDialect(jsonschema.Draft202012),
)
if err != nil {
    return err
}

schema, err := compiler.Compile(
    context.Background(),
    []byte(`{"type":"object","required":["name"],"properties":{"name":{"type":"string"}}}`),
)
if err != nil {
    return err
}

result, err := schema.Validate(
    context.Background(),
    []byte(`{"name":"Ada"}`),
)
if err != nil {
    return err
}
fmt.Println(result.Valid) // true
```

Compilation validates the schema against the embedded official meta-schema.
Compiled schemas are immutable and reusable concurrently. `Validate` accepts
raw JSON; `ValidateValue` accepts Go values and preserves `json.Number` text.
`ValidateOutput` and `ValidateValueOutput` provide Flag, Basic, Detailed, and
Verbose output units. `CollectAnnotations` returns retained successful-path
annotations as a flat deterministic list.

Format keywords are annotations by default; enable recognized assertions with
`WithFormatAssertion`. Content keywords are annotations by default, and
`WithContentAssertion` enables only Draft 7's optional assertion behavior.
Draft 2019-09 and Draft 2020-12 content processing never changes the enclosing
schema result. Remote references require an explicit `ResourceLoader`; the
core never performs network I/O.

## Documentation

Start with the [documentation index](docs/README.md). It links the quickstart,
API guide, dialect support, conformance contract, secure resolver guidance,
examples, operations, and maintainer references.

## Development

`make check` runs formatting, module, vet, tests, fixture provenance,
conformance-manifest, and Bowtie protocol gates offline after dependencies are
available. `go test -race ./...` is the concurrency gate. See
[CONTRIBUTING.md](CONTRIBUTING.md) for fixture and behavior changes.

## License

MIT. See [LICENSE](LICENSE) and [NOTICE](NOTICE).

## Related packages

- [OpenAPI](https://github.com/faustbrian/go-openapi) uses this package for
  OpenAPI 3.1 and 3.2 Schema Objects.
- [OpenRPC](https://github.com/faustbrian/go-openrpc) provides OpenRPC document
  modeling and validation.
