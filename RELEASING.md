# Releasing

The root module is released from a verified commit with a `v<version>` tag.
Only an independently releasable nested module uses a
`<module-directory>/v<version>` tag. The internal `benchmarks/comparison`
module is not released independently.

Before creating a release, require all of these to be true:

- every mandatory and optional released-dialect official case passes with
  zero failures and zero unexplained skips;
- every Core and Validation requirement maps to implementation, tests, and
  documentation;
- official meta-validation, references, dynamic scope, vocabularies,
  annotations, formats, content, and all output forms have no known gap;
- Bowtie runs all six dialects reproducibly and published reports are current;
- meaningful production statement coverage is 100%;
- race, fuzz, mutation, vulnerability, API, documentation, benchmark,
  dependency, license, secret, workflow-security, and owned-analysis gates
  pass;
- security limits and threat behavior are tested and documented;
- every public API has documentation and executable examples;
- changelog, compatibility notes, and benchmark evidence are current.

For a release candidate, run every local gate from a clean checkout, compare
the generated conformance manifest, review dependencies and licenses, build
the Bowtie image, inspect the full diff from the prior tag, and obtain review.
Tag the verified root commit with the root semantic version.
Release notes must distinguish normative behavior, implementation policy,
optional capabilities, convenience APIs, and any remaining non-v1 limitation.
