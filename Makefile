GOLANGCI_LINT ?= golangci-lint
CUSTOM_GOLANGCI_LINT := ./bin/golangci-lint-passthrough

.PHONY: build lint test test-integration fmt-check ci clean

build:
	$(GOLANGCI_LINT) custom

lint: build
	$(CUSTOM_GOLANGCI_LINT) run ./...

test:
	go test ./...

test-integration: build
	bash scripts/integration-test.sh $(CUSTOM_GOLANGCI_LINT)

fmt-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './bin/*'))" || \
		(gofmt -d $$(find . -name '*.go' -not -path './bin/*'); exit 1)

ci: fmt-check test test-integration lint

clean:
	rm -rf bin
