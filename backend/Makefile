.PHONY: help install-tools lint proto schema-dump swag test-functional test-mocks test-unit

help: ## Show this help
	@echo 'usage: make [target]'
	@echo
	@echo 'available targets:'
	@egrep '^.+\:.*##\ .+' ${MAKEFILE_LIST} | sed -r 's/^(.+): .*##\ (.+)/  \1#\2/' | column -t -c 2 -s '#'

install-tools: ## Install the development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/envoyproxy/protoc-gen-validate@v0.9.1
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/swaggo/swag/cmd/swag@v1.8.3

lint: ## Run the linter over all the Go code
	@command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	golangci-lint run

proto: ## Generate Go packages from proto files
	protoc \
		-I ./api/proto \
		--go_out=./pkg/authz --go_opt=paths=source_relative \
		--go-grpc_out=./pkg/authz --go-grpc_opt=paths=source_relative \
		api.proto

schema-dump: ## Dump current database schema into schema.sql file
	docker-compose exec postgres pg_dump -x -s -U root -w -d root > schema.sql

# install needed dependency:
# go install github.com/swaggo/swag/cmd/swag@latest
swag: ## Generate OpenAPI documentation
	swag init \
		-g ./internal/http/routing.go \
		--parseInternal \
		--dir ./ \
		--output ./internal/http/docs/

test-functional: ## Run functional tests
	go test -count=1 --tags=functional -v ./functional --godog.tags=$(tags)

test-mocks: ## Generate unit test mocks
	mockgen -source=internal/helper/time/clock.go -destination=internal/helper/time/clock_mock.go -package=time
	mockgen -source=internal/security/jwt/manager.go -destination=internal/security/jwt/manager_mock.go -package=jwt

test-unit: ## Run unit tests
	go test -count=1 -race -v ./...
