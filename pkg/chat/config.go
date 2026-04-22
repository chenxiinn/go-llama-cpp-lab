package chat

import (
	"errors"
	"flag"
)

// RuntimeConfig is the stable runtime flag surface shared by the CLI and HTTP
// entrypoints. Native handles and chat session state will plug into this in
// later phases.
type RuntimeConfig struct {
	ModelPath   string `json:"model_path"`
	ContextSize int    `json:"n_ctx"`
	GPULayers   int    `json:"gpu_layers"`
	MaxTokens   int    `json:"max_tokens"`
}

func DefaultRuntimeConfig() RuntimeConfig {
	return RuntimeConfig{
		ContextSize: 4096,
		GPULayers:   0,
		MaxTokens:   256,
	}
}

func (c *RuntimeConfig) BindFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.ModelPath, "model", c.ModelPath, "Path to the local GGUF model.")
	fs.IntVar(&c.ContextSize, "n-ctx", c.ContextSize, "Context window size for the runtime.")
	fs.IntVar(&c.GPULayers, "gpu-layers", c.GPULayers, "Number of llama.cpp layers to place on GPU.")
	fs.IntVar(&c.MaxTokens, "max-tokens", c.MaxTokens, "Maximum generated tokens per response.")
}

func (c RuntimeConfig) Validate() error {
	if c.ModelPath == "" {
		return errors.New("model path must be set")
	}
	if c.ContextSize <= 0 {
		return errors.New("n-ctx must be positive")
	}
	if c.GPULayers < 0 {
		return errors.New("gpu-layers must be zero or greater")
	}
	if c.MaxTokens <= 0 {
		return errors.New("max-tokens must be positive")
	}
	return nil
}
