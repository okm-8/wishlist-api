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

build-image:
	docker build -t wishlist-api:latest .

deploy-helm:
	helm upgrade --install --wait --namespace wishlist --create-namespace wishlist-api ./charts

delete-helm:
	helm delete --namespace wishlist wishlist-api
