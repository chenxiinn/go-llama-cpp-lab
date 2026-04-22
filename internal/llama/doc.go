// Package llama reserves the native integration boundary for libllama.
//
// Phase 1 adds the first cgo bridge behind the `llama` build tag. The bridge
// keeps raw C pointers and buffers inside this package so callers only deal
// with Go values and translated errors.
package llama
