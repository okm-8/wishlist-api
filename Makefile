.PHONY: clean clean-test test build

clean:
	rm -rf ./bin

clean-test:
	go clean -testcache

path ?= ./internal/...
test:
	go test $(path)

build:
	go build -o ./bin/wishlist-api ./cmd

run-system:
	go run ./cmd system core

run-cli:
	go run ./cmd $(cmd)
