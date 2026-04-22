# 01 Local Build Notes

- Date: 2026-04-22
- Phase: 0 Foundation
- Status: Active convention for local development

## Goal

Phase 1 adds the first native bridge smoke path. This document now covers two
separate local concerns:

- compile-time header and library discovery for `cgo`
- runtime model selection through an untracked local config file

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
  - absolute path to the directory containing `libllama`
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

Phase 1 will verify the exact library directories present on this machine
before hardening the `cgo` link flags.

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

Example:

```json
{
  "model_path": "/absolute/path/to/model.gguf",
  "n_ctx": 4096,
  "gpu_layers": 0,
  "max_tokens": 256
}
```

The CLI and server now read `config/local.json` by default. Flags can still
override values after the file is loaded.

## Current Repo Build

Phase 0 compiles only the Go scaffold:

```bash
go build ./...
```

The placeholder binaries expose the future runtime flags but do not start
inference yet:

```bash
go run ./cmd/chat --config ./config/local.json --help
go run ./cmd/server --config ./config/local.json --help
```

## Phase 1 Native Bridge Verification

Phase 1 uses the `llama` build tag so the repo still builds without local
native artifacts by default.

Example verification command:

```bash
export LLAMA_INCLUDE_DIR="/Users/chenxiinn/Documents/workspace/dev/ai_learning/upmbot/.venv/lib/python3.13/site-packages/include"
export LLAMA_LIB_DIR="/Users/chenxiinn/Documents/workspace/dev/ai_learning/upmbot/.venv/lib/python3.13/site-packages/lib"

CGO_CFLAGS="-I$LLAMA_INCLUDE_DIR" \
CGO_LDFLAGS="-L$LLAMA_LIB_DIR -Wl,-rpath,$LLAMA_LIB_DIR -lllama" \
go test -tags llama ./internal/llama
```

## Planned Phase 2 Wiring

When the native bridge lands, the local `llama.cpp` paths above will drive:

- include paths for `libllama` headers
- linker search paths for the local build output
- developer troubleshooting when local artifacts are missing or stale

The exact `CGO_CFLAGS` and `CGO_LDFLAGS` values are intentionally deferred to
Phase 1, because they must match the actual build output layout of the local
`llama.cpp` checkout rather than a guessed static recipe.
