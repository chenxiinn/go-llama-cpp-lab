package chat

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfigPathExplicit(t *testing.T) {
	path, explicit := ResolveConfigPath([]string{"-config", "config/dev.json"}, DefaultConfigPath)
	if !explicit {
		t.Fatal("expected explicit config path")
	}
	if path != "config/dev.json" {
		t.Fatalf("path = %q, want %q", path, "config/dev.json")
	}
}

func TestResolveConfigPathDefault(t *testing.T) {
	path, explicit := ResolveConfigPath(nil, DefaultConfigPath)
	if explicit {
		t.Fatal("expected default config path")
	}
	if path != DefaultConfigPath {
		t.Fatalf("path = %q, want %q", path, DefaultConfigPath)
	}
}

func TestApplyJSONFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "local.json")
	if err := os.WriteFile(path, []byte(`{"model_path":"test.gguf","n_ctx":8192,"gpu_layers":8,"max_tokens":512}`), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	cfg := DefaultRuntimeConfig()
	if err := cfg.ApplyJSONFile(path); err != nil {
		t.Fatalf("apply config file: %v", err)
	}

	if cfg.ModelPath != "test.gguf" {
		t.Fatalf("model path = %q, want %q", cfg.ModelPath, "test.gguf")
	}
	if cfg.ContextSize != 8192 {
		t.Fatalf("context size = %d, want 8192", cfg.ContextSize)
	}
	if cfg.GPULayers != 8 {
		t.Fatalf("gpu layers = %d, want 8", cfg.GPULayers)
	}
	if cfg.MaxTokens != 512 {
		t.Fatalf("max tokens = %d, want 512", cfg.MaxTokens)
	}
}
