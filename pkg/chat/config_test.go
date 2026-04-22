package chat

import "testing"

func TestRuntimeConfigValidateRequiresModelPath(t *testing.T) {
	cfg := DefaultRuntimeConfig()

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected missing model path error")
	}
}
