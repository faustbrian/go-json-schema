SHELL := /usr/bin/env bash

.PHONY: benchmark conformance docs fuzz

benchmark:
	go test . -run '^$$' -bench Benchmark -benchmem -benchtime=100ms
	cd benchmarks/comparison && go test -run '^$$' -bench . -benchmem -benchtime=100ms

conformance:
	./scripts/check-official-suite.sh
	./scripts/check-official-meta-schemas.sh
	./scripts/check-conformance-manifest.sh
	go test . -run '^TestSpecificationManifestPinsEveryConformanceSource$$' -count=1
	go test . -run '^TestOfficial' -count=1
	go test ./cmd/bowtie-json-schema

docs:
	./scripts/check-docs.sh

fuzz:
	./scripts/check-fuzz.sh 2s
