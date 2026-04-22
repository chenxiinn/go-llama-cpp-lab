package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/chenxiinn/go-llama-cpp-lab/pkg/chat"
)

type config struct {
	configPath string
	runtime    chat.RuntimeConfig
}

func newConfig() config {
	return config{
		configPath: chat.DefaultConfigPath,
		runtime:    chat.DefaultRuntimeConfig(),
	}
}

func (c *config) bindFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.configPath, "config", c.configPath, "Path to the local JSON config file.")
	c.runtime.BindFlags(fs)
}

func main() {
	cfg := newConfig()
	var explicitConfig bool
	cfg.configPath, explicitConfig = chat.ResolveConfigPath(os.Args[1:], cfg.configPath)
	var err error
	if explicitConfig {
		err = cfg.runtime.ApplyJSONFile(cfg.configPath)
	} else {
		err = cfg.runtime.ApplyOptionalJSONFile(cfg.configPath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "load chat config: %v\n", err)
		os.Exit(2)
	}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg.bindFlags(fs)
	fs.Parse(os.Args[1:])

	if err := cfg.runtime.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "invalid chat config: %v\n", err)
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "chat placeholder: native runtime lands in phase 1+\n")
	fmt.Fprintf(os.Stdout, "config=%q model=%q n_ctx=%d gpu_layers=%d max_tokens=%d\n",
		cfg.configPath,
		cfg.runtime.ModelPath,
		cfg.runtime.ContextSize,
		cfg.runtime.GPULayers,
		cfg.runtime.MaxTokens,
	)
}
