GO ?= go
BENCH_TIME ?= 100ms

.PHONY: benchmark conformance dependencies docs interoperability

benchmark:
	BENCH_TIME="$(BENCH_TIME)" ./scripts/check-benchmarks.sh

conformance:
	./scripts/check-conformance.sh

dependencies:
	./scripts/check-dependencies.sh

docs:
	./scripts/check-docs.sh

interoperability:
	GOWORK=off $(GO) test ./... -run \
		'^(TestDecodeIndependentWriters|TestSymbolsDecodeWithIndependentReaders)$$' \
		-count=1
