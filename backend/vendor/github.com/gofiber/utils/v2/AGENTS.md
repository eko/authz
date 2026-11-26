# AGENTS.md

Guidance for AI coding agents working in this repository.

## What this repository is

gofiber/utils is a helper library for [Fiber](https://github.com/gofiber/fiber) and its middleware ecosystem. It exports fast, low-allocation helpers; the consumers live in other modules (fiber itself, middleware, user code).

## Reviewing and proposing API changes

- Exported functions without in-module callers are by design. Do not argue that a new public function or package should move to `internal/` because grep finds no non-test usage inside this module. In-module usage counts say nothing about the value of the API; exporting helpers for downstream modules is the purpose of this repository.
- Judge new exported API instead by: contract clarity (documented preconditions, aliasing rules, edge cases), test coverage (including property or fuzz tests for bit-twiddling code), naming consistency with the existing surface, and benchmark evidence for performance claims.
- Additive API is not a SemVer break, but it is a long-term maintenance commitment. Call out sharp-edged contracts (unchecked preconditions, approximate results, unsafe memory semantics) explicitly so maintainers approve them consciously.

## Conventions

- Performance claims are verified with benchstat (base vs. head on the same machine, at least `-count=10`), never with single runs.
- `make test` runs the suite with race detector and shuffle; `make lint` runs golangci-lint; `make format` applies gofumpt; `make benchfmt` aligns the benchmark tables in README.md and `make benchupdate` re-measures them.
- README.md is a function catalog with benchmark blocks. New exported functions need a README section, and benchmark numbers for touched paths should be regenerated in the same run.
- Each benchmark block starts with the `go test` command that produced it; `make benchupdate` re-runs that command and refreshes the values of the rows it reports, leaving rows it does not measure alone. A maintainer can trigger the same thing on a pull request by commenting `/bench-readme` (arm64) or `/bench-readme-amd64`.
- Code, comments, commit messages, and PR text are always written in English.
