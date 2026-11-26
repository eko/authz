## help: 💡 Display available commands
.PHONY: help
help:
	@echo '⚡️ GoFiber/Fiber Development:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## audit: 🚀 Conduct quality checks
.PHONY: audit
audit:
	go mod verify
	go vet ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

## benchmark: 📈 Benchmark code performance
.PHONY: benchmark
benchmark:
	go test ./... -benchmem -bench=. -count=4 -run=^Benchmark_$

## coverage: ☂️  Generate coverage report
.PHONY: coverage
coverage:
	go run gotest.tools/gotestsum@latest -f testname -- ./... -race -count=1 -coverprofile=/tmp/coverage.out -covermode=atomic
	go tool cover -html=/tmp/coverage.out

## format: 🎨 Fix code format issues
.PHONY: format
format:
	go run mvdan.cc/gofumpt@latest -w -l .

## markdown: 🎨 Find markdown format issues (Requires markdownlint-cli2)
.PHONY: markdown
markdown:
	markdownlint-cli2 "**/*.md" "#vendor"

## lint: 🚨 Run lint checks
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0 run ./...

## test: 🚦 Execute all tests
.PHONY: test
test:
	go run gotest.tools/gotestsum@latest -f testname -- ./... -race -count=1 -shuffle=on

## longtest: 🚦 Execute all tests 10x
.PHONY: longtest
longtest:
	go run gotest.tools/gotestsum@latest -f testname -- ./... -race -count=15 -shuffle=on

## tidy: 📌 Clean and tidy dependencies
.PHONY: tidy
tidy:
	go mod tidy -v

## betteralign: 📐 Optimize alignment of fields in structs
.PHONY: betteralign
betteralign:
	go run github.com/dkorunic/betteralign/cmd/betteralign@latest -test_files -generated_files -apply ./...

## generate: ⚡️ Generate msgp && interface implementations
.PHONY: generate
generate:
	go install github.com/tinylib/msgp@latest
	go install github.com/vburenin/ifacemaker@975a95966976eeb2d4365a7fb236e274c54da64c
	go generate ./...
