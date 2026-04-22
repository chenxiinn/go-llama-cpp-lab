package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chenxiinn/go-llama-cpp-lab/pkg/chat"
)

type config struct {
	runtime chat.RuntimeConfig
}

func newConfig() config {
	return config{
		runtime: chat.DefaultRuntimeConfig(),
	}
}

func (c *config) bindFlags(fs *flag.FlagSet) {
	c.runtime.BindFlags(fs)
}

func main() {
	cfg := newConfig()

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg.bindFlags(fs)
	fs.Parse(os.Args[1:])

	if err := cfg.runtime.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid chat config: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "chat placeholder: native runtime lands in phase 1+\n")
	fmt.Fprintf(os.Stdout, "model=%q n_ctx=%d gpu_layers=%d max_tokens=%d\n",
		cfg.runtime.ModelPath,
		cfg.runtime.ContextSize,
		cfg.runtime.GPULayers,
		cfg.runtime.MaxTokens,
	)
}
