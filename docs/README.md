# Documentation

`json-schema` compiles and validates JSON Schema documents across Draft 3
through Draft 2020-12. Start with the quickstart, then use the sections below
for the contract you need.

The public `jsonschema` package is the application entry point. The
`cmd/bowtie-json-schema` command is the interoperability adapter,
`internal/cmd/conformance-manifest` is a maintainer-only evidence generator,
and `benchmarks/comparison` is an internal, unreleased nested module. The
caller owns configuration and resource loading; compiled schemas are immutable
and require no shutdown.

There is no dedicated testing-helper package. Application tests can compose
deterministic seams with `MapLoader`, `ResourceLoaderFunc`, `FormatFunc`, and
explicit `KeywordCompiler` callbacks registered through `WithVocabulary`; the
[API guide](api.md#extension-points) documents their ownership.

## Getting started

- [Quickstart](quickstart.md)
- [API guide](api.md)
- [When to use this package](adoption.md)
- [Cookbook](cookbook.md)

## Specifications and behavior

- [Dialect selection and migration](dialects.md)
- [Dialect and keyword matrices](matrices.md)
- [Conformance](conformance.md)
- [Specification decisions](specification-decisions.md)
- [Validation output](output.md)
- [Custom vocabularies, keywords, and formats](extensions.md)

## Architecture and operations

- [Architecture](architecture.md)
- [Resolvers and secure loading](resolvers.md)
- [Resource limits](limits.md)
- [Security](security.md)
- [Performance](performance.md)
- [Troubleshooting](troubleshooting.md)

## Reference and maintenance

- [Compatibility and versioning](versioning.md)
- [Dependencies](dependencies.md)
- [FAQ](faq.md)
- [Bowtie interoperability](../bowtie/README.md)
- [Contributing](../CONTRIBUTING.md)
- [Releasing](../RELEASING.md)
- [Compatibility](../COMPATIBILITY.md)
- [Support](../SUPPORT.md)
- [Security reporting](../SECURITY.md)
- [Changelog](../CHANGELOG.md)
- [License](../LICENSE)
