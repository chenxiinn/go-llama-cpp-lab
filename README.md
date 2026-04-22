# go-llama-cpp-lab

`go-llama-cpp-lab` is a learning and validation project for embedding `llama.cpp` directly into a Go application through `cgo` and `libllama`.

The project goal is not to build a generic chat product. The goal is to understand, step by step, how a local `GGUF` model can be loaded, driven, observed, and exposed from Go.

## Project Goal

This repository is focused on one technical question:

How do we build a local Go application that:

- calls `libllama` directly through `cgo`
- loads a local `GGUF` model
- supports multi-turn interactions
- supports streaming token output
- exposes the same runtime through:
  - a reusable Go package
  - a CLI demo
  - a local HTTP server
- emits structured JSON logs for full tracing and debugging

## Documentation

Current docs:

- [Master Plan](docs/00-master-plan.md)
- [Design Spec](docs/superpowers/specs/2026-04-22-go-llama-cpp-mvp-design.md)

## MVP Scope

The planned MVP is intentionally narrow:

- local only
- text only
- `Go + cgo + libllama`
- one loaded model
- one active generation at a time
- multi-turn state managed in Go
- streaming output
- JSON Lines logging

Out of scope for the MVP:

- multimodal
- tool calling
- automatic summarization
- concurrent generation
- distributed serving

## Planned Architecture

```text
GGUF model
   |
   v
libllama (llama.cpp)
   |
   v
internal/llama
   |
   v
pkg/chat
   |
   +--> cmd/chat
   |
   +--> cmd/server
```

Planned package layout:

```text
cmd/
  chat/
  server/
internal/
  llama/
pkg/
  chat/
docs/
```

## Expected Local Environment

The current design assumes:

- Go installed locally
- a local `llama.cpp` checkout available on disk
- a local `GGUF` model available on disk
- macOS or another environment where `cgo` can link against the local `llama.cpp` build output

The exact build instructions will be added as the native bridge and runtime are implemented.

Phase 0 has landed the Go module scaffold and placeholder binaries. The local
build and path-discovery convention now lives in [docs/01-local-build.md](docs/01-local-build.md).

## Why This Repo Exists

There are many ways to use `llama.cpp` from Go, including shelling out to a server or CLI wrapper.

This repository exists specifically to learn the direct path:

- native lifecycle management
- `cgo` boundaries
- model loading
- tokenization
- decode and sample loops
- streaming
- debugging and observability

That is why the repo is named as a lab rather than an application.

## Repository Hygiene

Two conventions are already in place:

- planning docs live under `docs/`
- local machine noise and build outputs should stay untracked via `.gitignore`

## Current Status

Phase 0 is implemented:

- `go.mod` exists
- `cmd/chat` and `cmd/server` compile as placeholders
- `pkg/chat` owns the initial runtime flag surface
- `internal/llama` reserves the native boundary and local path convention

Current verification command:

```bash
go build ./...
```
