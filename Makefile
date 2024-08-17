.PHONY: clean clean-test test build

clean:
	rm -rf ./bin

clean-test:
	go clean -testcache

test:
	go test ./internal/...

build:
	go build -o ./bin/wishlist-api ./cmd

run-system:
	go run ./cmd system core
