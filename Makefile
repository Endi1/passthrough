GOLANGCI_LINT ?= golangci-lint
BINARY := ./bin/passthrough

.PHONY: build lint test vet fmt-check ci clean

build:
	go build -o $(BINARY) ./cmd/passthrough

lint:
	$(GOLANGCI_LINT) run ./...

test:
	go test ./...

vet:
	go vet ./...

fmt-check:
	@test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './bin/*'))" || \
		(gofmt -d $$(find . -name '*.go' -not -path './bin/*'); exit 1)

ci: fmt-check test vet build lint

clean:
	rm -rf bin
