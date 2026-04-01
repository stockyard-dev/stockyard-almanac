VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo dev)
.PHONY: build run clean
build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION)" -o almanac ./cmd/almanac/
run: build
	./almanac
clean:
	rm -f almanac
