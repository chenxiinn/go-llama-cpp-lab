//go:build llama
// +build llama

package llama

/*
#cgo LDFLAGS: -lggml
#include "llama.h"
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"runtime"
	"unsafe"
)

type Token int32

type ModelConfig struct {
	GPULayers int
	UseMMap   bool
	UseMLock  bool
}

type Model struct {
	ptr   *C.struct_llama_model
	vocab *C.struct_llama_vocab
}

func DefaultModelConfig() ModelConfig {
	return ModelConfig{
		GPULayers: 0,
		UseMMap:   true,
		UseMLock:  false,
	}
}

func LoadModel(path string, config ModelConfig) (*Model, error) {
	if path == "" {
		return nil, newError("model_load", "model path is empty")
	}

	ensureBackendInitialized()

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	params := C.llama_model_default_params()
	params.n_gpu_layers = C.int32_t(config.GPULayers)
	params.use_mmap = C.bool(config.UseMMap)
	params.use_mlock = C.bool(config.UseMLock)

	var devices [2]C.ggml_backend_dev_t
	if config.GPULayers <= 0 {
		// Leaving devices nil means libllama can still consider every
		// registered backend, which triggers Metal initialization on this
		// machine even for CPU-only runs.
		cpu := C.ggml_backend_dev_by_type(C.GGML_BACKEND_DEVICE_TYPE_CPU)
		if cpu == nil {
			return nil, newError("model_load", "ggml_backend_dev_by_type returned nil for CPU")
		}
		devices[0] = cpu
		params.devices = &devices[0]
	}

	ptr := C.llama_model_load_from_file(cpath, params)
	runtime.KeepAlive(devices)
	if ptr == nil {
		return nil, newError("model_load", fmt.Sprintf("llama_model_load_from_file returned nil for %q", path))
	}

	vocab := C.llama_model_get_vocab(ptr)
	if vocab == nil {
		C.llama_model_free(ptr)
		return nil, newError("model_load", "model vocabulary is nil")
	}

	return &Model{
		ptr:   ptr,
		vocab: (*C.struct_llama_vocab)(unsafe.Pointer(vocab)),
	}, nil
}

func (m *Model) Close() error {
	if m == nil || m.ptr == nil {
		return nil
	}

	C.llama_model_free(m.ptr)
	m.ptr = nil
	m.vocab = nil
	return nil
}

func (m *Model) ChatTemplate() (string, error) {
	if m == nil || m.ptr == nil {
		return "", newError("chat_template", "model is nil")
	}

	tmpl := C.llama_model_chat_template(m.ptr, nil)
	if tmpl == nil {
		return "", newError("chat_template", "model does not expose a default chat template")
	}

	value := C.GoString(tmpl)
	if value == "" {
		return "", newError("chat_template", "model chat template is empty")
	}

	return value, nil
}

func (m *Model) MetaString(key string) (string, error) {
	if m == nil || m.ptr == nil {
		return "", newError("model_meta", "model is nil")
	}

	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))

	size := 256
	for {
		buf := make([]C.char, size)
		n := C.llama_model_meta_val_str(m.ptr, ckey, &buf[0], C.size_t(len(buf)))
		switch {
		case n < 0:
			return "", newError("model_meta", fmt.Sprintf("metadata key %q not found", key))
		case int(n) < size:
			return C.GoString(&buf[0]), nil
		default:
			size = int(n) + 1
		}
	}
}

func (m *Model) VocabSize() int {
	if m == nil || m.vocab == nil {
		return 0
	}
	return int(C.llama_vocab_n_tokens(m.vocab))
}

func (m *Model) Tokenize(text string, addSpecial, parseSpecial bool) ([]Token, error) {
	if m == nil || m.vocab == nil {
		return nil, newError("tokenize", "model vocabulary is nil")
	}

	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))

	size := len(text) + 8
	if size < 8 {
		size = 8
	}

	for {
		buf := make([]C.llama_token, size)
		n := C.llama_tokenize(
			m.vocab,
			ctext,
			C.int32_t(len(text)),
			(*C.llama_token)(unsafe.Pointer(&buf[0])),
			C.int32_t(len(buf)),
			C.bool(addSpecial),
			C.bool(parseSpecial),
		)
		if n >= 0 {
			tokens := make([]Token, int(n))
			for i := 0; i < int(n); i++ {
				tokens[i] = Token(buf[i])
			}
			return tokens, nil
		}

		size = int(-n)
		if size <= 0 {
			return nil, newError("tokenize", "llama_tokenize returned an invalid buffer requirement")
		}
	}
}

func (m *Model) TokenToPiece(token Token, lstrip int, special bool) (string, error) {
	if m == nil || m.vocab == nil {
		return "", newError("token_to_piece", "model vocabulary is nil")
	}

	size := 32
	for {
		buf := make([]C.char, size)
		n := C.llama_token_to_piece(
			m.vocab,
			C.llama_token(token),
			&buf[0],
			C.int32_t(len(buf)),
			C.int32_t(lstrip),
			C.bool(special),
		)
		if n >= 0 {
			return C.GoStringN(&buf[0], n), nil
		}

		size = int(-n)
		if size <= 0 {
			return "", newError("token_to_piece", fmt.Sprintf("token %d returned an invalid piece size", token))
		}
	}
}

func (m *Model) IsEndOfGeneration(token Token) bool {
	if m == nil || m.vocab == nil {
		return false
	}
	return bool(C.llama_vocab_is_eog(m.vocab, C.llama_token(token)))
}
