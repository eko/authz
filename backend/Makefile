.PHONY: lint schema-dump swag test-functional

lint:
	@command -v golangci-lint > /dev/null || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	golangci-lint run

schema-dump:
	docker-compose exec postgres pg_dump -x -s -U root -w -d root > schema.sql

# install needed dependency:
# go install github.com/swaggo/swag/cmd/swag@latest
swag:	
	swag init \
		-g ./internal/http/routing.go \
		--parseInternal \
		--dir ./ \
		--output ./internal/http/docs/

test-functional:
	go test -count=1 --tags=functional -v ./functional --godog.tags=$(tags)