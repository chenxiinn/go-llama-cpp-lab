# 00 Master Plan

- Date: 2026-04-22
- Status: Active implementation master plan
- Source spec: `docs/superpowers/specs/2026-04-22-go-llama-cpp-mvp-design.md`

## Purpose

This file is the implementation roadmap for the project.

The goal is to build, step by step, a local text-only LLM application in Go that:

- calls `libllama` directly through `cgo`
- loads a local `GGUF` model
- supports multi-turn chat
- supports streaming output
- exposes the same core through:
  - a reusable Go package
  - a CLI chat app
  - a local HTTP server
- emits full JSON Lines logs for tracing and debugging

This file is intentionally phase-based. Each phase has one clear purpose and one clear exit condition. Do not start the next phase until the current phase is stable enough to verify.

## How To Use This File

Implementation should follow the phases in order.

For each phase:

1. complete the listed work
2. verify the exit criteria
3. only then move to the next phase

This prevents the project from turning into a large, hard-to-debug integration effort where too many unknowns are introduced at once.

## Scope Lock

This master plan is only for the approved MVP:

- local only
- text only
- `Go + cgo + libllama`
- one loaded model
- one active generation at a time
- multi-turn chat
- streaming output
- JSON structured logs

Out of scope for this plan:

- multimodal
- tool calling
- auto summarization
- concurrent generation
- distributed serving

## Phase Status

| Phase | Status | Notes |
| --- | --- | --- |
| 0 | Completed on 2026-04-22 | Go module scaffold, placeholder binaries, and local build notes landed |
| 1 | Completed on 2026-04-22 | Minimal `cgo` bridge, smoke test, and local config-file model path workflow landed |
| 2 | Completed on 2026-04-22 | Model/context/sampler wrappers landed and the CPU-only runtime smoke test passes against a local GGUF |
| 3 | Pending | Waiting on Phase 2 |
| 4 | Pending | Waiting on Phase 3 |
| 5 | Pending | Waiting on Phase 4 |
| 6 | Pending | Waiting on Phase 4 |
| 7 | Pending | Waiting on Phases 5 and 6 |
| 8 | Pending | Waiting on earlier phases |

## Phase Summary

| Phase | Name | Main Outcome |
| --- | --- | --- |
| 0 | Foundation | Repo and build skeleton exists and can compile a minimal Go program with `cgo` wiring |
| 1 | Native Bridge | Go can call a tiny `libllama` wrapper and prove the C boundary works |
| 2 | Model Runtime | Go can load a `GGUF` model, create context/sampler, and run one-shot inference smoke tests |
| 3 | Chat Core | Go can manage sessions, message history, chat template rendering, and prompt delta logic |
| 4 | Streaming Loop | Go can generate tokens incrementally and stop correctly |
| 5 | CLI Demo | `cmd/chat` can run a local multi-turn streaming conversation |
| 6 | HTTP Server | `cmd/server` can create sessions and stream responses over HTTP |
| 7 | Observability | Full JSON log chain exists across app, chat, and llama layers |
| 8 | Hardening | Tests, failure handling, and final acceptance checks are in place |

## Phase 0: Foundation

### Why this phase exists

Before touching `llama.cpp`, the project needs a stable Go module layout and a predictable local build path. This phase removes "empty repo" uncertainty.

### What to do

- Initialize the Go module
- Create the top-level folder structure:
  - `cmd/chat`
  - `cmd/server`
  - `internal/llama`
  - `pkg/chat`
  - `docs`
- Add a minimal `main.go` placeholder for CLI and server
- Add the first build notes describing where `llama.cpp` lives locally and how the project links to it
- Decide how the project will find `llama.cpp` headers and built artifacts during local development
- Add basic config parsing skeleton so future phases have a stable place to plug in flags

### What this phase is really proving

It proves the project is no longer "an idea in an empty folder". It becomes a real Go codebase with a defined shape.

### Deliverables

- `go.mod`
- initial directory skeleton
- compilable placeholder binaries
- local build instructions

### Exit criteria

- `go build ./...` succeeds even if the llama integration is still stubbed
- the repo layout matches the design boundaries
- the next phase can focus only on the native boundary

## Phase 1: Native Bridge

### Why this phase exists

The highest-risk unknown is not chat logic. It is whether Go can cleanly talk to `libllama` through `cgo` in this repo and on this machine.

### What to do

- Add the first `cgo` bridge files under `internal/llama`
- Link against the local `llama.cpp` build output
- Expose the smallest possible native functions first
- Start with something intentionally tiny, such as:
  - backend init
  - a version/info probe if available
  - a no-op wrapper proving C function calls work
- Centralize C pointer handling in one place
- Set up error translation from C boundary failures into Go errors

### What this phase is really proving

It proves the toolchain works:

- Go compiler
- `cgo`
- include paths
- linker paths
- local `llama.cpp` artifacts

If this phase fails, nothing above it matters yet.

### Deliverables

- `internal/llama/bridge.go`
- `internal/llama/bridge.c`
- `internal/llama/bridge.h`
- first smoke test or tiny executable path proving the bridge works

### Exit criteria

- a Go test or tiny binary can call into the bridge successfully
- build and link errors are resolved
- C resources and ownership rules are documented in code comments or package docs

## Phase 2: Model Runtime

### Why this phase exists

After the C boundary is proven, the next risk is the actual runtime lifecycle:

- load model
- create context
- create sampler
- run one minimal decode path

This phase turns the bridge into a usable inference runtime.

### What to do

- Implement model loading from `GGUF`
- Implement context creation with config like `n_ctx`
- Implement sampler creation with first-pass sampling settings
- Read model metadata, especially the default chat template
- Implement tokenize
- Implement single prompt decode
- Implement sample-next-token
- Implement token-to-piece conversion
- Add cleanup for model, context, and sampler

### What this phase is really proving

It proves that Go can perform actual inference operations through `libllama`, not just call a native stub.

### Deliverables

- model wrapper
- context wrapper
- sampler wrapper
- smoke tests for:
  - model load
  - chat template read
  - tokenize
  - one decode step
  - one sampled token

### Exit criteria

- the runtime can load the chosen local `GGUF`
- the model chat template is readable
- a one-shot inference smoke path works end to end
- resource cleanup is reliable

## Phase 3: Chat Core

### Why this phase exists

Raw inference is not enough. The application needs a stable Go-native abstraction for multi-turn chat so that CLI and HTTP can reuse the same behavior.

### What to do

- Create `pkg/chat`
- Define `Runtime`
- Define `Session`
- Define `Message`
- Store ordered conversation history in Go
- Render prompt from model chat template
- Track `prev_prompt_len`
- Compute prompt delta for each new user turn
- Append assistant output back into session history after generation
- Add stop reason handling at the chat layer

### What this phase is really proving

It proves the system can preserve multi-turn state in Go without mixing session semantics into the native layer.

### Deliverables

- `pkg/chat/runtime.go`
- `pkg/chat/session.go`
- `pkg/chat/types.go`
- tests for:
  - message accumulation
  - prompt delta logic
  - assistant message append

### Exit criteria

- a session can hold multiple turns
- prompt rendering is deterministic
- prompt delta logic works across at least two turns

## Phase 4: Streaming Loop

### Why this phase exists

This phase turns "inference exists" into "usable incremental chat output exists".

Without this phase, the project may generate text internally but still fail the user-facing requirement of streaming output.

### What to do

- Build the generation loop around:
  - decode
  - sample
  - token-to-piece
  - emit piece
  - feed sampled token back into decode
- Define a streaming abstraction used by both CLI and HTTP
- Add generation stop conditions:
  - end-of-generation token
  - max token limit
  - context cancellation
  - context exhaustion
  - stream write failure
- Ensure final assistant text is accumulated while pieces are streamed

### What this phase is really proving

It proves the core engine can behave like a chat system instead of just returning a single batch result.

### Deliverables

- `pkg/chat/stream.go`
- streaming-capable send/generate path
- cancellation-aware generation flow

### Exit criteria

- token pieces are emitted incrementally
- a caller can cancel generation
- the final assistant text remains available for session history

## Phase 5: CLI Demo

### Why this phase exists

The CLI is the fastest real end-to-end way to validate human interaction with the system.

It exposes whether the runtime actually feels correct when used turn by turn.

### What to do

- Implement `cmd/chat`
- Parse runtime flags such as:
  - model path
  - context size
  - GPU layer count
  - max tokens
- Create one session at startup
- Read terminal input in a loop
- Stream assistant output to `stdout`
- Keep logs separate from chat output
- Support user interruption of the active generation

### What this phase is really proving

It proves the local developer can actually sit down and chat with the model using the Go program.

### Deliverables

- `cmd/chat/main.go`
- local usage instructions
- CLI smoke test notes

### Exit criteria

- `go run ./cmd/chat --model /path/to/model.gguf` starts
- at least two turns work in one session
- the second answer reflects the first turn
- output is visibly streamed, not buffered until the end

## Phase 6: HTTP Server

### Why this phase exists

The HTTP layer proves the runtime is reusable beyond the terminal and can serve as an integration point for other local tools or frontends.

### What to do

- Implement `cmd/server`
- Add in-memory session store
- Add:
  - `POST /sessions`
  - `POST /sessions/{id}/messages/stream`
- Use SSE for streaming output
- Attach request-level `trace_id`
- Handle:
  - session not found
  - busy runtime
  - client disconnect
  - invalid input

### What this phase is really proving

It proves the chat core is not coupled to the CLI and can act as a local service boundary.

### Deliverables

- `cmd/server/main.go`
- HTTP handlers
- SSE stream output
- curl examples for manual testing

### Exit criteria

- server starts locally
- a client can create a session
- a client can send a streamed chat request
- disconnect or cancellation does not kill the process

## Phase 7: Observability

### Why this phase exists

This project explicitly requires a full debug trail. Without this phase, problems inside `cgo`, inference, prompt handling, or HTTP flow will be expensive to diagnose.

### What to do

- Add a shared JSON Lines logger
- Emit logs at three layers:
  - app
  - chat
  - llama
- Add correlation fields:
  - `trace_id`
  - `session_id`
  - `component`
  - `event`
  - `elapsed_ms`
- Add standard lifecycle events for:
  - runtime init
  - model load
  - session creation
  - chat turn start
  - prompt tokenization
  - decode steps
  - streamed token events
  - stop reason
  - request completion
- Add optional debug modes:
  - prompt dump
  - token trace
- Ensure CLI answer text is not mixed with JSON logs

### What this phase is really proving

It proves the system is explainable. You can follow a request from entry to failure or completion and see where time and state changed.

### Deliverables

- shared logging package or module
- JSON log schema in code
- documented event names
- sample log snippets in docs or test notes

### Exit criteria

- every request can be traced end to end
- failures are visible at the correct layer
- logs are valid JSON objects, one per line

## Phase 8: Hardening

### Why this phase exists

The project is not finished just because it "works once". This phase converts a working prototype into a stable MVP with explicit acceptance checks.

### What to do

- Add missing tests
- Add startup error coverage
- Add request error coverage
- Verify cancellation behavior
- Verify busy-runtime behavior
- Verify context exhaustion behavior
- Run CLI smoke tests
- Run HTTP smoke tests
- Validate log output quality
- Remove dead paths and simplify any unstable abstractions introduced earlier

### What this phase is really proving

It proves the project meets the agreed MVP contract rather than only passing happy-path demos.

### Deliverables

- test suite updates
- manual verification checklist
- cleaned-up error messages
- final acceptance notes

### Exit criteria

- all MVP acceptance criteria from the source spec are satisfied
- typical failures can be diagnosed from logs
- the repo is in a state where future iteration can start from a stable base

## Phase Dependencies

The phases are intentionally ordered by risk:

1. Foundation before native work
2. Native work before real model runtime
3. Model runtime before chat state
4. Chat state before streaming UX
5. Core streaming before CLI and HTTP shells
6. Shells before full observability and hardening

If a later phase reveals a design flaw in an earlier one, fix the earlier phase first. Do not paper over foundational problems at the edge layers.

## Suggested Implementation Rhythm

Each phase should ideally follow the same loop:

1. create the smallest code that proves the phase goal
2. run the narrowest possible verification
3. fix correctness and lifecycle issues
4. only then add the next layer

This keeps the project convergent and prevents the common failure mode where CLI, HTTP, logging, and native integration are all built at once and become impossible to debug.

## First Practical Milestone

The first serious milestone is not the HTTP server.

It is:

- a Go binary that links through `cgo`
- loads the target `GGUF`
- reads the model chat template
- tokenizes a prompt
- decodes at least one step
- samples at least one token

If that milestone is not solid, every later phase is built on sand.

## Final Deliverable Definition

The project is considered complete for this master plan when all of the following are true:

- a reusable Go chat runtime exists
- a CLI demo exists
- a local HTTP server exists
- both support multi-turn chat
- both support streaming output
- logs are JSON Lines and usable for debugging
- failures can be localized to startup, template, tokenize, decode, stream, cancel, or context exhaustion

## Future Work After This Plan

Not part of this master plan, but natural next steps after MVP:

- concurrent generation support
- better session persistence
- context trimming or summarization
- metrics export
- benchmarks
- model hot-swap
- tool calling
- multimodal support
