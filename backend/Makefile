.PHONY: help install-tools lint proto mysql-schema-dump postgres-schema-dump swag test-functional test-mocks test-unit

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

mysql-schema-dump: ## Dump current database schema into schema.sql file
	docker-compose -f ../docker-compose.mysql.yaml exec mysql \
		mysqldump -u root -ptoor --no-data root > schema.sql

postgres-schema-dump: ## Dump current database schema into schema.sql file
	docker-compose -f ../docker-compose.postgres.yaml exec postgres \
		pg_dump -x -s -U root -w -d root > schema.sql

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
	mockgen -source=internal/compile/compiler.go -destination=internal/compile/compiler_mock.go -package=compile
	mockgen -source=internal/event/dispatcher.go -destination=internal/event/dispatcher_mock.go -package=event
	mockgen -source=internal/entity/manager/action.go -destination=internal/entity/manager/action_mock.go -package=manager
	mockgen -source=internal/entity/manager/attribute.go -destination=internal/entity/manager/attribute_mock.go -package=manager
	mockgen -source=internal/entity/manager/client.go -destination=internal/entity/manager/client_mock.go -package=manager
	mockgen -source=internal/entity/manager/compiled.go -destination=internal/entity/manager/compiled_mock.go -package=manager
	mockgen -source=internal/entity/manager/policy.go -destination=internal/entity/manager/policy_mock.go -package=manager
	mockgen -source=internal/entity/manager/principal.go -destination=internal/entity/manager/principal_mock.go -package=manager
	mockgen -source=internal/entity/manager/resource.go -destination=internal/entity/manager/resource_mock.go -package=manager
	mockgen -source=internal/entity/manager/role.go -destination=internal/entity/manager/role_mock.go -package=manager
	mockgen -source=internal/entity/manager/stats.go -destination=internal/entity/manager/stats_mock.go -package=manager
	mockgen -source=internal/entity/manager/user.go -destination=internal/entity/manager/user_mock.go -package=manager
	mockgen -source=internal/helper/time/clock.go -destination=internal/helper/time/clock_mock.go -package=time
	mockgen -source=internal/helper/token/generator.go -destination=internal/helper/token/generator_mock.go -package=token
	mockgen -source=internal/security/jwt/manager.go -destination=internal/security/jwt/manager_mock.go -package=jwt

test-unit: ## Run unit tests
	go test -count=1 -race -v ./...
