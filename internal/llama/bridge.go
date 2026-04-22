//go:build llama
// +build llama

package llama

/*
#include "bridge.h"
*/
import "C"

const (
	bridgeErrorBufferSize = 256
	systemInfoBufferSize  = 4096
)

type BackendProbe struct {
	MaxDevices int
	SystemInfo string
}

func ProbeBackend() (BackendProbe, error) {
	var errBuf [bridgeErrorBufferSize]C.char

	if code := C.bridge_llama_backend_init(&errBuf[0], C.size_t(len(errBuf))); code != 0 {
		return BackendProbe{}, newError("backend_init", C.GoString(&errBuf[0]))
	}
	defer C.bridge_llama_backend_free()

	maxDevices, err := bridgeMaxDevices()
	if err != nil {
		return BackendProbe{}, err
	}

	systemInfo, err := bridgeSystemInfo()
	if err != nil {
		return BackendProbe{}, err
	}

	return BackendProbe{
		MaxDevices: maxDevices,
		SystemInfo: systemInfo,
	}, nil
}

func bridgeMaxDevices() (int, error) {
	var errBuf [bridgeErrorBufferSize]C.char
	var out C.size_t

	if code := C.bridge_llama_max_devices(&out, &errBuf[0], C.size_t(len(errBuf))); code != 0 {
		return 0, newError("max_devices", C.GoString(&errBuf[0]))
	}

	return int(out), nil
}

func bridgeSystemInfo() (string, error) {
	var errBuf [bridgeErrorBufferSize]C.char
	var buf [systemInfoBufferSize]C.char

	if code := C.bridge_llama_print_system_info(&buf[0], C.size_t(len(buf)), &errBuf[0], C.size_t(len(errBuf))); code != 0 {
		return "", newError("print_system_info", C.GoString(&errBuf[0]))
	}

	return C.GoString(&buf[0]), nil
}
