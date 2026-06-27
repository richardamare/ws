.PHONY: build test fmt vet check install snapshot

build:
	go build -o bin/ws ./cmd/ws

test:
	go test -race -count=1 ./...

fmt:
	gofmt -w cmd internal

vet:
	go vet ./...

# what CI runs
check: vet test
	@test -z "$$(gofmt -l cmd internal)" || { echo "run 'make fmt'"; exit 1; }

install:
	go install ./cmd/ws

# local release dry-run (requires goreleaser); does not publish
snapshot:
	goreleaser release --snapshot --clean
