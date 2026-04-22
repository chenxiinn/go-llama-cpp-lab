# 01 Local Build Notes

- Date: 2026-04-22
- Phase: 0 Foundation
- Status: Active convention for local development

## Goal

Phase 0 does not link to `llama.cpp` yet. It defines the local filesystem
contract that Phase 1 will use for the first `cgo` bridge.

## Local Path Convention

Use these environment variables during local development:

- `LLAMA_CPP_DIR`
  - absolute path to the local `llama.cpp` checkout
- `LLAMA_CPP_BUILD_DIR`
  - absolute path to the local `llama.cpp` build output
  - if unset, the project assumes `${LLAMA_CPP_DIR}/build`

Example:

```bash
export LLAMA_CPP_DIR="$HOME/src/llama.cpp"
export LLAMA_CPP_BUILD_DIR="$LLAMA_CPP_DIR/build"
```

## Expected Upstream Layout

The project assumes `llama.cpp` is built outside this repo with its normal
CMake workflow and that headers remain under the source tree while libraries
land under the build tree.

Phase 1 will verify the exact library directories present on this machine
before hardening the `cgo` link flags.

## Current Repo Build

Phase 0 compiles only the Go scaffold:

```bash
go build ./...
```

The placeholder binaries expose the future runtime flags but do not start
inference yet:

```bash
go run ./cmd/chat --help
go run ./cmd/server --help
```

## Planned Phase 1 Wiring

When the native bridge lands, the local `llama.cpp` paths above will drive:

- include paths for `libllama` headers
- linker search paths for the local build output
- developer troubleshooting when local artifacts are missing or stale

The exact `CGO_CFLAGS` and `CGO_LDFLAGS` values are intentionally deferred to
Phase 1, because they must match the actual build output layout of the local
`llama.cpp` checkout rather than a guessed static recipe.
