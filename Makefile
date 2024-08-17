.PHONY: test build

test:
	go test ./internal/...

build:
	go build -o ./bin/wishlist-api ./cmd
