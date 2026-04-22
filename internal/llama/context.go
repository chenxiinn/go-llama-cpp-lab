//go:build llama
// +build llama

package llama

/*
#include "llama.h"
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type ContextConfig struct {
	ContextSize  int
	BatchSize    int
	Threads      int
	ThreadsBatch int
	OffloadKQV   bool
	OpOffload    bool
}

type Context struct {
	ptr   *C.struct_llama_context
	model *Model
}

func DefaultContextConfig() ContextConfig {
	return ContextConfig{
		ContextSize:  1024,
		BatchSize:    512,
		Threads:      runtime.GOMAXPROCS(0),
		ThreadsBatch: runtime.GOMAXPROCS(0),
		OffloadKQV:   false,
		OpOffload:    false,
	}
}

func (m *Model) NewContext(config ContextConfig) (*Context, error) {
	if m == nil || m.ptr == nil {
		return nil, newError("context_init", "model is nil")
	}

	params := C.llama_context_default_params()
	if config.ContextSize > 0 {
		params.n_ctx = C.uint32_t(config.ContextSize)
	}
	if config.BatchSize > 0 {
		params.n_batch = C.uint32_t(config.BatchSize)
		params.n_ubatch = C.uint32_t(config.BatchSize)
	}
	if config.Threads > 0 {
		params.n_threads = C.int32_t(config.Threads)
	}
	if config.ThreadsBatch > 0 {
		params.n_threads_batch = C.int32_t(config.ThreadsBatch)
	}
	params.offload_kqv = C.bool(config.OffloadKQV)
	params.op_offload = C.bool(config.OpOffload)

	ptr := C.llama_init_from_model(m.ptr, params)
	if ptr == nil {
		return nil, newError("context_init", "llama_init_from_model returned nil")
	}

	return &Context{
		ptr:   ptr,
		model: m,
	}, nil
}

func (c *Context) Close() error {
	if c == nil || c.ptr == nil {
		return nil
	}

	C.llama_free(c.ptr)
	c.ptr = nil
	c.model = nil
	return nil
}

func (c *Context) Decode(tokens []Token) error {
	if c == nil || c.ptr == nil {
		return newError("decode", "context is nil")
	}
	if len(tokens) == 0 {
		return newError("decode", "tokens are empty")
	}

	ctokens := make([]C.llama_token, len(tokens))
	for i, token := range tokens {
		ctokens[i] = C.llama_token(token)
	}

	batch := C.llama_batch_get_one((*C.llama_token)(unsafe.Pointer(&ctokens[0])), C.int32_t(len(ctokens)))
	code := C.llama_decode(c.ptr, batch)
	runtime.KeepAlive(ctokens)

	switch code {
	case 0:
		return nil
	case 1:
		return newError("decode", "llama_decode could not find a KV slot for the batch")
	case 2:
		return newError("decode", "llama_decode aborted")
	case -1:
		return newError("decode", "llama_decode rejected the input batch")
	default:
		return newError("decode", fmt.Sprintf("llama_decode returned %d", int(code)))
	}
}

func (c *Context) DecodeOne(token Token) error {
	return c.Decode([]Token{token})
}

func (c *Context) LogitsAt(index int) unsafe.Pointer {
	if c == nil || c.ptr == nil {
		return nil
	}
	return unsafe.Pointer(C.llama_get_logits_ith(c.ptr, C.int32_t(index)))
}
