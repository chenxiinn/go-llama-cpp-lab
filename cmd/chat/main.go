package main

import (
	"fmt"
	"os"

	"github.com/chenxiinn/go-llama-cpp-lab/internal/appconfig"
)

type config struct {
	bootstrap appconfig.RuntimeBootstrap
}

func newConfig() config {
	return config{
		bootstrap: appconfig.NewRuntimeBootstrap(),
	}
}

func main() {
	cfg := newConfig()
	loadedPath, err := cfg.bootstrap.LoadAndParse(os.Args[0], os.Args[1:], nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load chat config: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "chat placeholder: native runtime lands in phase 1+\n")
	fmt.Fprintf(os.Stdout, "config=%q model=%q n_ctx=%d gpu_layers=%d max_tokens=%d\n",
		loadedPath,
		cfg.bootstrap.Runtime.ModelPath,
		cfg.bootstrap.Runtime.ContextSize,
		cfg.bootstrap.Runtime.GPULayers,
		cfg.bootstrap.Runtime.MaxTokens,
	)
}
