BINARY := sysprobe

build:
	go build -o bin/$(BINARY) ./cmd/sysprobe

run: build
	./bin/$(BINARY) status

clean:
	rm -rf bin/

.PHONY: build run clean
