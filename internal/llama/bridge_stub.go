//go:build !llama
// +build !llama

package llama

type BackendProbe struct {
	MaxDevices int
	SystemInfo string
}

func ProbeBackend() (BackendProbe, error) {
	return BackendProbe{}, newError(
		"probe_backend",
		"native bridge disabled; build with -tags llama and set CGO_CFLAGS/CGO_LDFLAGS for local libllama headers and libraries",
	)
}
