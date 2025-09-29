.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: run
run:
	@if [ ! -f .env ]; then \
		echo ".env not found. Use .env.example..."; \
		cp .env.example .env; \
	fi
	docker-compose up --build

.PHONY: tests
tests:
	go test -v ./... -race

.PHONY: lint
lint:
	go vet -v ./...

.PHONY: stop
stop:
	docker-compose stop

.PHONY: up-env
up-env:
	docker-compose up -d lines-provider redis

.PHONY: run-app
run-app:
	go run ./cmd/main.go

.PHONY: protoc
protoc:
	@echo "Generating grpc"
	protoc -I api \
  	--go_out=internal/proto --go_opt=paths=source_relative \
  	--go-grpc_out=internal/proto --go-grpc_opt=paths=source_relative \
  	api/*.proto

.PHONY: coverage-html
cover-html: ### run test with coverage and open html report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: coverage
cover: ### run test with coverage
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

mockgen:
	mockgen -source=internal/ready/service.go -destination=internal/ready/mocks/service.go -package=readymocks