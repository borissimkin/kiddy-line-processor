help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

run:
	@if [ ! -f .env ]; then \
		echo ".env not found. Use .env.example..."; \
		cp .env.example .env; \
	fi
	docker-compose up --build

tests:
	go test ./...

lint:
	go vet ./...

stop:
	docker-compose stop

up-env:
	docker-compose up -d lines-provider redis

run-app:
	go run ./cmd/main.go

protoc:
	@echo "Generating Go files"
	cd internal/proto && protoc --go_out=. --go-grpc_out=. \
		--go-grpc_opt=paths=source_relative --go_opt=paths=source_relative *.proto