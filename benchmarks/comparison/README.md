# JSON Schema comparison benchmarks

This non-releasable benchmark and interoperability module isolates maintained
peer dependencies from the public JSON Schema module. It compares Draft
2020-12 compilation and validation with `kaptinlin/jsonschema` and
`santhosh-tekuri/jsonschema/v6`, and verifies the deliberately different
outcomes recorded for selected specification decisions.

The implementations do not expose identical compilation contracts. The local
compiler validates schemas against the official meta-schema and enforces
explicit context, resource, exact-number, output, and callback bounds. Results
are comparable only when those semantic differences remain visible.

Run the differential check and repeated benchmarks from this directory:

```sh
go test -run '^TestMaintainedPeerDecisionOutcomes$' -count=1
go test -run '^$' -bench . -benchmem -count=10
```

Record the Go version, dependency versions, GOOS/GOARCH, CPU, benchmark count,
limits, and raw output. The checked-in baseline and interpretation guidance
live in the parent module's
[performance guide](../../docs/performance.md). Shared construction,
ownership, lifecycle, and composition expectations are in the versioned
[Golib ecosystem index](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/README.md).
