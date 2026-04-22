package appconfig

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/chenxiinn/go-llama-cpp-lab/pkg/chat"
)

func TestResolveExplicitConfigPath(t *testing.T) {
	path, explicit, err := ResolveExplicitConfigPath([]string{"--config", "config/dev.json"})
	if err != nil {
		t.Fatalf("resolve explicit config path: %v", err)
	}
	if !explicit {
		t.Fatal("expected explicit config path")
	}
	if path != "config/dev.json" {
		t.Fatalf("path = %q, want %q", path, "config/dev.json")
	}
}

func TestResolveExplicitConfigPathMissingValue(t *testing.T) {
	_, _, err := ResolveExplicitConfigPath([]string{"--config"})
	if err == nil {
		t.Fatal("expected missing config value error")
	}
}

func TestFindDefaultConfigPathPrefersLocal(t *testing.T) {
	workDir := t.TempDir()
	homeDir := t.TempDir()
	localPath := filepath.Join(workDir, LocalConfigPath)
	globalPath := UserConfigPath(homeDir)

	mustWriteFile(t, localPath, `{"model_path":"local.gguf"}`)
	mustWriteFile(t, globalPath, `{"model_path":"global.gguf"}`)

	got, err := FindDefaultConfigPath(workDir, homeDir)
	if err != nil {
		t.Fatalf("find default config path: %v", err)
	}
	if got != localPath {
		t.Fatalf("path = %q, want %q", got, localPath)
	}
}

func TestFindDefaultConfigPathFallsBackToUserConfig(t *testing.T) {
	workDir := t.TempDir()
	homeDir := t.TempDir()
	globalPath := UserConfigPath(homeDir)

	mustWriteFile(t, globalPath, `{"model_path":"global.gguf"}`)

	got, err := FindDefaultConfigPath(workDir, homeDir)
	if err != nil {
		t.Fatalf("find default config path: %v", err)
	}
	if got != globalPath {
		t.Fatalf("path = %q, want %q", got, globalPath)
	}
}

func TestFindDefaultConfigPathErrorsWhenMissing(t *testing.T) {
	workDir := t.TempDir()
	homeDir := t.TempDir()

	_, err := FindDefaultConfigPath(workDir, homeDir)
	if err == nil {
		t.Fatal("expected missing config error")
	}
	if !strings.Contains(err.Error(), filepath.Join(workDir, LocalConfigPath)) {
		t.Fatalf("error = %q, missing local candidate", err)
	}
	if !strings.Contains(err.Error(), UserConfigPath(homeDir)) {
		t.Fatalf("error = %q, missing global candidate", err)
	}
}

func TestRuntimeBootstrapLoadUsesExplicitConfigOnly(t *testing.T) {
	workDir := t.TempDir()
	homeDir := t.TempDir()
	explicitPath := filepath.Join(workDir, "custom.json")
	localPath := filepath.Join(workDir, LocalConfigPath)

	mustWriteFile(t, explicitPath, `{"model_path":"explicit.gguf","n_ctx":8192}`)
	mustWriteFile(t, localPath, `{"model_path":"local.gguf","n_ctx":4096}`)

	bootstrap := RuntimeBootstrap{
		HomeDir: homeDir,
		Runtime: chat.DefaultRuntimeConfig(),
		WorkDir: workDir,
	}

	loadedPath, err := bootstrap.Load([]string{"--config", "custom.json"})
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if loadedPath != explicitPath {
		t.Fatalf("loaded path = %q, want %q", loadedPath, explicitPath)
	}
	if bootstrap.Runtime.ModelPath != "explicit.gguf" {
		t.Fatalf("model path = %q, want %q", bootstrap.Runtime.ModelPath, "explicit.gguf")
	}
	if bootstrap.Runtime.ContextSize != 8192 {
		t.Fatalf("context size = %d, want 8192", bootstrap.Runtime.ContextSize)
	}
}

func TestRuntimeBootstrapLoadErrorsForMissingExplicitConfig(t *testing.T) {
	bootstrap := RuntimeBootstrap{
		HomeDir: t.TempDir(),
		Runtime: chat.DefaultRuntimeConfig(),
		WorkDir: t.TempDir(),
	}

	_, err := bootstrap.Load([]string{"--config", "missing.json"})
	if err == nil {
		t.Fatal("expected missing explicit config error")
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir %q: %v", filepath.Dir(path), err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %q: %v", path, err)
	}
}
