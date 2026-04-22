# Repository Guidelines

## Project Structure & Module Organization

This repository is currently docs-first. The main planning documents are [README.md](README.md), [docs/00-master-plan.md](docs/00-master-plan.md), and [docs/superpowers/specs/2026-04-22-go-llama-cpp-mvp-design.md](docs/superpowers/specs/2026-04-22-go-llama-cpp-mvp-design.md).

Implementation should follow the planned layout:

- `cmd/chat/` for the CLI entrypoint
- `cmd/server/` for the local HTTP server
- `internal/llama/` for the `cgo` and `libllama` bridge
- `pkg/chat/` for Go-native session and streaming logic
- `docs/` for design, planning, and contributor-facing notes

Keep architecture changes aligned with `docs/00-master-plan.md`.

## Build, Test, and Development Commands

Phase 0 has landed the Go module scaffold, so do not add ad hoc scripts. Standard commands are:

- `go build ./...` to compile all packages
- `go test ./...` to run the full test suite
- `go run ./cmd/chat --model /path/to/model.gguf` to start the CLI
- `go run ./cmd/server --model /path/to/model.gguf` to start the local server

If you introduce a new command, document it in both `README.md` and this file.

Local llama.cpp path and build-linking conventions are documented in `docs/01-local-build.md`.

## Coding Style & Naming Conventions

Use standard Go formatting and keep files ASCII unless a file already requires Unicode. Run `gofmt` on edited Go files. Prefer short, explicit package boundaries: native interop stays in `internal/llama`, session logic stays in `pkg/chat`.

Naming guidelines:

- packages: lowercase, no underscores
- exported Go symbols: `CamelCase`
- unexported helpers: `camelCase`
- files: purpose-driven names such as `runtime.go`, `session.go`, `bridge.go`

## Testing Guidelines

Prefer small, focused tests. Keep native bridge tests near `internal/llama`; keep session and streaming tests near `pkg/chat`. Name tests with standard Go patterns such as `TestModelLoad` and `TestSessionPromptDelta`.

Add tests with each behavioral change, especially around tokenization, decoding, cancellation, and logging.

## Commit & Pull Request Guidelines

Use Conventional Commits with `type(scope): subject`. Example: `docs(repo): add contributor guide` or `feat(chat): add session runtime`.

Keep PRs phase-focused. Include:

- a short summary
- linked issue or plan section when relevant
- test evidence (`go test ./...`, manual CLI/HTTP checks)
- log samples or screenshots only when they clarify behavior

## Security & Configuration Tips

Do not commit local model paths, credentials, or machine-specific linker settings. Keep `GGUF` paths and `llama.cpp` locations configurable, and make local-only files stay untracked through `.gitignore`.
