test:
	go test ./internal/...

run:
	go run ./cmd

build:
	go build -o ./bin/wishlist-api ./cmd
