# Go + llama.cpp MVP Design

- Date: 2026-04-22
- Status: Approved for planning
- Scope: Minimal MVP design only, no implementation yet

## Summary

Build a minimal local LLM application in Go that integrates directly with `llama.cpp` through `cgo` and `libllama`, using a local `GGUF` model.

The MVP must support:

- Multi-turn chat
- Streaming token output
- Three entry points built on the same core runtime:
  - reusable Go package
  - CLI chat demo
  - local HTTP server
- Full structured logging in JSON Lines format for end-to-end tracing and debugging

The design deliberately stays text-only. It does not include multimodal inputs, tool calling, automatic context summarization, or concurrent generation.

## Goals

- Prove that Go can call `libllama` directly through `cgo`
- Load a local `GGUF` model from Go
- Maintain multi-turn chat state in Go
- Stream generated text incrementally to callers
- Reuse one inference core across library, CLI, and HTTP layers
- Provide enough structured logs to explain how a request moved through the system and where it failed

## Non-Goals

- No multimodal support
- No Python sidecar or `llama-server` dependency for inference
- No tool calling or function calling
- No distributed execution
- No advanced context window management such as summarization or automatic truncation
- No concurrent generation on the same runtime in the first version
- No production deployment concerns beyond local use

## Constraints

- The model format is `GGUF`
- Inference must happen through `Go + cgo + libllama`
- The program must be locally runnable on the developer machine
- The first MVP must prioritize observability and debuggability over feature breadth
- Logs must be JSON only

## Recommendation

Use a layered design:

1. `internal/llama`
   The thinnest possible `cgo` wrapper over `libllama`
2. `pkg/chat`
   Go-native chat runtime and session abstraction
3. `cmd/chat`
   Interactive CLI built on `pkg/chat`
4. `cmd/server`
   Local HTTP server built on `pkg/chat`

This is preferred over a single large binary with mixed responsibilities because it verifies the direct Go integration while keeping CLI and HTTP as thin shells over one core inference path.

## High-Level Architecture

```text
GGUF model
   |
   v
libllama (llama.cpp)
   |
   v
internal/llama        <- cgo wrapper around model/context/sampler/token APIs
   |
   v
pkg/chat              <- sessions, message history, prompt rendering, streaming API
   |
   +--> cmd/chat      <- terminal UI
   |
   +--> cmd/server    <- local HTTP API
```

## Proposed Project Structure

```text
cmd/
  chat/
    main.go
  server/
    main.go
internal/
  llama/
    bridge.go
    bridge.c
    bridge.h
    model.go
    context.go
    sampler.go
    errors.go
pkg/
  chat/
    runtime.go
    session.go
    stream.go
    types.go
    config.go
    errors.go
docs/
  superpowers/
    specs/
      2026-04-22-go-llama-cpp-mvp-design.md
```

The exact file split may change slightly during implementation, but the boundary between `internal/llama` and `pkg/chat` should remain stable.

## Core Components

### `internal/llama`

Purpose:
Expose the minimum `libllama` functionality needed by the MVP while keeping raw C details out of the rest of the codebase.

Responsibilities:

- Initialize llama backend
- Load model from `GGUF`
- Create and free context
- Create and free sampler
- Read model metadata such as chat template
- Tokenize prompt text
- Decode batches
- Sample next token
- Convert token IDs back to text pieces
- Surface `libllama` errors as Go errors
- Register low-level llama logging callback and bridge it into structured application logs

Deliberately out of scope:

- Business logic
- Session semantics
- HTTP concerns
- CLI interaction

### `pkg/chat`

Purpose:
Provide a Go-native API for multi-turn streaming chat.

Responsibilities:

- Own runtime configuration
- Expose session creation
- Maintain message history per session
- Render prompt using the model chat template
- Compute prompt delta for subsequent turns
- Drive generation loop
- Stream text pieces to callers
- Track stop reason
- Emit structured logs for each stage

### `cmd/chat`

Purpose:
Provide a local interactive terminal demo.

Responsibilities:

- Parse startup flags
- Load runtime
- Create one session
- Read user input line by line
- Stream assistant output to `stdout`
- Send logs to `stderr` or a log file

### `cmd/server`

Purpose:
Provide a local HTTP API on top of the same runtime.

Responsibilities:

- Parse startup flags
- Load runtime
- Manage in-memory sessions
- Expose streaming chat endpoint
- Return structured HTTP errors
- Emit request-scoped logs with `trace_id`

## Runtime Model

The runtime keeps a single model loaded in memory and exposes session creation.

### Runtime

Fields:

- model handle
- context handle
- sampler handle
- runtime config
- process-wide logger
- mutex guarding generation

Important first-version rule:

- Multiple sessions may exist
- Only one generation may run at a time

This avoids early complexity around concurrent access to the same `llama_context` and KV state.

### Session

Fields:

- `session_id`
- ordered `messages`
- `prev_prompt_len`
- timestamps
- optional session metadata for logs

Each session owns its own chat history in Go. The first MVP does not persist sessions to disk.

## Chat Data Model

Messages are stored in Go as:

```go
type Message struct {
    Role    string
    Content string
}
```

Supported roles:

- `system`
- `user`
- `assistant`

The first MVP does not support tool role messages.

## Generation Flow

For one chat turn:

1. Caller invokes `Session.Send(ctx, userText, stream)`
2. User message is appended to the session history
3. The model chat template is read from the loaded model
4. Full conversation is rendered into prompt text
5. `prev_prompt_len` is used to derive only the incremental prompt for this turn
6. Incremental prompt is tokenized through `libllama`
7. Prompt tokens are decoded into the context
8. Generation loop begins:
   - sample token
   - stop if end-of-generation token
   - convert token to text piece
   - emit piece to stream
   - feed sampled token back into decode
9. Final assistant text is appended to session history
10. `prev_prompt_len` is recomputed from the full formatted conversation for the next round

This follows the incremental prompt pattern demonstrated by the current `llama.cpp` simple chat example and keeps the Go side aligned with the upstream library behavior.

## Streaming Model

The chat package exposes one streaming abstraction used by both CLI and HTTP.

Possible shape:

```go
type Streamer interface {
    OnToken(text string) error
    OnComplete(result Result) error
    OnError(err error) error
}
```

The exact API can vary, but the behavior must stay the same:

- token pieces are emitted incrementally
- completion is explicit
- cancellation is propagated through `context.Context`

### CLI Streaming

- assistant text goes to `stdout`
- structured logs go to `stderr` or file
- `Ctrl+C` cancels the active generation without corrupting the process

### HTTP Streaming

Preferred transport:

- Server-Sent Events

Reason:

- simpler than WebSocket for first-version one-way token streaming
- easy to inspect with curl and browser tooling

Chosen endpoint shape:

- `POST /sessions`
- `POST /sessions/{id}/messages/stream`

These routes are part of the MVP contract and should not be renamed during implementation unless the spec is updated.

## Prompt Handling

Prompt construction must use the model-provided chat template from the loaded `GGUF`.

Rules:

- If the model has no chat template, startup fails
- No hand-written fallback prompt template is introduced in the MVP
- Prompt dumping for debugging is optional and controlled by config
- Prompt dumps are disabled by default or truncated to avoid noisy logs

This keeps behavior deterministic and closer to the model’s intended formatting.

## Configuration

The runtime should accept startup configuration through flags and optionally environment variables.

Required config:

- `model_path`
- `n_ctx`
- `n_gpu_layers`
- `max_tokens`
- `temperature`
- `min_p`
- log output destination

Optional config:

- token trace enabled
- prompt dump enabled
- HTTP bind address

The first MVP should prefer explicit startup flags so behavior is obvious while testing.

## Logging and Traceability

Logging is a first-class MVP requirement.

### Log Format

- JSON Lines only
- one event per line

### Log Destinations

- default: `stderr`
- optional: file output

### Correlation Fields

Every relevant event should include:

- `ts`
- `level`
- `component`
- `event`
- `trace_id`
- `session_id`
- `elapsed_ms`
- `model_path`
- `error`

Event-specific fields may include:

- `prompt_chars`
- `prompt_tokens`
- `generated_tokens`
- `token_id`
- `piece`
- `stop_reason`
- `http_method`
- `http_path`
- `status_code`

### Log Layers

#### App Layer

- process startup
- flag parsing
- CLI command lifecycle
- HTTP request lifecycle
- session creation
- cancellation

#### Chat Layer

- turn started
- template applied
- prompt delta computed
- tokenization complete
- decode step progressed
- token streamed
- generation stopped
- turn completed

#### Llama Layer

- backend init
- model load started/completed/failed
- context init started/completed/failed
- sampler init started/completed/failed
- low-level llama error callback messages

### Standard Event Names

- `runtime.init.started`
- `runtime.init.completed`
- `runtime.init.failed`
- `model.load.started`
- `model.load.completed`
- `model.load.failed`
- `session.created`
- `chat.turn.started`
- `chat.template.applied`
- `chat.prompt.delta_computed`
- `chat.prompt.tokenized`
- `chat.decode.step`
- `chat.token.streamed`
- `chat.generation.stopped`
- `chat.turn.completed`
- `chat.turn.failed`
- `chat.cancelled`
- `http.request.started`
- `http.request.completed`

### Debug Modes

Two optional high-verbosity log modes should be supported:

- `prompt dump`
  Records the rendered prompt text, preferably truncated
- `token trace`
  Records each emitted token id and text piece

These modes are for debugging and should be explicitly enabled.

## Error Handling

The MVP should optimize for diagnosable failure.

### Startup Errors

- model file missing
- invalid `GGUF`
- dynamic library or backend init failure
- context creation failure
- sampler creation failure
- missing chat template

These should fail fast at startup with structured errors.

### Request Errors

- session not found
- empty user input
- runtime busy
- context exhausted
- tokenization failure
- decode failure
- stream write failure
- request cancelled

These should return a clear error to the caller and a structured event in logs.

### Cancellation Behavior

- CLI cancellation interrupts current generation
- HTTP client disconnect cancels current generation
- cancellation must not crash the runtime

### Context Window Behavior

The first MVP does not implement summarization or auto-truncation.

If a turn exceeds context capacity:

- generation stops
- an explicit `context_exhausted` error is returned
- logs record where the overflow was detected
- user can create a new session or shorten history

This keeps the first version predictable and avoids hidden behavior.

## Concurrency Model

First version policy:

- one loaded model
- one active generation at a time
- multiple sessions may be stored in memory
- concurrent generation requests are rejected immediately with a clear `busy` response

This makes behavior easier to understand during debugging.

## Resource Lifecycle

The code must centralize ownership of all C resources.

Objects with explicit lifecycle:

- llama backend
- model
- context
- sampler
- C strings allocated for bridge calls

Requirements:

- creation and destruction happen in one layer
- callers outside `internal/llama` never touch raw C pointers
- cleanup order is deterministic
- shutdown logs include success or failure

## Testing Strategy

### `internal/llama`

Use small integration-oriented tests against the real library where practical.

Validate:

- model load
- chat template read
- tokenize
- decode one prompt
- sample one token
- token to piece conversion

The point of this layer is real interop, so mocking here should be minimal.

### `pkg/chat`

Use Go unit tests with a fake runner abstraction where useful.

Validate:

- message accumulation
- prompt delta logic
- assistant response append
- stop reason propagation
- cancellation propagation

### `cmd/server`

Use integration tests for:

- session creation
- streaming response
- invalid session id
- busy runtime response
- cancellation on client disconnect

### `cmd/chat`

Keep automated testing lightweight.

Validate:

- process startup
- basic prompt path

Use manual smoke testing for the interactive loop.

### Logging Tests

Validate:

- every request has a `trace_id`
- every session event has a `session_id`
- logs are valid JSON objects, one per line
- key lifecycle events appear in the expected order

## MVP Acceptance Criteria

The MVP is complete when all of the following are true:

1. `go run ./cmd/chat --model /path/to/model.gguf` starts successfully
2. The CLI supports at least two chat rounds in one session
3. The second round reflects the first round’s context
4. Assistant output streams incrementally in the CLI
5. `go run ./cmd/server --model /path/to/model.gguf` starts successfully
6. The HTTP API supports session creation and streamed chat responses
7. Cancellation stops active generation without killing the process
8. Logs are JSON Lines and contain enough information to trace a request end-to-end
9. Common failures can be located by stage: startup, template, tokenize, decode, stream, cancel, context exhaustion

## Risks and Tradeoffs

### `cgo` Complexity

Direct `libllama` integration proves the desired architecture, but it increases lifecycle and debugging complexity compared with calling `llama-server`.

This is acceptable because direct Go invocation is the explicit goal of the MVP.

### Single Active Generation

This is a deliberate limitation.

It reduces throughput, but it makes the first version easier to reason about and significantly lowers the risk of subtle context or threading bugs.

### No Automatic Context Management

This makes the MVP stricter for users, but the failure mode is visible and deterministic.

That is preferable to hidden truncation in the first version.

## Implementation Guidance for the Next Step

The next step after this spec is an implementation plan, not coding directly from the raw conversation.

That plan should break work into:

1. project bootstrap and build wiring for `cgo` + `llama.cpp`
2. `internal/llama` bridge
3. `pkg/chat` session runtime
4. CLI wrapper
5. HTTP wrapper
6. logging and diagnostics
7. test and acceptance pass

## Final Scope Lock

This MVP is specifically:

- local
- text-only
- `Go + cgo + libllama`
- multi-turn
- streaming
- reusable as package, CLI, and HTTP server
- fully traceable through JSON structured logs

Everything outside those boundaries is deferred.
