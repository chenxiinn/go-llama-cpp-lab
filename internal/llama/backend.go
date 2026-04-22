//go:build llama
// +build llama

package llama

/*
#include "llama.h"
*/
import "C"

import "sync"

var backendInitOnce sync.Once

func ensureBackendInitialized() {
	backendInitOnce.Do(func() {
		C.llama_backend_init()
	})
}
