# AGENTS.md

## Build, Lint, and Test Commands
- **Build:** `go build ./...`
- **Run:** `go run ./cmd/dccprint/main.go`
- **Test all:** `go test ./...`
- **Test single file:** `go test ./internal/components/scripts/scripts_test.go`
- **Test single function:** `go test -run ^TestRemoveGeneratedScripts$ ./internal/components/scripts/scripts_test.go`
- **Dependencies:** `go mod tidy`

## Code Style Guidelines
- **Imports:**
  - Standard library first, then third-party, then local modules, each group separated by a blank line.
- **Formatting:**
  - Use `gofmt` or `go fmt ./...` for formatting.
- **Types & Naming:**
  - Use `CamelCase` for types, structs, and exported functions.
  - Use `camelCase` for variables and unexported functions.
  - Constants use `CamelCase` or `ALL_CAPS` if appropriate.
- **Error Handling:**
  - Return errors, do not panic except in tests or main.
  - Use `log.Fatal` only in main or test code.
- **Tests:**
  - Place tests in `*_test.go` files, use Go's `testing` package.
- **General:**
  - Prefer explicit over implicit code.
  - Keep functions small and focused.
  - No Cursor or Copilot rules present.
