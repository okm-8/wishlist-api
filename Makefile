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

docker-build:
	docker build -t wishlist-api:latest .

helm-install:
	helm upgrade --install --wait --namespace wishlist --create-namespace wishlist-api ./charts

helm-delete:
	helm delete --namespace wishlist wishlist-api

kube-forward-private:
	kubectl port-forward --namespace wishlist service/wishlist-api-private 8081:8081 

kube-forward-public:
	kubectl port-forward --namespace wishlist service/wishlist-api-public 8080:8080

kube-forward-postgres:
	kubectl port-forward --namespace wishlist service/postgres 8001:5432

kube-forward-redis:
	kubectl port-forward --namespace wishlist service/redis 8002:6379
