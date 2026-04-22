# 01 Local Build Notes

- Date: 2026-04-22
- Phase: 2 Model Runtime
- Status: Active convention for local native development

## Goal

Phase 2 now includes both the first native bridge smoke path and the first
real runtime smoke path. This document covers two separate local concerns:

- compile-time header and library discovery for `cgo`
- runtime model selection through untracked local or per-user config files

## Local Path Convention

Use these environment variables during native bridge development:

- `LLAMA_CPP_DIR`
  - absolute path to the local `llama.cpp` checkout
- `LLAMA_CPP_BUILD_DIR`
  - absolute path to the local `llama.cpp` build output
  - if unset, the project assumes `${LLAMA_CPP_DIR}/build`
- `LLAMA_INCLUDE_DIR`
  - absolute path to the directory containing `llama.h`
  - if unset, the project assumes `${LLAMA_CPP_DIR}/include`
- `LLAMA_LIB_DIR`
  - absolute path to the directory containing `libllama` and companion `ggml` dylibs
  - if unset, the project assumes `${LLAMA_CPP_BUILD_DIR}/bin`

Example:

```bash
export LLAMA_CPP_DIR="$HOME/src/llama.cpp"
export LLAMA_CPP_BUILD_DIR="$LLAMA_CPP_DIR/build"
```

## Expected Upstream Layout

The project assumes `llama.cpp` is built outside this repo with its normal
CMake workflow and that headers remain under the source tree while libraries
land under the build tree.

Phase 1 verified the exact library directories needed for local native builds.

## Runtime Config File

Do not commit real model paths.

Commit the example file:

```text
config/local.example.json
```

Create your own untracked local file:

```text
config/local.json
```

Or keep a per-user default here:

```text
~/.go-llama-cpp-lab/config.json
```

Example:

```json
{
  "model_path": "/absolute/path/to/model.gguf",
  "n_ctx": 4096,
  "gpu_layers": 0,
  "max_tokens": 256
}
```

Startup search order is:

1. explicit `--config /path/to/config.json`
2. `./config/local.json`
3. `~/.go-llama-cpp-lab/config.json`

If none of the above exists, the program exits with an error. Flags still
override values after the chosen file is loaded.

## Current Repo Build

Phase 0 compiles only the Go scaffold:

```bash
go build ./...
```

The placeholder binaries expose the future runtime flags but do not start
inference yet:

```bash
go run ./cmd/chat
go run ./cmd/server
```

## Phase 1 Native Bridge Verification

Phase 1 uses the `llama` build tag so the repo still builds without local
native artifacts by default.

Example verification command:

```bash
export LLAMA_INCLUDE_DIR="/absolute/path/to/include"
export LLAMA_LIB_DIR="/absolute/path/to/lib"

CGO_CFLAGS="-I$LLAMA_INCLUDE_DIR" \
CGO_LDFLAGS="-L$LLAMA_LIB_DIR -Wl,-rpath,$LLAMA_LIB_DIR -lllama" \
go test -tags llama ./internal/llama
```

## Phase 2 Model Runtime Verification

Phase 2 keeps the same include and library discovery flow, but now verifies a
real `GGUF` runtime path as well.

Example:

```bash
export GO_LLAMA_CPP_LAB_TEST_MODEL="/absolute/path/to/model.gguf"

CGO_CFLAGS="-I$LLAMA_INCLUDE_DIR" \
CGO_LDFLAGS="-L$LLAMA_LIB_DIR -Wl,-rpath,$LLAMA_LIB_DIR -lllama" \
go test -tags llama ./internal/llama -run TestPhase2RuntimeSmoke -v
```

## Current Native Path Wiring

The local `llama.cpp` paths above currently drive:

- include paths for `libllama` headers
- linker search paths for the local build output
- developer troubleshooting when local artifacts are missing or stale

The exact `CGO_CFLAGS` and `CGO_LDFLAGS` values intentionally stay local,
because they must match the actual build output layout of the local
`llama.cpp` checkout rather than a guessed static repo-wide recipe.
