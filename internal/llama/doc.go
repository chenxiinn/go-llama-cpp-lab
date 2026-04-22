// Package llama reserves the native integration boundary for libllama.
//
// Phase 1 adds the first cgo bridge behind the `llama` build tag. Phase 2
// builds on that with thin model, context, sampler, tokenization, and decode
// wrappers while keeping raw C pointers and buffers inside this package.
package llama
