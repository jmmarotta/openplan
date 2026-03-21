# OpenPlan Agents

## Working Rules

- Keep the CLI package thin; push parsing, validation, and filesystem policy into owned internal packages.
- Treat `.plans/` files as the source of truth. Do not mirror plan bodies into chat output.
- Preserve deterministic output so tests and editor integrations stay stable.
- Prefer focused tests next to the package that owns the behavior.

## Verification

- Run `go test ./...` after meaningful changes.
- Run `golangci-lint run` when available.
- For CLI changes, smoke test `init`, `list`, and `show` in a temp directory when practical.
