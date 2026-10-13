# json-schema documentation map

[Repository README](../README.md) · [Current Golib portal](https://github.com/faustbrian/go-library-tools/blob/main/docs/ecosystem/README.md) · [Published ecosystem snapshot](https://github.com/faustbrian/go-library-tools/blob/v1.8.4/docs/ecosystem/README.md)

This map links package-owned adoption and operational guidance. Start with
the repository README for purpose, installation, and ownership boundaries.
Published compatibility sets are frozen multi-package baselines; they do not
certify every current main revision or a deployed application.

## Module boundaries

Import only the modules needed by your application. A deprecated facade,
preview adapter, or internal harness is not interchangeable with a stable
root API; follow its own README and migration guidance before composing it.

| Module | Status | Owner detail |
| --- | --- | --- |
| `github.com/faustbrian/go-json-schema/v2` | stable | [README](../README.md) |
| `github.com/faustbrian/go-json-schema/benchmarks/comparison` | internal; non-releasable | [README](../benchmarks/comparison/README.md) |

## Choose and adopt

- [Documentation](README.md).
- [Adoption](adoption.md).
- [API guide](api.md).
- [Architecture](architecture.md).
- [Cookbook](cookbook.md).
- [Dependency review](dependencies.md).
- [Dialect selection and migration](dialects.md).
- [Custom vocabularies, keywords, and formats](extensions.md).
- [FAQ and troubleshooting](faq.md).
- [Resource limits](limits.md).
- [Dialect and keyword matrices](matrices.md).
- [Validation output](output.md).
- [Quick start](quickstart.md).
- [Resolvers and secure loading](resolvers.md).
- [Compatibility, deprecation, and versioning](versioning.md).

## Operate and evaluate boundaries

- [Security and threat model](security.md).
- [Troubleshooting](troubleshooting.md).

## Owned adapters, examples, and harnesses

- [Bowtie interoperability](../bowtie/README.md).

## Maintainer, specification, and evidence references

- [JSON Schema comparison benchmarks](../benchmarks/comparison/README.md).
- [Conformance](conformance.md).
- [Performance methodology](performance.md).
- [Source documentation audit](source-documentation-audit.md).
- [JSON Schema specification decisions](specification-decisions.md).
- [Normative and conformance inputs](../specification/README.md).

## Related next steps

[Contributor setup and checks](../CONTRIBUTING.md), [Compatibility policy](../COMPATIBILITY.md), [Release history](../CHANGELOG.md), [Private vulnerability reporting](../SECURITY.md), [Support](../SUPPORT.md).
