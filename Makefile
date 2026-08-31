BINARY := sysprobe

build:
	go build -o bin/$(BINARY) ./cmd/sysprobe

run: build
	./bin/sysprobe status

clean:
	rm -rf bin/

fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

test-race:
	go test -race ./...

coverage:
	go test -coverprofile=coverage.out ./internal/collector
	go tool cover -func=coverage.out

.PHONY: build run clean fmt vet test test-race coverage
