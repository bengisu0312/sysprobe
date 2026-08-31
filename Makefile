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

.PHONY: build run clean fmt vet
