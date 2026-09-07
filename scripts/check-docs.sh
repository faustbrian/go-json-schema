#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "${root}"

snippet_dir="$(mktemp -d "${TMPDIR:-/tmp}/go-json-schema-docs.XXXXXX")"
cleanup() {
	find "${snippet_dir}" -depth -delete
}
trap cleanup EXIT HUP INT TERM

required=(
	README.md CHANGELOG.md CONTRIBUTING.md RELEASING.md LICENSE NOTICE
	SECURITY.md SUPPORT.md COMPATIBILITY.md example_test.go doc.go api/baseline.txt
	docs/adoption.md docs/api.md docs/architecture.md docs/conformance.md
	docs/cookbook.md docs/dependencies.md docs/dialects.md docs/extensions.md
	docs/faq.md docs/limits.md docs/matrices.md docs/output.md
	docs/performance.md docs/quickstart.md docs/resolvers.md docs/security.md
	docs/specification-decisions.md docs/troubleshooting.md docs/versioning.md
	bowtie/README.md
)
for path in "${required[@]}"; do
	if [[ ! -s "${path}" ]]; then
		printf 'required documentation is missing or empty: %s\n' "${path}" >&2
		exit 1
	fi
done

python3 scripts/docs_contract_test.py
python3 scripts/docs_contract.py --root "${root}" --output "${snippet_dir}"

(
	cd "${snippet_dir}"
	GOWORK=off go test -mod=mod ./...
	standalone_count="$(cat standalone-count)"
	for ((index = 0; index < standalone_count; index++)); do
		expectation="standalone-${index}/.expected-output"
		if [[ ! -f "${expectation}" ]]; then
			continue
		fi
		actual="$(GOWORK=off go run -mod=mod "./standalone-${index}")"
		expected="$(cat "${expectation}")"
		if [[ "${actual}" != "${expected}" ]]; then
			printf 'designated documentation example %s returned %q, expected %q\n' \
				"${index}" "${actual}" "${expected}" >&2
			exit 1
		fi
	done
)

for document in README.md docs/quickstart.md; do
	if ! grep -Fq 'go get github.com/faustbrian/go-json-schema@v1' "${document}"; then
		printf 'canonical v1 installation is missing from %s\n' "${document}" >&2
		exit 1
	fi
done

if grep -Fq 'under active development' doc.go api/baseline.txt; then
	printf 'stable-v1 package documentation is stale\n' >&2
	exit 1
fi

while IFS= read -r package; do
	go doc "${package}" >/dev/null
done < <(go list ./...)

example_count="$(grep -Ec '^func Example[[:alnum:]_]*\(\)' example_test.go)"
output_count="$(grep -Ec '^[[:space:]]*// Output:' example_test.go)"
if [[ "${example_count}" -eq 0 || "${example_count}" -ne "${output_count}" ]]; then
	printf 'every executable example must be nonempty and declare output\n' >&2
	exit 1
fi

if grep -En '(, _ :=|_ = .*\.(Close|Compile|Load|Marshal|Validate|ValidateOutput|ValidateValue|ValidateValueOutput)\()' \
	README.md docs/*.md example_test.go; then
	printf 'public documentation must handle returned errors\n' >&2
	exit 1
fi

go test -mod=readonly ./... -run '^Example' -count=1
golib api check
printf 'documentation contract passed\n'
