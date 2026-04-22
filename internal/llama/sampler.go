//go:build llama
// +build llama

package llama

/*
#include "llama.h"
*/
import "C"

import "fmt"

const DefaultSamplerSeed = ^uint32(0)

type SamplerConfig struct {
	UseGreedy   bool
	TopK        int
	TopP        float32
	Temperature float32
	Seed        uint32
}

type Sampler struct {
	ptr *C.struct_llama_sampler
}

func DefaultSamplerConfig() SamplerConfig {
	return SamplerConfig{
		UseGreedy:   true,
		TopK:        40,
		TopP:        0.95,
		Temperature: 0.8,
		Seed:        DefaultSamplerSeed,
	}
}

func NewSampler(config SamplerConfig) (*Sampler, error) {
	params := C.llama_sampler_chain_default_params()
	chain := C.llama_sampler_chain_init(params)
	if chain == nil {
		return nil, newError("sampler_init", "llama_sampler_chain_init returned nil")
	}

	fail := func(message string) (*Sampler, error) {
		C.llama_sampler_free(chain)
		return nil, newError("sampler_init", message)
	}

	if !config.UseGreedy {
		if config.TopK > 0 {
			topK := C.llama_sampler_init_top_k(C.int32_t(config.TopK))
			if topK == nil {
				return fail("llama_sampler_init_top_k returned nil")
			}
			C.llama_sampler_chain_add(chain, topK)
		}
		if config.TopP > 0 {
			topP := C.llama_sampler_init_top_p(C.float(config.TopP), 1)
			if topP == nil {
				return fail("llama_sampler_init_top_p returned nil")
			}
			C.llama_sampler_chain_add(chain, topP)
		}
		if config.Temperature > 0 {
			temp := C.llama_sampler_init_temp(C.float(config.Temperature))
			if temp == nil {
				return fail("llama_sampler_init_temp returned nil")
			}
			C.llama_sampler_chain_add(chain, temp)
		}
		dist := C.llama_sampler_init_dist(C.uint32_t(config.Seed))
		if dist == nil {
			return fail("llama_sampler_init_dist returned nil")
		}
		C.llama_sampler_chain_add(chain, dist)
	} else {
		greedy := C.llama_sampler_init_greedy()
		if greedy == nil {
			return fail("llama_sampler_init_greedy returned nil")
		}
		C.llama_sampler_chain_add(chain, greedy)
	}

	return &Sampler{ptr: chain}, nil
}

func (s *Sampler) Close() error {
	if s == nil || s.ptr == nil {
		return nil
	}

	C.llama_sampler_free(s.ptr)
	s.ptr = nil
	return nil
}

func (s *Sampler) Sample(ctx *Context) (Token, error) {
	if s == nil || s.ptr == nil {
		return 0, newError("sample", "sampler is nil")
	}
	if ctx == nil || ctx.ptr == nil {
		return 0, newError("sample", "context is nil")
	}

	token := C.llama_sampler_sample(s.ptr, ctx.ptr, -1)
	if token == C.llama_token(C.LLAMA_TOKEN_NULL) {
		return 0, newError("sample", "llama_sampler_sample returned LLAMA_TOKEN_NULL")
	}

	return Token(token), nil
}

func (s *Sampler) DebugName() string {
	if s == nil || s.ptr == nil {
		return ""
	}
	name := C.llama_sampler_name(s.ptr)
	if name == nil {
		return ""
	}
	return C.GoString(name)
}

func (s *Sampler) Accept(token Token) error {
	if s == nil || s.ptr == nil {
		return newError("sampler_accept", "sampler is nil")
	}
	C.llama_sampler_accept(s.ptr, C.llama_token(token))
	return nil
}

func (s *Sampler) String() string {
	return fmt.Sprintf("Sampler(%s)", s.DebugName())
}
