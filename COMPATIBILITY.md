# Compatibility Policy

Each releasable directory is an independent Go module and follows semantic
versioning. The root module uses `v<version>` tags. An independently releasable
nested module uses `<module-directory>/v<version>` tags; the
`benchmarks/comparison` module is internal and is not released independently.

Before `v1`, minor releases MAY contain reviewed breaking changes, but every
break MUST be documented with migration guidance. Patch releases MUST remain
backward compatible. At and after `v1`, incompatible exported API or documented
behavior changes require a new major version.

The latest published patch release in the current stable major is the supported
version for defect and security fixes. Earlier releases remain available but
are not maintained after a newer patch is published. The project does not
promise a calendar-based long-term-support window; consumers should upgrade to
the latest stable patch before requesting support.

Compatibility includes exported Go APIs, error classification, serialization,
protocol behavior, persistence schemas, environment variables, command output,
resource ownership, ordering, retry/idempotency semantics, and documented
defaults. A compile-compatible change can still be behaviorally breaking.

Specification-backed modules MUST NOT diverge from their declared standards.
Ambiguities require documented decisions and stable tests. Deprecated APIs
follow [`DEPRECATION.md`](DEPRECATION.md).
